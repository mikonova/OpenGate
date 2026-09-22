package types

const (
	Hello                 byte = 1
	KeepAlive             byte = 2
	Abort                 byte = 4
	ClientSendRequest     byte = 8
	ClientReceiveResponse byte = 16
	MalformedReqErr       byte = 32
	ClientError           byte = 64
)
