// Copyright (C) 2026 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package ngap_test

import (
	"bytes"
	"testing"

	"github.com/omec-project/ngap/v2/aper"
	"github.com/omec-project/ngap/v2/ngapType"
)

// Extension additions such as INTEGER (1..65535, ..., 65536..109999) must not widen the root range.
func TestExtensibleIntegerRootRange(t *testing.T) {
	tests := []struct {
		name string
		in   any
		out  any
		want []byte
	}{
		{
			"ExtendedPacketDelayBudget root",
			ngapType.ExtendedPacketDelayBudget{Value: 100},
			&ngapType.ExtendedPacketDelayBudget{},
			[]byte{0x00, 0x00, 0x63},
		},
		{
			"ExtendedPacketDelayBudget extension",
			ngapType.ExtendedPacketDelayBudget{Value: 100000},
			&ngapType.ExtendedPacketDelayBudget{},
			[]byte{0x80, 0x03, 0x01, 0x86, 0xA0},
		},
		{
			"MaximumDataBurstVolume root",
			ngapType.MaximumDataBurstVolume{Value: 1000},
			&ngapType.MaximumDataBurstVolume{},
			[]byte{0x00, 0x03, 0xE8},
		},
		{
			"MaximumDataBurstVolume extension",
			ngapType.MaximumDataBurstVolume{Value: 2000000},
			&ngapType.MaximumDataBurstVolume{},
			[]byte{0x80, 0x03, 0x1E, 0x84, 0x80},
		},
	}
	for _, tc := range tests {
		got, err := aper.MarshalWithParams(tc.in, "")
		if err != nil {
			t.Fatalf("%s: marshal: %v", tc.name, err)
		}
		if !bytes.Equal(got, tc.want) {
			t.Errorf("%s: marshal = %x, want %x", tc.name, got, tc.want)
		}
		if err := aper.Unmarshal(tc.want, tc.out); err != nil {
			t.Fatalf("%s: unmarshal: %v", tc.name, err)
		}
	}
}
