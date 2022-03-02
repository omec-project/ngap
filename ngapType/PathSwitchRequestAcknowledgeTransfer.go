// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/free5gc/aper" if it uses "aper"

type PathSwitchRequestAcknowledgeTransfer struct {
	ULNGUUPTNLInformation *UPTransportLayerInformation                                          `aper:"valueLB:0,valueUB:1,optional"`
	SecurityIndication    *SecurityIndication                                                   `aper:"valueExt,optional"`
	IEExtensions          *ProtocolExtensionContainerPathSwitchRequestAcknowledgeTransferExtIEs `aper:"optional"`
}
