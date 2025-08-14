/*
*Copyright (c) 2024-2025, Alibaba Cloud and its affiliates;
*Licensed under the Apache License, Version 2.0 (the "License");
*you may not use this file except in compliance with the License.
*You may obtain a copy of the License at

*   http://www.apache.org/licenses/LICENSE-2.0

*Unless required by applicable law or agreed to in writing, software
*distributed under the License is distributed on an "AS IS" BASIS,
*WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*See the License for the specific language governing permissions and
*limitations under the License.
 */

package infrastructure

import (
	essv1alpha1 "github.com/AliyunContainerService/alibabacloud-provider-for-Cluster-API/api/ess/v1alpha1"
	ess20220222 "github.com/alibabacloud-go/ess-20220222/v2/client"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// getOwnerMachinePool 查询目标资源 obj (aliyunPool)的 MachinePool 属主并返回.
//
//	@param obj: aliyunPool.ObjectMeta 成员字段
//
// getOwnerMachinePool returns the MachinePool object owning the current resource.
//func getOwnerMachinePool(ctx context.Context, c client.Client, obj metav1.ObjectMeta) (*expclusterv1.MachinePool, error) {
//	for _, ref := range obj.OwnerReferences {
//		if ref.Kind != "MachinePool" {
//			continue
//		}
//		gv, err := schema.ParseGroupVersion(ref.APIVersion)
//		if err != nil {
//			return nil, errors.WithStack(err)
//		}
//		if gv.Group == expclusterv1.GroupVersion.Group {
//			return getMachinePoolByName(ctx, c, obj.Namespace, ref.Name)
//		}
//	}
//	return nil, nil
//}

func str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
func f(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}
func ptrIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func stringSlicePtrEqual(a, b []*string) bool {
	if len(a) != len(b) {
		return false
	}
	// 无序比较
	set := map[string]int{}
	for _, p := range a {
		set[str(p)]++
	}
	for _, p := range b {
		set[str(p)]--
	}
	for _, v := range set {
		if v != 0 {
			return false
		}
	}
	return true
}

func f64Ptr(v int) *float64 { f := float64(v); return &f }

func toPtrSlice(ss []string) []*string {
	out := make([]*string, 0, len(ss))
	for i := range ss {
		s := ss[i]
		out = append(out, &s)
	}
	return out
}
func providerRefEqual(a, b *xpv1.Reference) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Name == b.Name
}

// 只比较我们关心的字段，避免被 Upjet 的内部回写干扰幂等性
func equalScalingGroupParams(a, b essv1alpha1.ScalingGroupParameters) bool {
	if str(a.ScalingGroupName) != str(b.ScalingGroupName) {
		return false
	}
	if f(a.MinSize) != f(b.MinSize) {
		return false
	}
	if f(a.MaxSize) != f(b.MaxSize) {
		return false
	}
	if f(a.DesiredCapacity) != f(b.DesiredCapacity) {
		return false
	}
	if !stringSlicePtrEqual(a.VswitchIds, b.VswitchIds) {
		return false
	}
	if !stringSlicePtrEqual(a.LoadbalancerIds, b.LoadbalancerIds) {
		return false
	}
	return true
}
func equalScalingConfigurationParamsSliceAware(a, b essv1alpha1.ScalingConfigurationParameters) bool {
	if str(a.ScalingGroupID) != str(b.ScalingGroupID) {
		return false
	}
	if str(a.ImageID) != str(b.ImageID) {
		return false
	}
	if !stringSlicePtrEqual(a.InstanceTypes, b.InstanceTypes) {
		return false
	}
	if !stringSlicePtrEqual(a.SecurityGroupIds, b.SecurityGroupIds) {
		return false
	}
	if str(a.UserData) != str(b.UserData) {
		return false
	}
	if str(a.InternetChargeType) != str(b.InternetChargeType) {
		return false
	}
	if f(a.InternetMaxBandwidthIn) != f(b.InternetMaxBandwidthIn) {
		return false
	}
	if f(a.InternetMaxBandwidthOut) != f(b.InternetMaxBandwidthOut) {
		return false
	}
	return true
}

func calcReplicas(instances []*ess20220222.DescribeScalingInstancesResponseBodyScalingInstances) (int32, int32, error) {
	replicas := int32(len(instances))
	ready := int32(0)
	for _, inst := range instances {
		if inst.HealthStatus == nil || *inst.HealthStatus != "Healthy" {
			continue
		}
		if inst.InstanceId == nil {
			continue
		}
		ready++
	}
	return replicas, ready, nil
}
