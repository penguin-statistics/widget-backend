package errors

var (
	ErrInvalidServer = New("InvalidServer", "malformed parameter `server` provided", BlameUser)
	ErrCantMarshal   = New("CantMarshal", "failed to populate matrix with metadata provided", BlameServer)
)
