// Copyright 2026 Forsway Solutions AB
//
// SPDX-License-Identifier: Apache-2.0

package ngap_test

import (
	"testing"

	"github.com/omec-project/ngap"
	"github.com/omec-project/ngap/ngapType"
)

func TestHasUserLocationInformationNRExtension(t *testing.T) {
	tests := []struct {
		name string
		loc  *ngapType.UserLocationInformationNR
		id   int64
		want bool
	}{
		{
			name: "nil location",
			loc:  nil,
			id:   ngapType.ProtocolIEIDNRNTNTAIInformation,
			want: false,
		},
		{
			name: "nil extensions",
			loc:  &ngapType.UserLocationInformationNR{},
			id:   ngapType.ProtocolIEIDNRNTNTAIInformation,
			want: false,
		},
		{
			name: "empty extension list",
			loc: &ngapType.UserLocationInformationNR{
				IEExtensions: &ngapType.ProtocolExtensionContainerUserLocationInformationNRExtIEs{},
			},
			id:   ngapType.ProtocolIEIDNRNTNTAIInformation,
			want: false,
		},
		{
			name: "unrelated extension only",
			loc:  newNRLocationWithExtensionIDs(42),
			id:   ngapType.ProtocolIEIDNRNTNTAIInformation,
			want: false,
		},
		{
			name: "matching extension present",
			loc:  newNRLocationWithExtensionIDs(ngapType.ProtocolIEIDNRNTNTAIInformation),
			id:   ngapType.ProtocolIEIDNRNTNTAIInformation,
			want: true,
		},
		{
			name: "matching extension among several",
			loc:  newNRLocationWithExtensionIDs(42, 99, ngapType.ProtocolIEIDNRNTNTAIInformation, 150),
			id:   ngapType.ProtocolIEIDNRNTNTAIInformation,
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ngap.HasUserLocationInformationNRExtension(tc.loc, tc.id); got != tc.want {
				t.Errorf("HasUserLocationInformationNRExtension = %v, want %v", got, tc.want)
			}
		})
	}
}

func newNRLocationWithExtensionIDs(ids ...int64) *ngapType.UserLocationInformationNR {
	list := make([]ngapType.UserLocationInformationNRExtIEs, 0, len(ids))
	for _, id := range ids {
		list = append(list, ngapType.UserLocationInformationNRExtIEs{
			Id: ngapType.ProtocolExtensionID{Value: id},
		})
	}
	return &ngapType.UserLocationInformationNR{
		IEExtensions: &ngapType.ProtocolExtensionContainerUserLocationInformationNRExtIEs{
			List: list,
		},
	}
}
