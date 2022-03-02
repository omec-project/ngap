// Copyright 2019 Communication Service/Software Laboratory, National Chiao Tung University (free5gc.org)
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

// Need to import "github.com/free5gc/aper" if it uses "aper"

type UserLocationInformationEUTRA struct {
	EUTRACGI     EUTRACGI                                                      `aper:"valueExt"`
	TAI          TAI                                                           `aper:"valueExt"`
	TimeStamp    *TimeStamp                                                    `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerUserLocationInformationEUTRAExtIEs `aper:"optional"`
}
