// Copyright 2026 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"strings"
	"testing"

	"github.com/omec-project/ngap/v2/aper"
	"github.com/omec-project/ngap/v2/ngapType"
	"github.com/omec-project/openapi/v2"
	"github.com/omec-project/openapi/v2/models"
)

func TestPlmnIdToModels(t *testing.T) {
	tests := []struct {
		name      string
		value     []byte
		wantMcc   string
		wantMnc   string
		wantError bool
	}{
		{name: "valid two digit mnc", value: []byte{0x02, 0xf8, 0x39}, wantMcc: "208", wantMnc: "93"},
		{name: "valid three digit mnc", value: []byte{0x13, 0x20, 0x06}, wantMcc: "310", wantMnc: "260"},
		{name: "short value", value: []byte{0x02, 0xf8}, wantError: true},
		{name: "invalid mcc digit", value: []byte{0xa2, 0xf8, 0x39}, wantError: true},
		{name: "invalid mnc digit", value: []byte{0x02, 0xf8, 0xa9}, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plmnID, err := PlmnIdToModels(ngapType.PLMNIdentity{Value: tt.value})
			if tt.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if plmnID.Mcc != tt.wantMcc || plmnID.Mnc != tt.wantMnc {
				t.Fatalf("expected MCC/MNC %s/%s, got %s/%s", tt.wantMcc, tt.wantMnc, plmnID.Mcc, plmnID.Mnc)
			}
		})
	}
}

func TestSNssaiToModels(t *testing.T) {
	tests := []struct {
		name      string
		snssai    ngapType.SNSSAI
		wantSst   int32
		wantSd    string
		wantError bool
	}{
		{name: "sst only", snssai: ngapType.SNSSAI{SST: ngapType.SST{Value: []byte{0x01}}}, wantSst: 1},
		{name: "sst and sd", snssai: ngapType.SNSSAI{SST: ngapType.SST{Value: []byte{0x02}}, SD: &ngapType.SD{Value: []byte{0x11, 0x22, 0x33}}}, wantSst: 2, wantSd: "112233"},
		{name: "missing sst", snssai: ngapType.SNSSAI{}, wantError: true},
		{name: "invalid sst length", snssai: ngapType.SNSSAI{SST: ngapType.SST{Value: []byte{0x01, 0x02}}}, wantError: true},
		{name: "invalid sd length", snssai: ngapType.SNSSAI{SST: ngapType.SST{Value: []byte{0x01}}, SD: &ngapType.SD{Value: []byte{0xaa}}}, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snssai, err := SNssaiToModels(tt.snssai)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if snssai.GetSst() != tt.wantSst || snssai.GetSd() != tt.wantSd {
				t.Fatalf("expected S-NSSAI %d/%s, got %d/%s", tt.wantSst, tt.wantSd, snssai.GetSst(), snssai.GetSd())
			}
		})
	}
}

