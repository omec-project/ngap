// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

import "github.com/omec-project/ngap/aper"

// Need to import "github.com/omec-project/ngap/aper" if it uses "aper"

const (
	NotificationCausePresentFulfilled    aper.Enumerated = 0
	NotificationCausePresentNotFulfilled aper.Enumerated = 1
)

type NotificationCause struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
