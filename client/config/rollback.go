package config

type RollbackConfig struct {
	RemoveBlock     bool // See cosmos-sdk/server/rollback.go
	RollbackHeights uint64
}

func DefaultRollbackConfig() RollbackConfig {
	return RollbackConfig{
		RemoveBlock:     false,
		RollbackHeights: 1,
	}
}
