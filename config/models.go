package config

var (
	// Server consists all servers that is cacheable
	Server = []string{"CN", "US", "JP", "KO"}
)

func ValidServer(server string) bool {
	for _, s := range Server {
		if s == server {
			return true
		}
	}
	return false
}