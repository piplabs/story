package config

type RollbackConfig struct {
	RollbackHeights uint64
}

func DefaultRollbackConfig() RollbackConfig {
	return RollbackConfig{
		RollbackHeights: 1,
	}
}
