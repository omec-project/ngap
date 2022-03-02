// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/free5gc/aper" if it uses "aper"

const (
	UserLocationInformationPresentNothing int = iota /* No components present */
	UserLocationInformationPresentUserLocationInformationEUTRA
	UserLocationInformationPresentUserLocationInformationNR
	UserLocationInformationPresentUserLocationInformationN3IWF
	UserLocationInformationPresentChoiceExtensions
)

type UserLocationInformation struct {
	Present                      int
	UserLocationInformationEUTRA *UserLocationInformationEUTRA `aper:"valueExt"`
	UserLocationInformationNR    *UserLocationInformationNR    `aper:"valueExt"`
	UserLocationInformationN3IWF *UserLocationInformationN3IWF `aper:"valueExt"`
	ChoiceExtensions             *ProtocolIESingleContainerUserLocationInformationExtIEs
}
