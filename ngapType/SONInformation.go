// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/omec-project/ngap/aper" if it uses "aper"

const (
	SONInformationPresentNothing int = iota /* No components present */
	SONInformationPresentSONInformationRequest
	SONInformationPresentSONInformationReply
	SONInformationPresentChoiceExtensions
)

type SONInformation struct {
	Present               int
	SONInformationRequest *SONInformationRequest
	SONInformationReply   *SONInformationReply `aper:"valueExt"`
	ChoiceExtensions      *ProtocolIESingleContainerSONInformationExtIEs
}
