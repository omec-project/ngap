// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"encoding/hex"
	"fmt"

	"github.com/omec-project/ngap/ngapType"
	"github.com/omec-project/openapi/models"
)

func TaiToModels(tai ngapType.TAI) (models.Tai, error) {
	var modelsTai models.Tai

	plmnID, err := PlmnIdToModels(tai.PLMNIdentity)
	if err != nil {
		return models.Tai{}, fmt.Errorf("invalid TAI PLMN identity: %w", err)
	}
	modelsTai.PlmnId = &plmnID
	modelsTai.Tac = hex.EncodeToString(tai.TAC.Value)

	return modelsTai, nil
}
