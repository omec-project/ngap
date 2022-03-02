// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/free5gc/aper" if it uses "aper"

const (
	TargetIDPresentNothing int = iota /* No components present */
	TargetIDPresentTargetRANNodeID
	TargetIDPresentTargeteNBID
	TargetIDPresentChoiceExtensions
)

type TargetID struct {
	Present          int
	TargetRANNodeID  *TargetRANNodeID `aper:"valueExt"`
	TargeteNBID      *TargeteNBID     `aper:"valueExt"`
	ChoiceExtensions *ProtocolIESingleContainerTargetIDExtIEs
}