func TestTaiToModels(t *testing.T) {
	tests := []struct {
		name      string
		tai       ngapType.TAI
		wantTac   string
		wantMcc   string
		wantMnc   string
		wantError string
	}{
		{
			name: "valid tai",
			tai: ngapType.TAI{
				PLMNIdentity: ngapType.PLMNIdentity{Value: []byte{0x02, 0xf8, 0x39}},
				TAC:          ngapType.TAC{Value: []byte{0x00, 0x00, 0x07}},
			},
			wantTac: "000007",
			wantMcc: "208",
			wantMnc: "93",
		},
		{
			name: "invalid plmn propagates error",
			tai: ngapType.TAI{
				PLMNIdentity: ngapType.PLMNIdentity{Value: []byte{0x02, 0xf8}},
				TAC:          ngapType.TAC{Value: []byte{0x00, 0x00, 0x07}},
			},
			wantError: "invalid TAI PLMN identity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelsTai, err := TaiToModels(tt.tai)
			if tt.wantError != "" {
				if err == nil {
					t.Fatal("expected error")
				}
				if !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("expected error containing %q, got %v", tt.wantError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if modelsTai.PlmnId.GetMcc() != tt.wantMcc || modelsTai.PlmnId.GetMnc() != tt.wantMnc {
				t.Fatalf("expected MCC/MNC %s/%s, got %s/%s", tt.wantMcc, tt.wantMnc, modelsTai.PlmnId.GetMcc(), modelsTai.PlmnId.GetMnc())
			}
			if modelsTai.Tac != tt.wantTac {
				t.Fatalf("expected TAC %s, got %s", tt.wantTac, modelsTai.Tac)
			}
		})
	}
}

func TestRanIdToModels(t *testing.T) {
	tests := []struct {
		name      string
		ranNodeId ngapType.GlobalRANNodeID
		wantError string
	}{
		{
			name: "invalid gnb plmn propagates error",
			ranNodeId: ngapType.GlobalRANNodeID{
				Present: ngapType.GlobalRANNodeIDPresentGlobalGNBID,
				GlobalGNBID: &ngapType.GlobalGNBID{
					PLMNIdentity: ngapType.PLMNIdentity{Value: []byte{0x02, 0xf8}},
					GNBID: ngapType.GNBID{
						Present: ngapType.GNBIDPresentGNBID,
						GNBID:   &aper.BitString{Bytes: []byte{0x01}, BitLength: 8},
					},
				},
			},
			wantError: "invalid GlobalGNBID PLMN identity",
		},
		{
			name: "invalid ngenb plmn propagates error",
			ranNodeId: ngapType.GlobalRANNodeID{
				Present: ngapType.GlobalRANNodeIDPresentGlobalNgENBID,
				GlobalNgENBID: &ngapType.GlobalNgENBID{
					PLMNIdentity: ngapType.PLMNIdentity{Value: []byte{0x02, 0xf8}},
					NgENBID: ngapType.NgENBID{
						Present:      ngapType.NgENBIDPresentMacroNgENBID,
						MacroNgENBID: &aper.BitString{Bytes: []byte{0x01, 0x02, 0x03}, BitLength: 18},
					},
				},
			},
			wantError: "invalid GlobalNgENBID PLMN identity",
		},
		{
			name: "unsupported ngenb variant returns error",
			ranNodeId: ngapType.GlobalRANNodeID{
				Present: ngapType.GlobalRANNodeIDPresentGlobalNgENBID,
				GlobalNgENBID: &ngapType.GlobalNgENBID{
					PLMNIdentity: ngapType.PLMNIdentity{Value: []byte{0x02, 0xf8, 0x39}},
					NgENBID: ngapType.NgENBID{
						Present: 99,
					},
				},
			},
			wantError: "unsupported NgENBID present type",
		},
		{
			name: "invalid n3iwf plmn propagates error",
			ranNodeId: ngapType.GlobalRANNodeID{
				Present: ngapType.GlobalRANNodeIDPresentGlobalN3IWFID,
				GlobalN3IWFID: &ngapType.GlobalN3IWFID{
					PLMNIdentity: ngapType.PLMNIdentity{Value: []byte{0x02, 0xf8}},
					N3IWFID: ngapType.N3IWFID{
						Present: ngapType.N3IWFIDPresentN3IWFID,
						N3IWFID: &aper.BitString{Bytes: []byte{0x01}, BitLength: 8},
					},
				},
			},
			wantError: "invalid GlobalN3IWFID PLMN identity",
		},
		{
			name: "unsupported top-level variant returns error",
			ranNodeId: ngapType.GlobalRANNodeID{
				Present: 99,
			},
			wantError: "unsupported GlobalRANNodeID present type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RanIdToModels(tt.ranNodeId)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("expected error containing %q, got %v", tt.wantError, err)
			}
		})
	}
}

