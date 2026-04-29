// Copyright 2026 Forsway Solutions AB
//
// SPDX-License-Identifier: Apache-2.0

package ngap

import "github.com/omec-project/ngap/ngapType"

// HasUserLocationInformationNRExtension reports whether the given
// UserLocationInformationNR carries a protocol extension whose
// identifier matches the supplied id.
//
// NGAP shares a single integer numbering space between ProtocolIE-ID
// and ProtocolExtensionID, so any ProtocolIEID* constant from this
// package is a valid id regardless of whether the same identifier is
// also used as a ProtocolIE-ID in another context.
//
// A nil UserLocationInformationNR or a nil IEExtensions container both
// yield false.
func HasUserLocationInformationNRExtension(loc *ngapType.UserLocationInformationNR, id int64) bool {
	if loc == nil || loc.IEExtensions == nil {
		return false
	}
	for _, ext := range loc.IEExtensions.List {
		if ext.Id.Value == id {
			return true
		}
	}
	return false
}
