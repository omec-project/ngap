// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"encoding/hex"
	"fmt"

	"github.com/omec-project/ngap/v2/logger"
	"github.com/omec-project/ngap/v2/ngapType"
	"github.com/omec-project/openapi/v2"
	"github.com/omec-project/openapi/v2/models"
)

func SNssaiToModels(ngapSnssai ngapType.SNSSAI) (modelsSnssai models.Snssai, err error) {
	if len(ngapSnssai.SST.Value) != 1 {
		return models.Snssai{}, fmt.Errorf("invalid S-NSSAI SST length %d", len(ngapSnssai.SST.Value))
	}
	modelsSnssai.Sst = int32(ngapSnssai.SST.Value[0])
	if ngapSnssai.SD != nil {
		if len(ngapSnssai.SD.Value) != 3 {
			return models.Snssai{}, fmt.Errorf("invalid S-NSSAI SD length %d", len(ngapSnssai.SD.Value))
		}
		modelsSnssai.Sd = openapi.PtrString(hex.EncodeToString(ngapSnssai.SD.Value))
	}
	return modelsSnssai, nil
}

func SNssaiToNgap(modelsSnssai models.Snssai) ngapType.SNSSAI {
	var ngapSnssai ngapType.SNSSAI
	ngapSnssai.SST.Value = []byte{byte(modelsSnssai.Sst)}

	if modelsSnssai.GetSd() != "" {
		ngapSnssai.SD = new(ngapType.SD)
		if sdTmp, err := hex.DecodeString(modelsSnssai.GetSd()); err != nil {
			logger.NgapLog.Warnf("decode snssai.sd failed: %+v", err)
		} else {
			ngapSnssai.SD.Value = sdTmp
		}
	}
	return ngapSnssai
}
