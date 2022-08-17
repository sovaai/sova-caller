package sip

// IUserService : User callback functions
type IUserService interface {
	OnSipReady()
	OnRegState(uid string, isActive bool, code int)
	OnCallClosed(callId string, reason string)
	OnUnregister()
}
