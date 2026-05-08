// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"fmt"
	"strings"

	"github.com/omec-project/ngap/v2/aper"
	"github.com/omec-project/ngap/v2/logger"
	"github.com/omec-project/ngap/v2/ngapType"
	"github.com/omec-project/openapi/v2"
	"github.com/omec-project/openapi/v2/models"
)

func validBitString(bitString aper.BitString, bitLength int) bool {
	return bitString.BitLength == uint64(bitLength) && len(bitString.Bytes) == (bitLength+7)/8
}

func RanIdToModels(ranNodeId ngapType.GlobalRANNodeID) (ranId models.GlobalRanNodeId, err error) {
	switch ranNodeId.Present {
	case ngapType.GlobalRANNodeIDPresentGlobalGNBID:
		ranId.GNbId = models.NewGNbIdWithDefaults()
		gnbId := ranId.GNbId
		ngapGnbId := ranNodeId.GlobalGNBID
		plmnid, err := PlmnIdToModels(ngapGnbId.PLMNIdentity)
		if err != nil {
			return models.GlobalRanNodeId{}, fmt.Errorf("invalid GlobalGNBID PLMN identity: %w", err)
		}
		ranId.PlmnId = plmnid
		if ngapGnbId.GNBID.Present == ngapType.GNBIDPresentGNBID {
			choiceGnbId := ngapGnbId.GNBID.GNBID
			gnbId.BitLength = int32(choiceGnbId.BitLength)
			gnbId.GNBValue = BitStringToHex(choiceGnbId)
		}
	case ngapType.GlobalRANNodeIDPresentGlobalNgENBID:
		ngapNgENBID := ranNodeId.GlobalNgENBID
		plmnid, err := PlmnIdToModels(ngapNgENBID.PLMNIdentity)
		if err != nil {
			return models.GlobalRanNodeId{}, fmt.Errorf("invalid GlobalNgENBID PLMN identity: %w", err)
		}
		ranId.PlmnId = plmnid
		switch ngapNgENBID.NgENBID.Present {
		case ngapType.NgENBIDPresentMacroNgENBID:
			macroNgENBID := ngapNgENBID.NgENBID.MacroNgENBID
			ranId.NgeNbId = openapi.PtrString("MacroNGeNB-" + BitStringToHex(macroNgENBID))
		case ngapType.NgENBIDPresentShortMacroNgENBID:
			shortMacroNgENBID := ngapNgENBID.NgENBID.ShortMacroNgENBID
			ranId.NgeNbId = openapi.PtrString("SMacroNGeNB-" + BitStringToHex(shortMacroNgENBID))
		case ngapType.NgENBIDPresentLongMacroNgENBID:
			longMacroNgENBID := ngapNgENBID.NgENBID.LongMacroNgENBID
			ranId.NgeNbId = openapi.PtrString("LMacroNGeNB-" + BitStringToHex(longMacroNgENBID))
		default:
			return models.GlobalRanNodeId{}, fmt.Errorf("unsupported NgENBID present type %d", ngapNgENBID.NgENBID.Present)
		}
	case ngapType.GlobalRANNodeIDPresentGlobalN3IWFID:
		ngapN3IWFID := ranNodeId.GlobalN3IWFID
		plmnid, err := PlmnIdToModels(ngapN3IWFID.PLMNIdentity)
		if err != nil {
			return models.GlobalRanNodeId{}, fmt.Errorf("invalid GlobalN3IWFID PLMN identity: %w", err)
		}
		ranId.PlmnId = plmnid
		if ngapN3IWFID.N3IWFID.Present == ngapType.N3IWFIDPresentN3IWFID {
			choiceN3IWFID := ngapN3IWFID.N3IWFID.N3IWFID
			ranId.N3IwfId = openapi.PtrString(BitStringToHex(choiceN3IWFID))
		}
	default:
		return models.GlobalRanNodeId{}, fmt.Errorf("unsupported GlobalRANNodeID present type %d", ranNodeId.Present)
	}

	return ranId, nil
}

