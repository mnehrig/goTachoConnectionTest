package main

/*
	errorCodeAuthenticationGenerelReject:
	- communication timeout with remote company car (2*Trem) has expired.
	- Tauth timeout has expired or has been closed.
	- the VU is not able (for internal reason) to perform the remote card authentication, company shall try again later
	- too many errors o/*
	errorCodeAuthenticationGenerelReject:
	- communication timeout with remote company car (2*Trem) has expired.
	- Tauth timeout has expired or has been closed.
	- the VU is not able (for internal reason) to perform the remote card authentication, company shall try again later
	- too many errors on low communication layers.
	- a valid company card, control card or workshop card is inserted in the VU

	errorCodeAuthenticationServiceNotSupported:
	- the requested service is not supported

	errorCodeAuthenticationSubFunctionNotSupported:
	- the routineControlType is neither startRoutine, stopRoutine nor requestRoutineResults

	errorCodeAuthenticationIncorrectMessageLengthOrInvalidFormat:
	- self explaining, isn't it?

	errorCodeAuthenticationBusyRepeatRequest:
	- the VU is busy. The FMS shall perform repetition of this request

	errorCodeAuthenticationConditionsNotCorrect:
	- TAuth is active (another remote authentication process or data transfer is
	already in progress).
	- a differentt DownloadRequestList has already been received by the VU in the
	current authentication process.

	errorCodeAuthenticationRequestSequenceError:
	the sequence of the requests of the authentication process is not correct.
	the stopRoutine or requestRoutineResults subfunction is received, without having
	first received a startRoutine for the requests routineIdentifier

	errorCodeAuthenticationRequestOutOfRange:
	routine identifier parameter is not supported.
	routine identifer 0180 is used in subfunctions stopRoutine or requestRoutineResults.
	reh optional routineControlOptionRecord is not allowed in or contains invalid data
	for the requested routineIdentifier (e.g. max APDU size read in RemoteCompanyCardReady
	is strictly below 240 bytes or strictly above 250 bytes, or a period start in the activities of
	specified calendar day(s) parameter in a RemoteDownloadDataRequest is not immediately
	followed by a period stop...).

	errorCodeAuthenticationRequestCorrectlyReceivedResponsePending:
	Chill. Positive or negative response will follow

	errorCodeAuthenticationServiceNotSupportedInActiveSession
	current session does not support start Routine
	startRoutine RemoteTachographCardDataTransfer service only allowed in remote session.
	in low communication layers.
	- a valid company card, control card or workshop card is inserted in the VU

	errorCodeAuthenticationServiceNotSupported:
	- the requested service is not supported

	errorCodeAuthenticationSubFunctionNotSupported:
	- the routineControlType is neither startRoutine, stopRoutine nor requestRoutineResults

	errorCodeAuthenticationIncorrectMessageLengthOrInvalidFormat:
	- self explaining, isn't it?

	errorCodeAuthenticationBusyRepeatRequest:
	- the VU is busy. The FMS shall perform repetition of this request

	errorCodeAuthenticationConditionsNotCorrect:
	- TAuth is active (another remote authentication process or data transfer is
	already in progress).
	- a differentt DownloadRequestList has already been received by the VU in the
	current authentication process.

	errorCodeAuthenticationRequestSequenceError:
	the sequence of the requests of the authentication process is not correct.
	the stopRoutine or requestRoutineResults subfunction is received, without having
	first received a startRoutine for the requests routineIdentifier

	errorCodeAuthenticationRequestOutOfRange:
	routine identifier parameter is not supported.
	routine identifer 0180 is used in subfunctions stopRoutine or requestRoutineResults.
	reh optional routineControlOptionRecord is not allowed in or contains invalid data
	for the requested routineIdentifier (e.g. max APDU size read in RemoteCompanyCardReady
	is strictly below 240 bytes or strictly above 250 bytes, or a period start in the activities of
	specified calendar day(s) parameter in a RemoteDownloadDataRequest is not immediately
	followed by a period stop...).

	errorCodeAuthenticationRequestCorrectlyReceivedResponsePending:
	Chill. Positive or negative response will follow

	errorCodeAuthenticationServiceNotSupportedInActiveSession
	current session does not support start Routine
	startRoutine RemoteTachographCardDataTransfer service only allowed in remote session.
*/

type errorCodeDiagnosticSessionControl string

const (
	errorCodeDiagnosticSessionControlRequestCorrectlyReceivedResponsePending errorCodeAuthentication = "7F1078"
)

type errorCodeAuthentication string

const (
	errorCodeAuthenticationGenerelReject                           errorCodeAuthentication = "7F3110"
	errorCodeAuthenticationServiceNotSupported                     errorCodeAuthentication = "7F3111"
	errorCodeAuthenticationSubFunctionNotSupported                 errorCodeAuthentication = "7F3112"
	errorCodeAuthenticationIncorrectMessageLengthOrInvalidFormat   errorCodeAuthentication = "7F3113"
	errorCodeAuthenticationBusyRepeatRequest                       errorCodeAuthentication = "7F3121"
	errorCodeAuthenticationConditionsNotCorrect                    errorCodeAuthentication = "7F3122"
	errorCodeAuthenticationRequestSequenceError                    errorCodeAuthentication = "7F3124"
	errorCodeAuthenticationRequestOutOfRange                       errorCodeAuthentication = "7F3131"
	errorCodeAuthenticationRequestCorrectlyReceivedResponsePending errorCodeAuthentication = "7F3178"
	errorCodeAuthenticationServiceNotSupportedInActiveSession      errorCodeAuthentication = "7F317F"
)

