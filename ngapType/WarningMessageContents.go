// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

import "github.com/free5gc/aper"

// Need to import "github.com/free5gc/aper" if it uses "aper"

type WarningMessageContents struct {
	Value aper.OctetString `aper:"sizeLB:1,sizeUB:9600"`
}
