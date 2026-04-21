// Copyright 2026 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"strings"
	"testing"

	"github.com/omec-project/ngap/aper"
	"github.com/omec-project/ngap/ngapType"
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
			if snssai.Sst != tt.wantSst || snssai.Sd != tt.wantSd {
				t.Fatalf("expected S-NSSAI %d/%s, got %d/%s", tt.wantSst, tt.wantSd, snssai.Sst, snssai.Sd)
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
			if modelsTai.PlmnId == nil {
				t.Fatal("expected PlmnId to be set")
			}
			if modelsTai.PlmnId.Mcc != tt.wantMcc || modelsTai.PlmnId.Mnc != tt.wantMnc {
				t.Fatalf("expected MCC/MNC %s/%s, got %s/%s", tt.wantMcc, tt.wantMnc, modelsTai.PlmnId.Mcc, modelsTai.PlmnId.Mnc)
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
