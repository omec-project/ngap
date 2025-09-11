// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

import "github.com/omec-project/ngap/aper"

// Need to import "github.com/omec-project/ngap/aper" if it uses "aper"

const (
	ExpectedHOIntervalPresentSec15    aper.Enumerated = 0
	ExpectedHOIntervalPresentSec30    aper.Enumerated = 1
	ExpectedHOIntervalPresentSec60    aper.Enumerated = 2
	ExpectedHOIntervalPresentSec90    aper.Enumerated = 3
	ExpectedHOIntervalPresentSec120   aper.Enumerated = 4
	ExpectedHOIntervalPresentSec180   aper.Enumerated = 5
	ExpectedHOIntervalPresentLongTime aper.Enumerated = 6
)

type ExpectedHOInterval struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:6"`
}
