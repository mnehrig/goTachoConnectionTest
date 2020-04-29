package main

type canMessage string

const (
	TesterPresent                            = canMessage("023E80FFFFFFFFFF")
	RequestTransferExit                      = canMessage("023700FFFFFFFFFF")
	PositiveResponseRemoteSessionExit        = "T18DAF0EE8025001FFFFFFFFFF"
	PositiveResponseRemoteAuthenticationExit = "T18DAF0EE805710101800AFFFF"
	NegativeResponsePending                  = "T18DAF0EE8037F3178FFFFFFFF"
)
