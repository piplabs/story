package server

import "github.com/spf13/pflag"

// BindFlags binds the provided flags to the corresponding fields in the Config struct.
func BindFlags(flags *pflag.FlagSet, cfg *Config) {
	flags.BoolVar(&cfg.Enable, "api-enable", cfg.Enable, "Define if the API server should be enabled")
	flags.StringVar(&cfg.Address, "api-address", cfg.Address, "The API server address to listen on")
	flags.BoolVar(&cfg.EnableUnsafeCORS, "api-enable-unsafe-cors", cfg.EnableUnsafeCORS, "Enable unsafe CORS for API server")
	flags.UintVar(&cfg.ReadTimeout, "api-read-timeout", cfg.ReadTimeout, "Define the API server read timeout (in seconds)")
	flags.UintVar(&cfg.ReadHeaderTimeout, "api-read-header-timeout", cfg.ReadHeaderTimeout, "Define the API server read header timeout (in seconds)")
	flags.UintVar(&cfg.WriteTimeout, "api-write-timeout", cfg.WriteTimeout, "Define the API server write timeout (in seconds)")
	flags.UintVar(&cfg.IdleTimeout, "api-idle-timeout", cfg.IdleTimeout, "Define the API server idle timeout (in seconds)")
	flags.UintVar(&cfg.MaxHeaderBytes, "api-max-header-bytes", cfg.MaxHeaderBytes, "Define the API server max header (in bytes)")
}

type Config struct {
	Enable            bool
	Address           string
	EnableUnsafeCORS  bool
	ReadTimeout       uint
	ReadHeaderTimeout uint
	WriteTimeout      uint
	IdleTimeout       uint
	MaxHeaderBytes    uint
}

func DefaultConfig() Config {
	return Config{
		Enable:            false,
		Address:           "127.0.0.1:1317",
		EnableUnsafeCORS:  false,
		ReadTimeout:       10,
		ReadHeaderTimeout: 10,
		WriteTimeout:      10,
		IdleTimeout:       10,
		MaxHeaderBytes:    8 << 10,
	}
}
