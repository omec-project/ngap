// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/omec-project/ngap/aper" if it uses "aper"

const (
	AMFPagingTargetPresentNothing int = iota /* No components present */
	AMFPagingTargetPresentGlobalRANNodeID
	AMFPagingTargetPresentTAI
	AMFPagingTargetPresentChoiceExtensions
)

type AMFPagingTarget struct {
	Present          int
	GlobalRANNodeID  *GlobalRANNodeID `aper:"valueLB:0,valueUB:3"`
	TAI              *TAI             `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerAMFPagingTargetExtIEs
}
