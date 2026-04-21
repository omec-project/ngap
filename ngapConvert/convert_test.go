// Copyright 2026 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"testing"

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
