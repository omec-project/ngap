// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

import "github.com/omec-project/ngap/aper"

// Need to import "github.com/omec-project/ngap/aper" if it uses "aper"

type PDUSessionResourceHandoverItem struct {
	PDUSessionID            PDUSessionID
	HandoverCommandTransfer aper.OctetString
	IEExtensions            *ProtocolExtensionContainerPDUSessionResourceHandoverItemExtIEs `aper:"optional"`
}
