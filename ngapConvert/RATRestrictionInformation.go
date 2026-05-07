// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"github.com/omec-project/ngap/v2/aper"
	"github.com/omec-project/ngap/v2/ngapType"
	"github.com/omec-project/openapi/v2/models"
)

// TS 38.413 9.3.1.85
func RATRestrictionInformationToNgap(ratType models.RatType) (ratResInfo ngapType.RATRestrictionInformation) {
	// Only support EUTRA & NR in version15.2.0
	switch ratType {
	case models.RATTYPE_EUTRA:
		ratResInfo.Value = aper.BitString{
			Bytes:     []byte{0x80},
			BitLength: 8,
		}
	case models.RATTYPE_NR:
		ratResInfo.Value = aper.BitString{
			Bytes:     []byte{0x40},
			BitLength: 8,
		}
	}
	return
}
