// Copyright 2026 Forsway Solutions AB
//
// SPDX-License-Identifier: Apache-2.0

package ngap

import "github.com/omec-project/ngap/v2/ngapType"

// FindUserLocationInformationNRNTNTAIInformation returns the decoded
// NRNTNTAIInformation protocol extension carried by loc, if present.
//
// A nil UserLocationInformationNR, a nil IEExtensions container, or a
// location that does not carry the typed extension yields nil.
func FindUserLocationInformationNRNTNTAIInformation(loc *ngapType.UserLocationInformationNR) *ngapType.NRNTNTAIInformation {
	if loc == nil || loc.IEExtensions == nil {
		return nil
	}
	return loc.IEExtensions.FindNRNTNTAIInformation()
}
