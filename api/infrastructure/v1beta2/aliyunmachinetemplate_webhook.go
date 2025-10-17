package v1beta2

import (
	"errors"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var mtmpLog = logf.Log.WithName("aliyunmachinetemplate-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks.
func (r *AliyunMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-aliyunmachinetemplate,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=aliyunmachinetemplates,verbs=create;update,versions=v1beta2,name=maliyunmachinetemplate.kb.io,admissionReviewVersions=v1
var _ webhook.Defaulter = &AliyunMachineTemplate{}

// +kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-aliyunmachinetemplate,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=aliyunmachinetemplates,verbs=create;update,versions=v1beta2,name=valiyunmachinetemplate.kb.io,admissionReviewVersions=v1
var _ webhook.Validator = &AliyunMachineTemplate{}

var (
	zeroI32 int32 = 0

	defSystemDiskSizeI32 int32 = 120
	minSystemDiskSizeI32 int32 = 40
	maxSystemDiskSizeI32 int32 = 2048

	minDataDiskSizeI32 int32 = 40
	maxDataDiskSizeI32 int32 = 32768

	defInstanceChargeType = "PostPaid"     // 预设实例按量付费
	defInternetChargeType = "PayByTraffic" // 预设宽带按流量付费
)

// Defaulting
func (r *AliyunMachineTemplate) Default() {
	mtmpLog.V(5).Info("default", "name", r.Name)

	spec := &r.Spec.Template.Spec

	// 计费相关默认
	if spec.InstanceChargeType == nil || *spec.InstanceChargeType == "" {
		mtmpLog.Info("InstanceChargeType empty, use default", "value", defInstanceChargeType)
		spec.InstanceChargeType = &defInstanceChargeType
	}
	if spec.InternetChargeType == nil || *spec.InternetChargeType == "" {
		mtmpLog.Info("InternetChargeType empty, use default", "value", defInternetChargeType)
		spec.InternetChargeType = &defInternetChargeType
	}

	// 系统盘默认
	if spec.SystemDisk == nil {
		spec.SystemDisk = &SystemDiskSpec{}
	}
	if spec.SystemDisk.Category == nil || *spec.SystemDisk.Category == "" {
		def := defSystemDiskCategory
		mtmpLog.Info("SystemDisk.Category empty, use default", "value", def)
		spec.SystemDisk.Category = &def
	}
	if spec.SystemDisk.Size == nil || *spec.SystemDisk.Size == zeroI32 {
		def := defSystemDiskSizeI32
		mtmpLog.Info("SystemDisk.Size empty, use default", "value", def)
		spec.SystemDisk.Size = &def
	}
	if spec.SystemDisk.PerformanceLevel == nil || *spec.SystemDisk.PerformanceLevel == "" {
		def := defSystemDiskPerformanceLevel
		mtmpLog.Info("SystemDisk.PerformanceLevel empty, use default", "value", def)
		spec.SystemDisk.PerformanceLevel = &def
	}
}

func (r *AliyunMachineTemplate) validateRequiredFields() field.ErrorList {
	var allErrs field.ErrorList
	p := field.NewPath("spec", "template", "spec")

	spec := r.Spec.Template.Spec

	if spec.RegionID == "" {
		allErrs = append(allErrs, field.Required(p.Child("regionId"), "regionId is required"))
	}
	if spec.VSwitchID == "" {
		allErrs = append(allErrs, field.Required(p.Child("vSwitchId"), "vSwitchId is required"))
	}
	if spec.InstanceType == "" {
		allErrs = append(allErrs, field.Required(p.Child("instanceType"), "instanceType is required"))
	}
	return allErrs
}

func (r *AliyunMachineTemplate) validateSystemDisk() field.ErrorList {
	var allErrs field.ErrorList
	p := field.NewPath("spec", "template", "spec", "systemDisk")
	sd := r.Spec.Template.Spec.SystemDisk
	if sd == nil {
		// 创建逻辑会在创建时补上默认的 systemdisk
		return allErrs
	}
	if sd.Category != nil && !systemDiskCategories.Has(*sd.Category) {
		allErrs = append(allErrs, field.Invalid(p.Child("category"), sd.Category, fmt.Sprintf("category must be one of %v", systemDiskCategories.List())))
	}
	if sd.PerformanceLevel != nil && !systemDiskPerformanceLevels.Has(*sd.PerformanceLevel) {
		allErrs = append(allErrs, field.Invalid(p.Child("performanceLevel"), sd.PerformanceLevel, fmt.Sprintf("performanceLevel must be one of %v", systemDiskPerformanceLevels.List())))
	}
	if sd.Size != nil && (*sd.Size < minSystemDiskSizeI32 || *sd.Size > maxSystemDiskSizeI32) {
		allErrs = append(allErrs, field.Invalid(p.Child("size"), sd.Size, fmt.Sprintf("size must be in [%d, %d] GB", minSystemDiskSizeI32, maxSystemDiskSizeI32)))
	}
	return allErrs
}

func (r *AliyunMachineTemplate) validateDataDisks() field.ErrorList {
	var allErrs field.ErrorList
	base := field.NewPath("spec", "template", "spec", "dataDisks")
	for i := range r.Spec.Template.Spec.DataDisks {
		p := base.Index(i)
		dd := r.Spec.Template.Spec.DataDisks[i]
		if dd.Category != nil && !dataDiskCategories.Has(*dd.Category) {
			allErrs = append(allErrs, field.Invalid(p.Child("category"), dd.Category, fmt.Sprintf("category must be one of %v", dataDiskCategories.List())))
		}
		if dd.PerformanceLevel != nil && !dataDiskPerformanceLevels.Has(*dd.PerformanceLevel) {
			allErrs = append(allErrs, field.Invalid(p.Child("performanceLevel"), dd.PerformanceLevel, fmt.Sprintf("performanceLevel must be one of %v", dataDiskPerformanceLevels.List())))
		}
		if dd.Size != nil && (*dd.Size < minDataDiskSizeI32 || *dd.Size > maxDataDiskSizeI32) {
			allErrs = append(allErrs, field.Invalid(p.Child("size"), dd.Size, fmt.Sprintf("size must be in [%d, %d] GB", minDataDiskSizeI32, maxDataDiskSizeI32)))
		}
	}
	return allErrs
}

func (r *AliyunMachineTemplate) validateNetworkingAndCharge() field.ErrorList {
	var allErrs field.ErrorList
	p := field.NewPath("spec", "template", "spec")
	spec := r.Spec.Template.Spec

	// 公网带宽限制
	if spec.InternetMaxBandwidthOut != nil && *spec.InternetMaxBandwidthOut < 0 {
		allErrs = append(allErrs, field.Invalid(p.Child("internetMaxBandwidthOut"), *spec.InternetMaxBandwidthOut, "must be >= 0"))
	}
	if spec.InternetMaxBandwidthIn != nil && *spec.InternetMaxBandwidthIn < 0 {
		allErrs = append(allErrs, field.Invalid(p.Child("internetMaxBandwidthIn"), *spec.InternetMaxBandwidthIn, "must be >= 0"))
	}
	return allErrs
}

// Validator
// ValidateCreate implements webhook.Validator
func (r *AliyunMachineTemplate) ValidateCreate() (admission.Warnings, error) {
	mtmpLog.V(5).Info("validate create", "name", r.Name)

	var allErrs field.ErrorList
	allErrs = append(allErrs, r.validateRequiredFields()...)
	allErrs = append(allErrs, r.validateSystemDisk()...)
	allErrs = append(allErrs, r.validateDataDisks()...)
	allErrs = append(allErrs, r.validateNetworkingAndCharge()...)

	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate implements webhook.Validator
func (r *AliyunMachineTemplate) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	mtmpLog.V(5).Info("validate update", "name", r.Name)

	oldObj, ok := old.(*AliyunMachineTemplate)
	if !ok {
		return nil, apierrors.NewInvalid(
			GroupVersion.WithKind(AliyunMachineTemplateKind).GroupKind(),
			r.Name,
			field.ErrorList{field.InternalError(nil, errors.New("failed to convert old AliyunMachineTemplate to object"))},
		)
	}

	var allErrs field.ErrorList
	allErrs = append(allErrs, r.validateSystemDisk()...)
	allErrs = append(allErrs, r.validateDataDisks()...)
	allErrs = append(allErrs, r.validateNetworkingAndCharge()...)

	// 建议在模板中将以下字段视为不可变：
	p := field.NewPath("spec", "template", "spec")
	newSpec := r.Spec.Template.Spec
	oldSpec := oldObj.Spec.Template.Spec

	// 不可以改变 RegionID/VSwitchID/InstanceType/ImageID
	if newSpec.RegionID != oldSpec.RegionID {
		allErrs = append(allErrs, field.Invalid(p.Child("regionId"), newSpec.RegionID, "field is immutable"))
	}
	if newSpec.VSwitchID != oldSpec.VSwitchID {
		allErrs = append(allErrs, field.Invalid(p.Child("vSwitchId"), newSpec.VSwitchID, "field is immutable"))
	}
	if newSpec.InstanceType != oldSpec.InstanceType {
		allErrs = append(allErrs, field.Invalid(p.Child("instanceType"), newSpec.InstanceType, "field is immutable"))
	}
	if newSpec.ImageID != nil && oldSpec.ImageID != nil && *newSpec.ImageID != *oldSpec.ImageID {
		allErrs = append(allErrs, field.Invalid(p.Child("imageId"), newSpec.ImageID, "field is immutable"))
	}

	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateDelete implements webhook.Validator
func (r *AliyunMachineTemplate) ValidateDelete() (admission.Warnings, error) {
	mtmpLog.V(5).Info("validate delete", "name", r.Name)
	return nil, nil
}