type errorCodeAuthenticationPositiveResponse string

const (
	errorCodeAuthenticationThreeConsecutiveAPDUErrorsHaveOccured       errorCodeAuthenticationPositiveResponse = "710101800C"
	errorCodeAuthenticationCardAuthenticationHasFailed                 errorCodeAuthenticationPositiveResponse = "710101800E"
	errorCodeAuthentication5ConsecutiveAuthenticationErrorsHaveOccured errorCodeAuthenticationPositiveResponse = "7101018010"
)

// Negative Response codes during TransferData
type errorCodesRequestUpload string

const (
	errorCodeRequestUploadGeneralReject                           errorCodesRequestUpload = "7F3510"
	errorCodeRequestUploadServiceNotSupported                     errorCodesRequestUpload = "7F3511"
	errorCodeRequestUploadIncorrectMessageLengthOrInvalidFormat   errorCodesRequestUpload = "7F3513"
	errorCodeRequestUploadBusyRepeatRequest                       errorCodesRequestUpload = "7F3521"
	errorCodeRequestUploadConditionsNotCorrect                    errorCodesRequestUpload = "7F3522"
	errorCodeRequestUploadRequestOutOfRange                       errorCodesRequestUpload = "7F3531"
	errorCodeRequestUploadRequestCorrectlyReceivedResponsePending errorCodesRequestUpload = "7F3578"
	errorCodeRequestUploadServiceNotSupportedInActiveSession      errorCodesRequestUpload = "7F357F"
)

// Negative Response codes during TransferData
type errorCodesTransferData string

const (
	errorCodeTransferDataGeneralReject                           errorCodesTransferData = "7F3610"
	errorCodeTransferDataServiceNotSupported                     errorCodesTransferData = "7F3611"
	errorCodeTransferDataIncorrectMessageLengthOrInvalidFormat   errorCodesTransferData = "7F3613"
	errorCodeTransferDataBusyRepeatRequest                       errorCodesTransferData = "7F3621"
	errorCodeTransferDataConditionsNotCorrect                    errorCodesTransferData = "7F3622"
	errorCodeTransferDataRequestSequenceError                    errorCodesTransferData = "7F3624"
	errorCodeTransferDataRequestOutOfRange                       errorCodesTransferData = "7F3631"
	errorCodeTransferDataTransferDataSuspended                   errorCodesTransferData = "7F3671"
	errorCodeTransferDataWrotngBlockSqequenceCounter             errorCodesTransferData = "7F3673"
	errorCodeTransferDataRequestCorrectlyReceivedResponsePending errorCodesTransferData = "7F3678"
	errorCodeTransferDataServiceNotSupportedInActiveSession      errorCodesTransferData = "7F367F"
)

// Negative Response codes during TransferData
type errorCodesTransferExit string

const (
	errorCodeTransferExitGeneralReject                           errorCodesTransferExit = "7F3710"
	errorCodeTransferExitServiceNotSupported                     errorCodesTransferExit = "7F3711"
	errorCodeTransferExitIncorrectMessageLengthOrInvalidFormat   errorCodesTransferExit = "7F3713"
	errorCodeTransferExitBusyRepeatRequest                       errorCodesTransferExit = "7F3721"
	errorCodeTransferExitConditionsNotCorrect                    errorCodesTransferExit = "7F3722"
	errorCodeTransferExitRequestSequenceError                    errorCodesTransferExit = "7F3724"
	errorCodeTransferExitRequestOutOfRange                       errorCodesTransferExit = "7F3731"
	errorCodeTransferExitTransferDataSuspended                   errorCodesTransferExit = "7F3771"
	errorCodeTransferExitRequestCorrectlyReceivedResponsePending errorCodesTransferExit = "7F3778"
	errorCodeTransferExitServiceNotSupportedInActiveSession      errorCodesTransferExit = "7F377F"
)

type state string

const (

	// State while authentication:
	stateUndefined                    state = "stateUndefined"
	stateRemoteAuthenticationClosed   state = "stateRemoteAuthenticationClosed"
	stateVUReady                      state = "stateVUReady"
	stateVUToCompanyCardData          state = "stateVUToCompanyCardData"
	stateRemoteDownloadAccessGranted  state = "stateRemoteDownloadAccessGranted"
	stateRemoteAuthenticatedSucceeded state = "stateRemoteAuthenticatedSucceeded"

	// Error States while authentication:
	stateAPDUError                   state = "stateAPDUError"
	stateAuthenticationError         state = "stateAuthenticationError"
	stateTooManyAuthenticationErrors state = "stateTooManyAuthenticationErrors"

	// States while data download

	// Error States while authentication:
)

type resetLevel string

const (
	resetLevelTotal resetLevel = "resetLevelTotal"
)
