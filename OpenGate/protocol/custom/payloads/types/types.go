package types

const (
	Hello                 byte = 1   // hello type to identify other client
	KeepAlive             byte = 2   // Keepalive to help the connection stay open
	Abort                 byte = 4   // Stop communicating. Universal request
	ClientSend            byte = 8   // Client sends a messag
	ClientReceiveResponse byte = 16  // Answer that the message has been delivered
	ClientSyncReq         byte = 32  // Client asks the signal server for syncing. SPECIAL request type containing ip within the MESSAGE field
	MalformedReqErr       byte = 64  // Client failed to identify the message from server or from another client
	ClientError           byte = 128 // Client encountered an error
)