func RanIDToNgap(modelsRanNodeId models.GlobalRanNodeId) ngapType.GlobalRANNodeID {
	var ngapRanNodeId ngapType.GlobalRANNodeID

	if modelsRanNodeId.GNbId != nil && modelsRanNodeId.GNbId.BitLength != 0 {
		bitString := HexToBitString(modelsRanNodeId.GNbId.GNBValue, int(modelsRanNodeId.GNbId.BitLength))
		if validBitString(bitString, int(modelsRanNodeId.GNbId.BitLength)) {
			ngapRanNodeId.Present = ngapType.GlobalRANNodeIDPresentGlobalGNBID
			ngapRanNodeId.GlobalGNBID = new(ngapType.GlobalGNBID)
			globalGNBID := ngapRanNodeId.GlobalGNBID

			globalGNBID.PLMNIdentity = PlmnIdToNgap(modelsRanNodeId.PlmnId)
			globalGNBID.GNBID.Present = ngapType.GNBIDPresentGNBID
			globalGNBID.GNBID.GNBID = new(aper.BitString)
			*globalGNBID.GNBID.GNBID = bitString
		} else {
			logger.NgapLog.Warnf("RanIDToNgap: invalid GNBID payload %q with bit length %d", modelsRanNodeId.GNbId.GNBValue, modelsRanNodeId.GNbId.BitLength)
		}
	} else if modelsRanNodeId.GetNgeNbId() != "" {
		globalNgENBID := &ngapType.GlobalNgENBID{PLMNIdentity: PlmnIdToNgap(modelsRanNodeId.PlmnId)}
		ngENBID := &globalNgENBID.NgENBID
		ngeNbIdVal := modelsRanNodeId.GetNgeNbId()
		if strings.HasPrefix(ngeNbIdVal, "MacroNGeNB-") {
			bitString := HexToBitString(ngeNbIdVal[11:], 18)
			if validBitString(bitString, 18) {
				ngENBID.Present = ngapType.NgENBIDPresentMacroNgENBID
				ngENBID.MacroNgENBID = new(aper.BitString)
				*ngENBID.MacroNgENBID = bitString
				ngapRanNodeId.Present = ngapType.GlobalRANNodeIDPresentGlobalNgENBID
				ngapRanNodeId.GlobalNgENBID = globalNgENBID
			} else {
				logger.NgapLog.Warnf("RanIDToNgap: invalid MacroNgENBID payload %q", ngeNbIdVal)
			}
		} else if strings.HasPrefix(ngeNbIdVal, "SMacroNGeNB-") {
			bitString := HexToBitString(ngeNbIdVal[12:], 20)
			if validBitString(bitString, 20) {
				ngENBID.Present = ngapType.NgENBIDPresentShortMacroNgENBID
				ngENBID.ShortMacroNgENBID = new(aper.BitString)
				*ngENBID.ShortMacroNgENBID = bitString
				ngapRanNodeId.Present = ngapType.GlobalRANNodeIDPresentGlobalNgENBID
				ngapRanNodeId.GlobalNgENBID = globalNgENBID
			} else {
				logger.NgapLog.Warnf("RanIDToNgap: invalid ShortMacroNgENBID payload %q", ngeNbIdVal)
			}
		} else if strings.HasPrefix(ngeNbIdVal, "LMacroNGeNB-") {
			bitString := HexToBitString(ngeNbIdVal[12:], 21)
			if validBitString(bitString, 21) {
				ngENBID.Present = ngapType.NgENBIDPresentLongMacroNgENBID
				ngENBID.LongMacroNgENBID = new(aper.BitString)
				*ngENBID.LongMacroNgENBID = bitString
				ngapRanNodeId.Present = ngapType.GlobalRANNodeIDPresentGlobalNgENBID
				ngapRanNodeId.GlobalNgENBID = globalNgENBID
			} else {
				logger.NgapLog.Warnf("RanIDToNgap: invalid LongMacroNgENBID payload %q", ngeNbIdVal)
			}
		} else {
			logger.NgapLog.Warnf("RanIDToNgap: unexpected NgENBID format %q", ngeNbIdVal)
		}
	} else if modelsRanNodeId.GetN3IwfId() != "" {
		logger.NgapLog.Debugf("RanIDToNgap: converting GlobalN3IWFID %q", modelsRanNodeId.GetN3IwfId())
		bitLength := len(modelsRanNodeId.GetN3IwfId()) * 4
		bitString := HexToBitString(modelsRanNodeId.GetN3IwfId(), bitLength)
		if validBitString(bitString, bitLength) {
			ngapRanNodeId.Present = ngapType.GlobalRANNodeIDPresentGlobalN3IWFID
			ngapRanNodeId.GlobalN3IWFID = new(ngapType.GlobalN3IWFID)
			globalN3IWFID := ngapRanNodeId.GlobalN3IWFID

			globalN3IWFID.PLMNIdentity = PlmnIdToNgap(modelsRanNodeId.PlmnId)
			globalN3IWFID.N3IWFID.Present = ngapType.N3IWFIDPresentN3IWFID
			globalN3IWFID.N3IWFID.N3IWFID = new(aper.BitString)
			*globalN3IWFID.N3IWFID.N3IWFID = bitString
		} else {
			logger.NgapLog.Warnf("RanIDToNgap: invalid N3IWFID payload %q", modelsRanNodeId.GetN3IwfId())
		}
	} else {
		logger.NgapLog.Warnf("RanIDToNgap: no supported RAN node identifier found for PLMN %+v", modelsRanNodeId.PlmnId)
	}

	return ngapRanNodeId
}
