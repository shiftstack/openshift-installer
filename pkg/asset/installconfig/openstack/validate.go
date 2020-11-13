package openstack

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/asset/installconfig/openstack/validation"
	"github.com/openshift/installer/pkg/types"
)

// Validate validates the given installconfig for OpenStack platform
func Validate(ic *types.InstallConfig) error {
	ci, err := validation.GetCloudInfo(ic)
	if err != nil {
		return err
	}

	allErrs := field.ErrorList{}

	// NOTES:
	//
	// the idea is to pass ic.Platform.OpenStack to ValidateMachinePool() here
	// and in validation/machinepool.go ValidateMachinePool(), we would do the following logic:
	// if platform has rootVolume:
	// 		checkStorageFlavor = false
	//		validateFlavor(mpFlavor, checkStorageFlavor)
	// elseif platform has flavor:
	// 		checkStorageFlavor = true
	//		validateFlavor(mpFlavor, checkStorageFlavor)
	// else:
	//		validateFlavor(platformFlavor)
	//
	// Therefore, we would remove the validateFlavor call from ValidatePlatform() (not necessary anymore)
	//
	// Concerns / questions:
	// - so we would validate platform flavor multiple times, if rootVolume wasn't used
	// - what if no machine pool is defined? Where do we run the flavor validation?
	allErrs = append(allErrs, validation.ValidatePlatform(ic.Platform.OpenStack, ic.Networking, ci)...)

	// areMachinePools bool
	if ic.ControlPlane.Platform.OpenStack != nil {
		// areMachinePools = true
		allErrs = append(allErrs, validation.ValidateMachinePool(ic.ControlPlane.Platform.OpenStack, ci, true, field.NewPath("controlPlane", "platform", "openstack"))...)
	}
	for idx, compute := range ic.Compute {
		// areMachinePools = true
		// check if rootVolume or Flavor defined:
		// if yes: we don't check platform flavor
		// if no:l we check
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.OpenStack != nil {
			allErrs = append(
				allErrs,
				validation.ValidateMachinePool(compute.Platform.OpenStack, ci, false, fldPath.Child("platform", "openstack"))...)
		}
	}

	// if areMachinePools; {validate platform.flavor()}

	return allErrs.ToAggregate()
}
