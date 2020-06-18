// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package virtualmachine_test

import (
	"context"
	"time"

	azurev1alpha1 "github.wdf.sap.corp/kubernetes/remedy-controller/pkg/apis/azure/v1alpha1"
	"github.wdf.sap.corp/kubernetes/remedy-controller/pkg/apis/config"
	"github.wdf.sap.corp/kubernetes/remedy-controller/pkg/controller"
	"github.wdf.sap.corp/kubernetes/remedy-controller/pkg/controller/azure/virtualmachine"
	mockclient "github.wdf.sap.corp/kubernetes/remedy-controller/pkg/mock/controller-runtime/client"
	mockprometheus "github.wdf.sap.corp/kubernetes/remedy-controller/pkg/mock/prometheus"
	mockutilsazure "github.wdf.sap.corp/kubernetes/remedy-controller/pkg/mock/remedy-controller/utils/azure"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/go-logr/logr"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

var _ = Describe("Actuator", func() {
	const (
		nodeName                = "shoot--dev--test-vm1"
		hostname                = "shoot--dev--test-vm1"
		providerID              = "azure:///subscriptions/xxx/resourceGroups/shoot--dev--test/providers/Microsoft.Compute/virtualMachines/shoot--dev--test-vm1"
		azureVirtualMachineID   = "/subscriptions/xxx/resourceGroups/shoot--dev--test/providers/Microsoft.Compute/virtualMachines/shoot--dev--test-vm1"
		azureVirtualMachineName = "shoot--dev--test-vm1"
	)

	var (
		ctrl *gomock.Controller
		ctx  context.Context

		c                   *mockclient.MockClient
		sw                  *mockclient.MockStatusWriter
		vmUtils             *mockutilsazure.MockVirtualMachineUtils
		reappliedVMsCounter *mockprometheus.MockCounter

		cfg      config.AzureFailedVMRemedyConfiguration
		logger   logr.Logger
		actuator controller.Actuator

		newVM                  func(bool, bool, bool, compute.ProvisioningState) *azurev1alpha1.VirtualMachine
		newAzureVirtualMachine func(compute.ProvisioningState) *compute.VirtualMachine
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		ctx = context.TODO()

		c = mockclient.NewMockClient(ctrl)
		sw = mockclient.NewMockStatusWriter(ctrl)
		c.EXPECT().Status().Return(sw).AnyTimes()
		vmUtils = mockutilsazure.NewMockVirtualMachineUtils(ctrl)
		reappliedVMsCounter = mockprometheus.NewMockCounter(ctrl)

		cfg = config.AzureFailedVMRemedyConfiguration{
			RequeueInterval: metav1.Duration{Duration: 1 * time.Second},
		}
		logger = log.Log.WithName("test")
		actuator = virtualmachine.NewActuator(vmUtils, cfg, logger, reappliedVMsCounter)
		Expect(actuator.(inject.Client).InjectClient(c)).To(Succeed())

		newVM = func(ready, unreachable, withStatus bool, provisioningState compute.ProvisioningState) *azurev1alpha1.VirtualMachine {
			var status azurev1alpha1.VirtualMachineStatus
			if withStatus {
				status = azurev1alpha1.VirtualMachineStatus{
					Exists:            true,
					ID:                pointer.StringPtr(azureVirtualMachineID),
					Name:              pointer.StringPtr(azureVirtualMachineName),
					ProvisioningState: pointer.StringPtr(string(provisioningState)),
				}
			}
			return &azurev1alpha1.VirtualMachine{
				ObjectMeta: metav1.ObjectMeta{
					Name: nodeName,
				},
				Spec: azurev1alpha1.VirtualMachineSpec{
					Hostname:    hostname,
					ProviderID:  providerID,
					Ready:       ready,
					Unreachable: unreachable,
				},
				Status: status,
			}
		}
		newAzureVirtualMachine = func(provisioningState compute.ProvisioningState) *compute.VirtualMachine {
			return &compute.VirtualMachine{
				ID:   pointer.StringPtr(azureVirtualMachineID),
				Name: pointer.StringPtr(azureVirtualMachineName),
				VirtualMachineProperties: &compute.VirtualMachineProperties{
					ProvisioningState: pointer.StringPtr(string(provisioningState)),
				},
			}
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#CreateOrUpdate", func() {
		It("should update the VirtualMachine object status if the VM is found", func() {
			vm := newVM(true, false, false, "")
			vmWithStatus := newVM(true, false, true, compute.ProvisioningStateSucceeded)
			azureVirtualMachine := newAzureVirtualMachine(compute.ProvisioningStateSucceeded)
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(azureVirtualMachine, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vm.Name}, vm).Return(nil)
			sw.EXPECT().Update(ctx, vmWithStatus).Return(nil)

			requeueAfter, removeFinalizer, err := actuator.CreateOrUpdate(ctx, vm)
			Expect(err).NotTo(HaveOccurred())
			Expect(requeueAfter).To(Equal(time.Duration(0)))
			Expect(removeFinalizer).To(Equal(false))
		})

		It("should not update the VirtualMachine object status if the VM is not found", func() {
			vm := newVM(true, false, false, "")
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(nil, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vm.Name}, vm).Return(nil)

			requeueAfter, removeFinalizer, err := actuator.CreateOrUpdate(ctx, vm)
			Expect(err).NotTo(HaveOccurred())
			Expect(requeueAfter).To(Equal(1 * time.Second))
			Expect(removeFinalizer).To(Equal(false))
		})

		It("should not update the VirtualMachine object status if the VM is found and the status is already initialized", func() {
			vmWithStatus := newVM(true, false, true, compute.ProvisioningStateSucceeded)
			azureVirtualMachine := newAzureVirtualMachine(compute.ProvisioningStateSucceeded)
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(azureVirtualMachine, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vmWithStatus.Name}, vmWithStatus).Return(nil)

			requeueAfter, removeFinalizer, err := actuator.CreateOrUpdate(ctx, vmWithStatus)
			Expect(err).NotTo(HaveOccurred())
			Expect(requeueAfter).To(Equal(time.Duration(0)))
			Expect(removeFinalizer).To(Equal(false))
		})

		It("should update the VirtualMachine object status if the VM is not found and the status is already initialized", func() {
			vm := newVM(true, false, false, "")
			vmWithStatus := newVM(true, false, true, compute.ProvisioningStateSucceeded)
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(nil, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vmWithStatus.Name}, vmWithStatus).Return(nil)
			sw.EXPECT().Update(ctx, vm).Return(nil)

			requeueAfter, removeFinalizer, err := actuator.CreateOrUpdate(ctx, vmWithStatus)
			Expect(err).NotTo(HaveOccurred())
			Expect(requeueAfter).To(Equal(1 * time.Second))
			Expect(removeFinalizer).To(Equal(false))
		})

		It("should reapply the Azure VM if it's in a failed state", func() {
			vm := newVM(false, true, false, "")
			vmWithStatus := newVM(false, true, true, compute.ProvisioningStateFailed)
			azureVirtualMachine := newAzureVirtualMachine(compute.ProvisioningStateFailed)
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(azureVirtualMachine, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vm.Name}, vm).Return(nil)
			sw.EXPECT().Update(ctx, vmWithStatus).Return(nil)
			vmUtils.EXPECT().Reapply(ctx, azureVirtualMachineName).Return(nil)
			reappliedVMsCounter.EXPECT().Inc()

			requeueAfter, removeFinalizer, err := actuator.CreateOrUpdate(ctx, vm)
			Expect(err).NotTo(HaveOccurred())
			Expect(requeueAfter).To(Equal(time.Duration(0)))
			Expect(removeFinalizer).To(Equal(false))
		})

		It("should fail if getting the Azure VM fails", func() {
			vm := newVM(true, false, false, "")
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(nil, errors.New("test"))

			_, _, err := actuator.CreateOrUpdate(ctx, vm)
			Expect(err).To(MatchError("could not get Azure virtual machine: test"))
		})

		It("should fail if updating the VirtualMachine object status fails", func() {
			vm := newVM(true, false, false, "")
			vmWithStatus := newVM(true, false, true, compute.ProvisioningStateSucceeded)
			azureVirtualMachine := newAzureVirtualMachine(compute.ProvisioningStateSucceeded)
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(azureVirtualMachine, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vm.Name}, vm).Return(nil)
			sw.EXPECT().Update(ctx, vmWithStatus).Return(errors.New("test"))

			_, _, err := actuator.CreateOrUpdate(ctx, vm)
			Expect(err).To(MatchError("could not update virtualmachine status: test"))
		})

		It("should fail if reapplying the Azure VM fails", func() {
			vm := newVM(false, true, false, "")
			vmWithStatus := newVM(false, true, true, compute.ProvisioningStateFailed)
			azureVirtualMachine := newAzureVirtualMachine(compute.ProvisioningStateFailed)
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(azureVirtualMachine, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vm.Name}, vm).Return(nil)
			sw.EXPECT().Update(ctx, vmWithStatus).Return(nil)
			vmUtils.EXPECT().Reapply(ctx, azureVirtualMachineName).Return(errors.New("test"))

			_, _, err := actuator.CreateOrUpdate(ctx, vm)
			Expect(err).To(MatchError("could not reapply Azure virtual machine: test"))
		})
	})

	Describe("#Delete", func() {
		It("should update the VirtualMachine object status if the VM is found", func() {
			vm := newVM(true, false, false, "")
			vmWithStatus := newVM(true, false, true, compute.ProvisioningStateSucceeded)
			azureVirtualMachine := newAzureVirtualMachine(compute.ProvisioningStateSucceeded)
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(azureVirtualMachine, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vm.Name}, vm).Return(nil)
			sw.EXPECT().Update(ctx, vmWithStatus).Return(nil)

			Expect(actuator.Delete(ctx, vm)).To(Succeed())
		})

		It("should not update the VirtualMachine object status if the VM is not found", func() {
			vm := newVM(true, false, false, "")
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(nil, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vm.Name}, vm).Return(nil)

			Expect(actuator.Delete(ctx, vm)).To(Succeed())
		})

		It("should not update the VirtualMachine object status if the VM is found and the status is already initialized", func() {
			vmWithStatus := newVM(true, false, true, compute.ProvisioningStateSucceeded)
			azureVirtualMachine := newAzureVirtualMachine(compute.ProvisioningStateSucceeded)
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(azureVirtualMachine, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vmWithStatus.Name}, vmWithStatus).Return(nil)

			Expect(actuator.Delete(ctx, vmWithStatus)).To(Succeed())
		})

		It("should update the VirtualMachine object status if the VM is not found and the status is already initialized", func() {
			vm := newVM(true, false, false, "")
			vmWithStatus := newVM(true, false, true, compute.ProvisioningStateSucceeded)
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(nil, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vmWithStatus.Name}, vmWithStatus).Return(nil)
			sw.EXPECT().Update(ctx, vm).Return(nil)

			Expect(actuator.Delete(ctx, vmWithStatus)).To(Succeed())
		})

		It("should fail if getting the Azure VM fails", func() {
			vm := newVM(true, false, false, "")
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(nil, errors.New("test"))

			Expect(actuator.Delete(ctx, vm)).To(MatchError("could not get Azure virtual machine: test"))
		})

		It("should fail if updating the VirtualMachine object status fails", func() {
			vm := newVM(true, false, false, "")
			vmWithStatus := newVM(true, false, true, compute.ProvisioningStateSucceeded)
			azureVirtualMachine := newAzureVirtualMachine(compute.ProvisioningStateSucceeded)
			vmUtils.EXPECT().Get(ctx, azureVirtualMachineName).Return(azureVirtualMachine, nil)
			c.EXPECT().Get(ctx, client.ObjectKey{Name: vm.Name}, vm).Return(nil)
			sw.EXPECT().Update(ctx, vmWithStatus).Return(errors.New("test"))

			Expect(actuator.Delete(ctx, vm)).To(MatchError("could not update virtualmachine status: test"))
		})
	})
})