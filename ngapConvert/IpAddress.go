// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapConvert

import (
	"net"

	"github.com/omec-project/aper"
	"github.com/omec-project/ngap/logger"
	"github.com/omec-project/ngap/ngapType"
)

func IPAddressToString(ipAddr ngapType.TransportLayerAddress) (ipv4Addr, ipv6Addr string) {
	ip := ipAddr.Value

	// Described in 38.414
	switch ip.BitLength {
	case 32: // ipv4
		netIP := net.IPv4(ip.Bytes[0], ip.Bytes[1], ip.Bytes[2], ip.Bytes[3])
		ipv4Addr = netIP.String()
	case 128: // ipv6
		netIP := net.IP{}
		for i := range ip.Bytes {
			netIP = append(netIP, ip.Bytes[i])
		}
		ipv6Addr = netIP.String()
	case 160: // ipv4 + ipv6, and ipv4 is contained in the first 32 bits
		netIPv4 := net.IPv4(ip.Bytes[0], ip.Bytes[1], ip.Bytes[2], ip.Bytes[3])
		netIPv6 := net.IP{}
		for i := range ip.Bytes {
			netIPv6 = append(netIPv6, ip.Bytes[i+4])
		}
		ipv4Addr = netIPv4.String()
		ipv6Addr = netIPv6.String()
	}
	return
}

func IPAddressToNgap(ipv4Addr, ipv6Addr string) ngapType.TransportLayerAddress {
	var ipAddr ngapType.TransportLayerAddress

	switch {
	case ipv4Addr != "" && ipv6Addr != "":
		// Both ipv4 & ipv6
		ipv4NetIP := net.ParseIP(ipv4Addr).To4()
		ipv6NetIP := net.ParseIP(ipv6Addr).To16()

		ipBytes := []byte{ipv4NetIP[0], ipv4NetIP[1], ipv4NetIP[2], ipv4NetIP[3]}
		for _, byteTmp := range ipv6NetIP[:16] {
			ipBytes = append(ipBytes, byteTmp)
		}

		ipAddr.Value = aper.BitString{
			Bytes:     ipBytes,
			BitLength: 160,
		}
	case ipv4Addr != "" && ipv6Addr == "":
		// ipv4
		ipv4NetIP := net.ParseIP(ipv4Addr).To4()

		ipBytes := []byte{ipv4NetIP[0], ipv4NetIP[1], ipv4NetIP[2], ipv4NetIP[3]}

		ipAddr.Value = aper.BitString{
			Bytes:     ipBytes,
			BitLength: 32,
		}
	case ipv4Addr == "" && ipv6Addr != "":
		// ipv6
		ipv6NetIP := net.ParseIP(ipv6Addr).To16()

		ipBytes := []byte{}
		for _, byteTmp := range ipv6NetIP[:16] {
			ipBytes = append(ipBytes, byteTmp)
		}

		ipAddr.Value = aper.BitString{
			Bytes:     ipBytes,
			BitLength: 128,
		}
	default:
		logger.NgapLog.Warnln("IPAddressToNgap: Both ipv4 and ipv6 are empty string")
	}

	return ipAddr
}
