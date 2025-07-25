// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"github.com/omec-project/aper"
	"github.com/omec-project/ngap/logger"
	"github.com/omec-project/ngap/ngapType"
	"github.com/omec-project/openapi/models"
)

func RanIdToModels(ranNodeId ngapType.GlobalRANNodeID) (ranId models.GlobalRanNodeId) {
	switch ranNodeId.Present {
	case ngapType.GlobalRANNodeIDPresentGlobalGNBID:
		ranId.GNbId = new(models.GNbId)
		gnbId := ranId.GNbId
		ngapGnbId := ranNodeId.GlobalGNBID
		plmnid := PlmnIdToModels(ngapGnbId.PLMNIdentity)
		ranId.PlmnId = &plmnid
		if ngapGnbId.GNBID.Present == ngapType.GNBIDPresentGNBID {
			choiceGnbId := ngapGnbId.GNBID.GNBID
			gnbId.BitLength = int32(choiceGnbId.BitLength)
			gnbId.GNBValue = BitStringToHex(choiceGnbId)
		}
	case ngapType.GlobalRANNodeIDPresentGlobalNgENBID:
		ngapNgENBID := ranNodeId.GlobalNgENBID
		plmnid := PlmnIdToModels(ngapNgENBID.PLMNIdentity)
		ranId.PlmnId = &plmnid
		switch ngapNgENBID.NgENBID.Present {
		case ngapType.NgENBIDPresentMacroNgENBID:
			macroNgENBID := ngapNgENBID.NgENBID.MacroNgENBID
			ranId.NgeNbId = "MacroNGeNB-" + BitStringToHex(macroNgENBID)
		case ngapType.NgENBIDPresentShortMacroNgENBID:
			shortMacroNgENBID := ngapNgENBID.NgENBID.ShortMacroNgENBID
			ranId.NgeNbId = "SMacroNGeNB-" + BitStringToHex(shortMacroNgENBID)
		case ngapType.NgENBIDPresentLongMacroNgENBID:
			longMacroNgENBID := ngapNgENBID.NgENBID.LongMacroNgENBID
			ranId.NgeNbId = "LMacroNGeNB-" + BitStringToHex(longMacroNgENBID)
		default:
			logger.NgapLog.Warnf("RanIdToModels: Unexpected NgENBID present type %d", ngapNgENBID.NgENBID.Present)
		}
	case ngapType.GlobalRANNodeIDPresentGlobalN3IWFID:
		ngapN3IWFID := ranNodeId.GlobalN3IWFID
		plmnid := PlmnIdToModels(ngapN3IWFID.PLMNIdentity)
		ranId.PlmnId = &plmnid
		if ngapN3IWFID.N3IWFID.Present == ngapType.N3IWFIDPresentN3IWFID {
			choiceN3IWFID := ngapN3IWFID.N3IWFID.N3IWFID
			ranId.N3IwfId = BitStringToHex(choiceN3IWFID)
		}
	default:
		logger.NgapLog.Warnf("RanIdToModels: Unexpected GlobalRANNodeID present type %d", ranNodeId.Present)
	}

	return ranId
}

func RanIDToNgap(modelsRanNodeId models.GlobalRanNodeId) ngapType.GlobalRANNodeID {
	var ngapRanNodeId ngapType.GlobalRANNodeID

	if modelsRanNodeId.GNbId.BitLength != 0 {
		ngapRanNodeId.Present = ngapType.GlobalRANNodeIDPresentGlobalGNBID
		ngapRanNodeId.GlobalGNBID = new(ngapType.GlobalGNBID)
		globalGNBID := ngapRanNodeId.GlobalGNBID

		globalGNBID.PLMNIdentity = PlmnIdToNgap(*modelsRanNodeId.PlmnId)
		globalGNBID.GNBID.Present = ngapType.GNBIDPresentGNBID
		globalGNBID.GNBID.GNBID = new(aper.BitString)
		*globalGNBID.GNBID.GNBID = HexToBitString(modelsRanNodeId.GNbId.GNBValue, int(modelsRanNodeId.GNbId.BitLength))
	} else if modelsRanNodeId.NgeNbId != "" {
		ngapRanNodeId.Present = ngapType.GlobalRANNodeIDPresentGlobalNgENBID
		ngapRanNodeId.GlobalNgENBID = new(ngapType.GlobalNgENBID)
		globalNgENBID := ngapRanNodeId.GlobalNgENBID

		globalNgENBID.PLMNIdentity = PlmnIdToNgap(*modelsRanNodeId.PlmnId)
		ngENBID := &globalNgENBID.NgENBID
		if modelsRanNodeId.NgeNbId[:11] == "MacroNGeNB-" {
			ngENBID.Present = ngapType.NgENBIDPresentMacroNgENBID
			ngENBID.MacroNgENBID = new(aper.BitString)
			*ngENBID.MacroNgENBID = HexToBitString(modelsRanNodeId.NgeNbId[11:], 18)
		} else if modelsRanNodeId.NgeNbId[:12] == "SMacroNGeNB-" {
			ngENBID.Present = ngapType.NgENBIDPresentShortMacroNgENBID
			ngENBID.ShortMacroNgENBID = new(aper.BitString)
			*ngENBID.ShortMacroNgENBID = HexToBitString(modelsRanNodeId.NgeNbId[12:], 20)
		} else if modelsRanNodeId.NgeNbId[:12] == "LMacroNGeNB-" {
			ngENBID.Present = ngapType.NgENBIDPresentLongMacroNgENBID
			ngENBID.LongMacroNgENBID = new(aper.BitString)
			*ngENBID.LongMacroNgENBID = HexToBitString(modelsRanNodeId.NgeNbId[12:], 21)
		}
	} else if modelsRanNodeId.N3IwfId != "" {
		ngapRanNodeId.Present = ngapType.GlobalRANNodeIDPresentGlobalN3IWFID
		ngapRanNodeId.GlobalN3IWFID = new(ngapType.GlobalN3IWFID)
		globalN3IWFID := ngapRanNodeId.GlobalN3IWFID

		globalN3IWFID.PLMNIdentity = PlmnIdToNgap(*modelsRanNodeId.PlmnId)
		globalN3IWFID.N3IWFID.Present = ngapType.N3IWFIDPresentN3IWFID
		globalN3IWFID.N3IWFID.N3IWFID = new(aper.BitString)
		*globalN3IWFID.N3IWFID.N3IWFID = HexToBitString(modelsRanNodeId.N3IwfId, len(modelsRanNodeId.N3IwfId)*4)
	}

	return ngapRanNodeId
}
