package proxy

// Config is the generic configuration for the server
type Config struct {
	Host        string
	Port        int
	BackendAddr string
}
