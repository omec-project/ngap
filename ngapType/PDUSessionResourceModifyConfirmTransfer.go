// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/omec-project/ngap/aper" if it uses "aper"

type PDUSessionResourceModifyConfirmTransfer struct {
	QosFlowModifyConfirmList      QosFlowModifyConfirmList
	ULNGUUPTNLInformation         UPTransportLayerInformation                                              `aper:"valueLB:0,valueUB:1"`
	AdditionalNGUUPTNLInformation *UPTransportLayerInformationPairList                                     `aper:"optional"`
	QosFlowFailedToModifyList     *QosFlowListWithCause                                                    `aper:"optional"`
	IEExtensions                  *ProtocolExtensionContainerPDUSessionResourceModifyConfirmTransferExtIEs `aper:"optional"`
}
