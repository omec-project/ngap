// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/omec-project/ngap/logger"
	"github.com/omec-project/ngap/ngapType"
	"github.com/omec-project/openapi/models"
)

func PlmnIdToModels(ngapPlmnId ngapType.PLMNIdentity) (modelsPlmnid models.PlmnId, err error) {
	value := ngapPlmnId.Value
	if len(value) != 3 {
		return models.PlmnId{}, fmt.Errorf("invalid PLMNIdentity length %d", len(value))
	}

	hexString := hex.EncodeToString(value)
	if len(hexString) != 6 {
		return models.PlmnId{}, fmt.Errorf("invalid encoded PLMNIdentity length %d", len(hexString))
	}

	mcc := []byte{hexString[1], hexString[0], hexString[3]}
	for _, digit := range mcc {
		if digit < '0' || digit > '9' {
			return models.PlmnId{}, fmt.Errorf("invalid MCC digit %q in PLMNIdentity %q", digit, hexString)
		}
	}
	modelsPlmnid.Mcc = string(mcc)

	if hexString[2] == 'f' {
		if hexString[5] < '0' || hexString[5] > '9' || hexString[4] < '0' || hexString[4] > '9' {
			return models.PlmnId{}, fmt.Errorf("invalid two-digit MNC in PLMNIdentity %q", hexString)
		}
		modelsPlmnid.Mnc = string([]byte{hexString[5], hexString[4]})
	} else {
		mnc := []byte{hexString[2], hexString[5], hexString[4]}
		for _, digit := range mnc {
			if digit < '0' || digit > '9' {
				return models.PlmnId{}, fmt.Errorf("invalid three-digit MNC in PLMNIdentity %q", hexString)
			}
		}
		modelsPlmnid.Mnc = string(mnc)
	}
	return modelsPlmnid, nil
}

func PlmnIdToNgap(modelsPlmnid models.PlmnId) ngapType.PLMNIdentity {
	var hexString string
	mcc := strings.Split(modelsPlmnid.Mcc, "")
	mnc := strings.Split(modelsPlmnid.Mnc, "")
	if len(modelsPlmnid.Mnc) == 2 {
		hexString = mcc[1] + mcc[0] + "f" + mcc[2] + mnc[1] + mnc[0]
	} else {
		hexString = mcc[1] + mcc[0] + mnc[0] + mcc[2] + mnc[2] + mnc[1]
	}

	var ngapPlmnId ngapType.PLMNIdentity
	if plmnId, err := hex.DecodeString(hexString); err != nil {
		logger.NgapLog.Warnf("decode plmn failed: %+v", err)
	} else {
		ngapPlmnId.Value = plmnId
	}
	return ngapPlmnId
}
