// Copyright 2026 Forsway Solutions AB
//
// SPDX-License-Identifier: Apache-2.0

package ngap_test

import (
	"testing"

	"github.com/omec-project/ngap/v2"
	"github.com/omec-project/ngap/v2/ngapType"
)

func TestFindUserLocationInformationNRNTNTAIInformation(t *testing.T) {
	expected := &ngapType.NRNTNTAIInformation{}
	other := &ngapType.NRNTNTAIInformation{}

	tests := []struct {
		name string
		loc  *ngapType.UserLocationInformationNR
		want *ngapType.NRNTNTAIInformation
	}{
		{
			name: "nil location",
			loc:  nil,
			want: nil,
		},
		{
			name: "nil extensions",
			loc:  &ngapType.UserLocationInformationNR{},
			want: nil,
		},
		{
			name: "empty extension list",
			loc: &ngapType.UserLocationInformationNR{
				IEExtensions: &ngapType.ProtocolExtensionContainerUserLocationInformationNRExtIEs{},
			},
			want: nil,
		},
		{
			name: "unrelated extension only",
			loc:  newNRLocationWithExtensions(userLocationInformationNRExtensionNGRANCGI()),
			want: nil,
		},
		{
			name: "matching extension present",
			loc:  newNRLocationWithExtensions(userLocationInformationNRExtensionNRNTNTAIInformation(expected)),
			want: expected,
		},
		{
			name: "matching extension among several",
			loc: newNRLocationWithExtensions(
				userLocationInformationNRExtensionNGRANCGI(),
				userLocationInformationNRExtensionNRNTNTAIInformation(expected),
				userLocationInformationNRExtensionNRNTNTAIInformation(other),
			),
			want: expected,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ngap.FindUserLocationInformationNRNTNTAIInformation(tc.loc); got != tc.want {
				t.Errorf("FindUserLocationInformationNRNTNTAIInformation() = %p, want %p", got, tc.want)
			}
			var got *ngapType.NRNTNTAIInformation
			if tc.loc != nil && tc.loc.IEExtensions != nil {
				got = tc.loc.IEExtensions.FindNRNTNTAIInformation()
			}
			if got != tc.want {
				t.Errorf("IEExtensions.FindNRNTNTAIInformation() = %p, want %p", got, tc.want)
			}
		})
	}
}

func newNRLocationWithExtensions(exts ...ngapType.UserLocationInformationNRExtIEs) *ngapType.UserLocationInformationNR {
	return &ngapType.UserLocationInformationNR{
		IEExtensions: &ngapType.ProtocolExtensionContainerUserLocationInformationNRExtIEs{
			List: exts,
		},
	}
}

func userLocationInformationNRExtensionNGRANCGI() ngapType.UserLocationInformationNRExtIEs {
	return ngapType.UserLocationInformationNRExtIEs{
		Id: ngapType.ProtocolExtensionID{Value: ngapType.ProtocolIEIDNGRANCGI},
		ExtensionValue: ngapType.UserLocationInformationNRExtIEsExtensionValue{
			Present:           ngapType.UserLocationInformationNRExtIEsPresentPSCellInformation,
			PSCellInformation: &ngapType.NGRANCGI{},
		},
	}
}

func userLocationInformationNRExtensionNRNTNTAIInformation(value *ngapType.NRNTNTAIInformation) ngapType.UserLocationInformationNRExtIEs {
	return ngapType.UserLocationInformationNRExtIEs{
		Id: ngapType.ProtocolExtensionID{Value: ngapType.ProtocolIEIDNRNTNTAIInformation},
		ExtensionValue: ngapType.UserLocationInformationNRExtIEsExtensionValue{
			Present:             ngapType.UserLocationInformationNRExtIEsPresentNRNTNTAIInformation,
			NRNTNTAIInformation: value,
		},
	}
}