func TestRanIDToNgap(t *testing.T) {
	tests := []struct {
		name  string
		ranID models.GlobalRanNodeId
		check func(*testing.T, ngapType.GlobalRANNodeID)
	}{
		{
			name: "gnb",
			ranID: models.GlobalRanNodeId{
				PlmnId: models.PlmnId{Mcc: "208", Mnc: "93"},
				GNbId:  models.NewGNbId(22, "123456"),
			},
			check: func(t *testing.T, got ngapType.GlobalRANNodeID) {
				t.Helper()
				if got.Present != ngapType.GlobalRANNodeIDPresentGlobalGNBID {
					t.Fatalf("expected GlobalGNBID present, got %d", got.Present)
				}
				if got.GlobalGNBID == nil || got.GlobalGNBID.GNBID.GNBID == nil {
					t.Fatal("expected GlobalGNBID to be populated")
				}
				if got.GlobalGNBID.GNBID.Present != ngapType.GNBIDPresentGNBID {
					t.Fatalf("expected GNBID present, got %d", got.GlobalGNBID.GNBID.Present)
				}
				if got.GlobalGNBID.GNBID.GNBID.BitLength != 22 {
					t.Fatalf("expected GNB bit length 22, got %d", got.GlobalGNBID.GNBID.GNBID.BitLength)
				}
			},
		},
		{
			name: "macro ngenb",
			ranID: models.GlobalRanNodeId{
				PlmnId:  models.PlmnId{Mcc: "208", Mnc: "93"},
				NgeNbId: openapi.PtrString("MacroNGeNB-12345"),
			},
			check: func(t *testing.T, got ngapType.GlobalRANNodeID) {
				t.Helper()
				if got.Present != ngapType.GlobalRANNodeIDPresentGlobalNgENBID {
					t.Fatalf("expected GlobalNgENBID present, got %d", got.Present)
				}
				if got.GlobalNgENBID == nil || got.GlobalNgENBID.NgENBID.MacroNgENBID == nil {
					t.Fatal("expected macro NgENBID to be populated")
				}
				if got.GlobalNgENBID.NgENBID.Present != ngapType.NgENBIDPresentMacroNgENBID {
					t.Fatalf("expected macro NgENBID present, got %d", got.GlobalNgENBID.NgENBID.Present)
				}
				if got.GlobalNgENBID.NgENBID.MacroNgENBID.BitLength != 18 {
					t.Fatalf("expected macro NgENBID bit length 18, got %d", got.GlobalNgENBID.NgENBID.MacroNgENBID.BitLength)
				}
			},
		},
		{
			name: "short macro ngenb",
			ranID: models.GlobalRanNodeId{
				PlmnId:  models.PlmnId{Mcc: "208", Mnc: "93"},
				NgeNbId: openapi.PtrString("SMacroNGeNB-12345"),
			},
			check: func(t *testing.T, got ngapType.GlobalRANNodeID) {
				t.Helper()
				if got.Present != ngapType.GlobalRANNodeIDPresentGlobalNgENBID {
					t.Fatalf("expected GlobalNgENBID present, got %d", got.Present)
				}
				if got.GlobalNgENBID == nil || got.GlobalNgENBID.NgENBID.ShortMacroNgENBID == nil {
					t.Fatal("expected short macro NgENBID to be populated")
				}
				if got.GlobalNgENBID.NgENBID.Present != ngapType.NgENBIDPresentShortMacroNgENBID {
					t.Fatalf("expected short macro NgENBID present, got %d", got.GlobalNgENBID.NgENBID.Present)
				}
				if got.GlobalNgENBID.NgENBID.ShortMacroNgENBID.BitLength != 20 {
					t.Fatalf("expected short macro NgENBID bit length 20, got %d", got.GlobalNgENBID.NgENBID.ShortMacroNgENBID.BitLength)
				}
			},
		},
		{
			name: "long macro ngenb",
			ranID: models.GlobalRanNodeId{
				PlmnId:  models.PlmnId{Mcc: "208", Mnc: "93"},
				NgeNbId: openapi.PtrString("LMacroNGeNB-123456"),
			},
			check: func(t *testing.T, got ngapType.GlobalRANNodeID) {
				t.Helper()
				if got.Present != ngapType.GlobalRANNodeIDPresentGlobalNgENBID {
					t.Fatalf("expected GlobalNgENBID present, got %d", got.Present)
				}
				if got.GlobalNgENBID == nil || got.GlobalNgENBID.NgENBID.LongMacroNgENBID == nil {
					t.Fatal("expected long macro NgENBID to be populated")
				}
				if got.GlobalNgENBID.NgENBID.Present != ngapType.NgENBIDPresentLongMacroNgENBID {
					t.Fatalf("expected long macro NgENBID present, got %d", got.GlobalNgENBID.NgENBID.Present)
				}
				if got.GlobalNgENBID.NgENBID.LongMacroNgENBID.BitLength != 21 {
					t.Fatalf("expected long macro NgENBID bit length 21, got %d", got.GlobalNgENBID.NgENBID.LongMacroNgENBID.BitLength)
				}
			},
		},
		{
			name: "n3iwf",
			ranID: models.GlobalRanNodeId{
				PlmnId:  models.PlmnId{Mcc: "208", Mnc: "93"},
				N3IwfId: openapi.PtrString("1234abcd"),
			},
			check: func(t *testing.T, got ngapType.GlobalRANNodeID) {
				t.Helper()
				if got.Present != ngapType.GlobalRANNodeIDPresentGlobalN3IWFID {
					t.Fatalf("expected GlobalN3IWFID present, got %d", got.Present)
				}
				if got.GlobalN3IWFID == nil || got.GlobalN3IWFID.N3IWFID.N3IWFID == nil {
					t.Fatal("expected N3IWFID to be populated")
				}
				if got.GlobalN3IWFID.N3IWFID.Present != ngapType.N3IWFIDPresentN3IWFID {
					t.Fatalf("expected N3IWFID present, got %d", got.GlobalN3IWFID.N3IWFID.Present)
				}
				if got.GlobalN3IWFID.N3IWFID.N3IWFID.BitLength != 32 {
					t.Fatalf("expected N3IWFID bit length 32, got %d", got.GlobalN3IWFID.N3IWFID.N3IWFID.BitLength)
				}
			},
		},
		{
			name: "invalid gnb hex returns zero value",
			ranID: models.GlobalRanNodeId{
				PlmnId: models.PlmnId{Mcc: "208", Mnc: "93"},
				GNbId:  models.NewGNbId(22, "zzzzzz"),
			},
			check: func(t *testing.T, got ngapType.GlobalRANNodeID) {
				t.Helper()
				if got.Present != 0 {
					t.Fatalf("expected zero-value GlobalRANNodeID, got present %d", got.Present)
				}
				if got.GlobalGNBID != nil {
					t.Fatal("expected GlobalGNBID to remain nil for invalid GNB payload")
				}
			},
		},
		{
			name: "invalid ngenb format returns zero value",
			ranID: models.GlobalRanNodeId{
				PlmnId:  models.PlmnId{Mcc: "208", Mnc: "93"},
				NgeNbId: openapi.PtrString("Unexpected-12345"),
			},
			check: func(t *testing.T, got ngapType.GlobalRANNodeID) {
				t.Helper()
				if got.Present != 0 {
					t.Fatalf("expected zero-value GlobalRANNodeID, got present %d", got.Present)
				}
				if got.GlobalNgENBID != nil {
					t.Fatal("expected GlobalNgENBID to remain nil for invalid format")
				}
			},
		},
		{
			name: "invalid ngenb hex returns zero value",
			ranID: models.GlobalRanNodeId{
				PlmnId:  models.PlmnId{Mcc: "208", Mnc: "93"},
				NgeNbId: openapi.PtrString("MacroNGeNB-zzzzz"),
			},
			check: func(t *testing.T, got ngapType.GlobalRANNodeID) {
				t.Helper()
				if got.Present != 0 {
					t.Fatalf("expected zero-value GlobalRANNodeID, got present %d", got.Present)
				}
				if got.GlobalNgENBID != nil {
					t.Fatal("expected GlobalNgENBID to remain nil for invalid hex payload")
				}
			},
		},
		{
			name: "invalid n3iwf hex returns zero value",
			ranID: models.GlobalRanNodeId{
				PlmnId:  models.PlmnId{Mcc: "208", Mnc: "93"},
				N3IwfId: openapi.PtrString("zzzz"),
			},
			check: func(t *testing.T, got ngapType.GlobalRANNodeID) {
				t.Helper()
				if got.Present != 0 {
					t.Fatalf("expected zero-value GlobalRANNodeID, got present %d", got.Present)
				}
				if got.GlobalN3IWFID != nil {
					t.Fatal("expected GlobalN3IWFID to remain nil for invalid N3IWF payload")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, RanIDToNgap(tt.ranID))
		})
	}
}

