package server

import "github.com/spf13/pflag"

// BindFlags binds the provided flags to the corresponding fields in the Config struct.
func BindFlags(flags *pflag.FlagSet, cfg *Config) {
	flags.BoolVar(&cfg.APIEnable, "api-enable", cfg.APIEnable, "Define if the API server should be enabled")
	flags.StringVar(&cfg.APIAddress, "api-address", cfg.APIAddress, "The API server address to listen on")
	flags.BoolVar(&cfg.EnableUnsafeCORS, "enabled-unsafe-cors", cfg.EnableUnsafeCORS, "Enable unsafe CORS for API server")
	flags.UintVar(&cfg.ReadTimeout, "read-timeout", cfg.ReadTimeout, "Define the API server read timeout (in seconds)")
	flags.UintVar(&cfg.ReadHeaderTimeout, "read-header-timeout", cfg.ReadHeaderTimeout, "Define the API server read header timeout (in seconds)")
	flags.UintVar(&cfg.WriteTimeout, "write-timeout", cfg.WriteTimeout, "Define the API server write timeout (in seconds)")
	flags.UintVar(&cfg.IdleTimeout, "idle-timeout", cfg.IdleTimeout, "Define the API server idle timeout (in seconds)")
	flags.UintVar(&cfg.MaxHeaderBytes, "max-header-bytes", cfg.MaxHeaderBytes, "Define the API server max header (in bytes)")
}

type Config struct {
	APIEnable         bool
	APIAddress        string
	EnableUnsafeCORS  bool
	ReadTimeout       uint
	ReadHeaderTimeout uint
	WriteTimeout      uint
	IdleTimeout       uint
	MaxHeaderBytes    uint
}

func DefaultConfig() Config {
	return Config{
		APIEnable:         false,
		APIAddress:        "127.0.0.1:1317",
		EnableUnsafeCORS:  false,
		ReadTimeout:       10,
		ReadHeaderTimeout: 10,
		WriteTimeout:      10,
		IdleTimeout:       10,
		MaxHeaderBytes:    8 << 10,
	}
}
