// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/omec-project/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct UPTransportLayerInformationList */
/* UPTransportLayerInformationItem */
type UPTransportLayerInformationList struct {
	List []UPTransportLayerInformationItem `aper:"valueExt,sizeLB:1,sizeUB:3"`
}
