// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/omec-project/ngap/aper" if it uses "aper"

type SONInformationReply struct {
	XnTNLConfigurationInfo *XnTNLConfigurationInfo                              `aper:"valueExt,optional"`
	IEExtensions           *ProtocolExtensionContainerSONInformationReplyExtIEs `aper:"optional"`
}
