package errors

var (
	// ErrInvalidServer describes an error in which a malformed `server` provided
	ErrInvalidServer = New("InvalidServer", "malformed `server` provided", BlameUser)
)