func TestTraceDataToNgap(t *testing.T) {
	tests := []struct {
		name          string
		interfaceList *string
		wantBitLength uint64
		wantBytes     []byte
	}{
		{
			name:          "unset interface list remains unset",
			wantBitLength: 8,
			wantBytes:     []byte{0x00},
		},
		{
			name:          "valid interface list encodes one octet",
			interfaceList: openapi.PtrString("01"),
			wantBitLength: 8,
			wantBytes:     []byte{0x01},
		},
		{
			name:          "invalid interface list remains unset",
			interfaceList: openapi.PtrString("abcd"),
			wantBitLength: 8,
			wantBytes:     []byte{0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traceData := *models.NewTraceData("20893-a1b2c3", models.TRACEDEPTH_MINIMUM, "01", "01")
			traceData.InterfaceList = tt.interfaceList

			got := TraceDataToNgap(traceData, "abcd")
			if got.InterfacesToTrace.Value.BitLength != tt.wantBitLength {
				t.Fatalf("expected interface bit length %d, got %d", tt.wantBitLength, got.InterfacesToTrace.Value.BitLength)
			}
			if string(got.InterfacesToTrace.Value.Bytes) != string(tt.wantBytes) {
				t.Fatalf("expected interface bytes %v, got %v", tt.wantBytes, got.InterfacesToTrace.Value.Bytes)
			}
		})
	}
}
