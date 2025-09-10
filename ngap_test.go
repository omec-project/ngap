// Copyright 2024 Canonical, Ltd.
//
// SPDX-License-Identifier: Apache-2.0

package ngap_test

import (
	"testing"

	"github.com/omec-project/ngap"
	"github.com/omec-project/ngap/aper"
	"github.com/omec-project/ngap/ngapType"
)

// globalRANIDIE for MCC 208, MNC 93
var globalRANIDIE ngapType.NGSetupRequestIEs = ngapType.NGSetupRequestIEs{
	Id:          ngapType.ProtocolIEID{Value: ngapType.ProtocolIEIDGlobalRANNodeID},
	Criticality: ngapType.Criticality{Value: ngapType.CriticalityPresentReject},
	Value: ngapType.NGSetupRequestIEsValue{
		Present: ngapType.GlobalRANNodeIDPresentGlobalGNBID,
		GlobalRANNodeID: &ngapType.GlobalRANNodeID{
			Present: ngapType.GlobalRANNodeIDPresentGlobalGNBID,
			GlobalGNBID: &ngapType.GlobalGNBID{
				PLMNIdentity: ngapType.PLMNIdentity{Value: aper.OctetString{0x02, 0xF8, 0x39}},
				GNBID: ngapType.GNBID{
					Present: ngapType.GNBIDPresentGNBID,
					GNBID: &aper.BitString{
						Bytes:     []byte{0x00, 0x01, 0x02},
						BitLength: 24,
					},
				},
			},
		},
	},
}

func TestSimplePDUEncoding(t *testing.T) {
	// NGSetupRequest PDU
	pdu := ngapType.NGAPPDU{
		Present: ngapType.NGAPPDUPresentInitiatingMessage,
		InitiatingMessage: &ngapType.InitiatingMessage{
			ProcedureCode: ngapType.ProcedureCode{Value: ngapType.ProcedureCodeNGSetup},
			Criticality:   ngapType.Criticality{Value: ngapType.CriticalityPresentReject},
			Value: ngapType.InitiatingMessageValue{
				Present: ngapType.InitiatingMessagePresentNGSetupRequest,
				NGSetupRequest: &ngapType.NGSetupRequest{
					ProtocolIEs: ngapType.ProtocolIEContainerNGSetupRequestIEs{
						List: []ngapType.NGSetupRequestIEs{
							globalRANIDIE,
						},
					},
				},
			},
		},
	}

	result, err := ngap.Encoder(pdu)
	if err != nil {
		t.Errorf("Could not encode simple PDU: %v; got error: %v\n", pdu, err)
	}

	expected := []byte{
		0x00, 0x15, 0x00, 0x0F, 0x00, 0x00, 0x01, 0x00, 0x1B, 0x00,
		0x08, 0x00, 0x02, 0xF8, 0x39, 0x10, 0x00, 0x01, 0x02,
	}

	if len(result) != len(expected) {
		t.Errorf("Got wrong result length: %v\n", len(result))
	}

	for i, b := range result {
		if b != expected[i] {
			t.Errorf("Byte %d was %v, expected %v\n", i, b, expected[i])
		}
	}
}

func TestSimplePDUDecoding(t *testing.T) {
	// ASN.1 PER encoded PDU
	pdu_bytes := []byte{
		0x00, 0x15, 0x00, 0x0F, 0x00, 0x00, 0x01, 0x00, 0x1B, 0x00,
		0x08, 0x00, 0x02, 0xF8, 0x39, 0x10, 0x00, 0x01, 0x02,
	}

	pdu, err := ngap.Decoder(pdu_bytes)
	if err != nil {
		t.Errorf("Could not decode simple PDU bytes: %v; got error: %v\n", pdu_bytes, err)
	}

	ie := pdu.InitiatingMessage.Value.NGSetupRequest.ProtocolIEs.List[0]
	plmn := ie.Value.GlobalRANNodeID.GlobalGNBID.PLMNIdentity.Value

	if len(plmn) != 3 || plmn[0] != 0x02 || plmn[1] != 0xF8 || plmn[2] != 0x39 {
		t.Errorf("Failed decoding simple PDU: %v; Expected PLMN: %v, got %v\n",
			pdu_bytes, []byte{0x02, 0xF8, 0x39}, plmn)
	}
}
