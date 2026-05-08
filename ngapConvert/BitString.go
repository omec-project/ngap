// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"encoding/hex"

	"github.com/omec-project/ngap/v2/aper"
	"github.com/omec-project/ngap/v2/logger"
)

func BitStringToHex(bitString *aper.BitString) (hexString string) {
	hexString = hex.EncodeToString(bitString.Bytes)
	hexLen := (bitString.BitLength + 3) / 4
	hexString = hexString[:hexLen]
	return
}

func HexToBitString(hexString string, bitLength int) (bitString aper.BitString) {
	if bitLength <= 0 {
		return bitString
	}

	hexLen := len(hexString)
	if hexLen != (bitLength+3)/4 {
		logger.NgapLog.Warnf("hexLen[%d] doesn't match bitLength[%d]", hexLen, bitLength)
		return bitString
	}
	if hexLen%2 == 1 {
		hexString += "0"
	}
	if byteTmp, err := hex.DecodeString(hexString); err != nil {
		logger.NgapLog.Warnf("Decode byteString failed: %+v", err)
		return bitString
	} else {
		bitString.Bytes = byteTmp
	}
	if len(bitString.Bytes) != (bitLength+7)/8 {
		logger.NgapLog.Warnf("decoded byte length[%d] doesn't match bitLength[%d]", len(bitString.Bytes), bitLength)
		bitString.Bytes = nil
		return bitString
	}
	bitString.BitLength = uint64(bitLength)
	if bitLength%8 != 0 {
		mask := byte(0xff)
		mask = mask << uint(8-bitLength%8)
		bitString.Bytes[len(bitString.Bytes)-1] &= mask
	}
	return bitString
}

func ByteToBitString(byteArray []byte, bitLength int) (bitString aper.BitString) {
	byteLen := (bitLength + 7) / 8
	if byteLen > len(byteArray) {
		logger.NgapLog.Warnf("bitLength[%d] is beyond byteArray size[%d]", bitLength, len(byteArray))
		return
	}
	bitString.Bytes = byteArray
	bitString.BitLength = uint64(bitLength)
	return
}
