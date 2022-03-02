// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/free5gc/aper" if it uses "aper"

type IntendedNumberOfPagingAttempts struct {
	Value int64 `aper:"valueExt,valueLB:1,valueUB:16"`
}
