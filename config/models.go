package config

var (
	// Server consists all servers that is cacheable
	Server = []string{"CN", "US", "JP", "KR"}
)

// ValidServer validates if server counts as a valid server
func ValidServer(server string) bool {
	for _, s := range Server {
		if s == server {
			return true
		}
	}
	return false
}
