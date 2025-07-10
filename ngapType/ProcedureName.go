// Copyright 2025 Canonical, Ltd.
//
// SPDX-License-Identifier: Apache-2.0

package ngapType

import "fmt"

func ProcedureName(code int64) string {
	switch code {
	case ProcedureCodeAMFConfigurationUpdate:
		return "AMFConfigurationUpdate"
	case ProcedureCodeAMFStatusIndication:
		return "AMFStatusIndication"
	case ProcedureCodeCellTrafficTrace:
		return "CellTrafficTrace"
	case ProcedureCodeDeactivateTrace:
		return "DeactivateTrace"
	case ProcedureCodeDownlinkNASTransport:
		return "DownlinkNASTransport"
	case ProcedureCodeDownlinkNonUEAssociatedNRPPaTransport:
		return "DownlinkNonUEAssociatedNRPPaTransport"
	case ProcedureCodeDownlinkRANConfigurationTransfer:
		return "DownlinkRANConfigurationTransfer"
	case ProcedureCodeDownlinkRANStatusTransfer:
		return "DownlinkRANStatusTransfer"
	case ProcedureCodeDownlinkUEAssociatedNRPPaTransport:
		return "DownlinkUEAssociatedNRPPaTransport"
	case ProcedureCodeErrorIndication:
		return "ErrorIndication"
	case ProcedureCodeHandoverCancel:
		return "HandoverCancel"
	case ProcedureCodeHandoverNotification:
		return "HandoverNotification"
	case ProcedureCodeHandoverPreparation:
		return "HandoverPreparation"
	case ProcedureCodeHandoverResourceAllocation:
		return "HandoverResourceAllocation"
	case ProcedureCodeInitialContextSetup:
		return "InitialContextSetup"
	case ProcedureCodeInitialUEMessage:
		return "InitialUEMessage"
	case ProcedureCodeLocationReportingControl:
		return "LocationReportingControl"
	case ProcedureCodeLocationReportingFailureIndication:
		return "LocationReportingFailureIndication"
	case ProcedureCodeLocationReport:
		return "LocationReport"
	case ProcedureCodeNASNonDeliveryIndication:
		return "NASNonDeliveryIndication"
	case ProcedureCodeNGReset:
		return "NGReset"
	case ProcedureCodeNGSetup:
		return "NGSetup"
	case ProcedureCodeOverloadStart:
		return "OverloadStart"
	case ProcedureCodeOverloadStop:
		return "OverloadStop"
	case ProcedureCodePaging:
		return "Paging"
	case ProcedureCodePathSwitchRequest:
		return "PathSwitchRequest"
	case ProcedureCodePDUSessionResourceModify:
		return "PDUSessionResourceModify"
	case ProcedureCodePDUSessionResourceModifyIndication:
		return "PDUSessionResourceModifyIndication"
	case ProcedureCodePDUSessionResourceRelease:
		return "PDUSessionResourceRelease"
	case ProcedureCodePDUSessionResourceSetup:
		return "PDUSessionResourceSetup"
	case ProcedureCodePDUSessionResourceNotify:
		return "PDUSessionResourceNotify"
	case ProcedureCodePrivateMessage:
		return "PrivateMessage"
	case ProcedureCodePWSCancel:
		return "PWSCancel"
	case ProcedureCodePWSFailureIndication:
		return "PWSFailureIndication"
	case ProcedureCodePWSRestartIndication:
		return "PWSRestartIndication"
	case ProcedureCodeRANConfigurationUpdate:
		return "RANConfigurationUpdate"
	case ProcedureCodeRerouteNASRequest:
		return "RerouteNASRequest"
	case ProcedureCodeRRCInactiveTransitionReport:
		return "RRCInactiveTransitionReport"
	case ProcedureCodeTraceFailureIndication:
		return "TraceFailureIndication"
	case ProcedureCodeTraceStart:
		return "TraceStart"
	case ProcedureCodeUEContextModification:
		return "UEContextModification"
	case ProcedureCodeUEContextRelease:
		return "UEContextRelease"
	case ProcedureCodeUEContextReleaseRequest:
		return "UEContextReleaseRequest"
	case ProcedureCodeUERadioCapabilityCheck:
		return "UERadioCapabilityCheck"
	case ProcedureCodeUERadioCapabilityInfoIndication:
		return "UERadioCapabilityInfoIndication"
	case ProcedureCodeUETNLABindingRelease:
		return "UETNLABindingRelease"
	case ProcedureCodeUplinkNASTransport:
		return "UplinkNASTransport"
	case ProcedureCodeUplinkNonUEAssociatedNRPPaTransport:
		return "UplinkNonUEAssociatedNRPPaTransport"
	case ProcedureCodeUplinkRANConfigurationTransfer:
		return "UplinkRANConfigurationTransfer"
	case ProcedureCodeUplinkRANStatusTransfer:
		return "UplinkRANStatusTransfer"
	case ProcedureCodeUplinkUEAssociatedNRPPaTransport:
		return "UplinkUEAssociatedNRPPaTransport"
	case ProcedureCodeWriteReplaceWarning:
		return "WriteReplaceWarning"
	case ProcedureCodeSecondaryRATDataUsageReport:
		return "SecondaryRATDataUsageReport"
	default:
		return fmt.Sprintf("unknown procedure code %d", code)
	}
}
