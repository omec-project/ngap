// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/free5gc/aper" if it uses "aper"

const (
	CPTransportLayerInformationPresentNothing int = iota /* No components present */
	CPTransportLayerInformationPresentEndpointIPAddress
	CPTransportLayerInformationPresentChoiceExtensions
)

type CPTransportLayerInformation struct {
	Present           int
	EndpointIPAddress *TransportLayerAddress
	ChoiceExtensions  *ProtocolIESingleContainerCPTransportLayerInformationExtIEs
}
