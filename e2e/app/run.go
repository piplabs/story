package app

import (
	"context"

	"github.com/piplabs/story/e2e/types"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/log"
)

func DefaultDeployConfig() DeployConfig {
	return DeployConfig{}
}

type DeployConfig struct {
	// Internal use parameters (no command line flags).
	testConfig bool
}

// Deploy a new e2e network.
func Deploy(ctx context.Context, def Definition, cfg DeployConfig) error {
	if def.Testnet.Network.IsProtected() {
		// If a protected network needs to be deployed temporarily comment out this check.
		return errors.New("cannot deploy protected network", "network", def.Testnet.Network)
	}

	if err := deployPublicCreate3(ctx, def); err != nil {
		return err
	}

	if err := Setup(ctx, def, cfg); err != nil {
		return err
	}

	// Only stop and delete existing network right before actually starting new ones.
	if err := CleanInfra(ctx, def); err != nil {
		return err
	}

	if err := StartInitial(ctx, def.Testnet.Testnet, def.Infra); err != nil {
		return err
	}

	if err := fundAccounts(ctx, def); err != nil {
		return err
	}

	if err := deployPrivateCreate3(ctx, def); err != nil {
		return err
	}

	//nolint:revive // Will add more logic after this if check
	if err := FundValidatorsForTesting(ctx, def); err != nil {
		return err
	}

	return nil
}

// E2ETestConfig is the configuration required to run a full e2e test.
type E2ETestConfig struct {
	Preserve bool
}

// DefaultE2ETestConfig returns a default configuration for a e2e test.
func DefaultE2ETestConfig() E2ETestConfig {
	return E2ETestConfig{}
}

// E2ETest runs a full e2e test.
func E2ETest(ctx context.Context, def Definition, cfg E2ETestConfig) error {
	stopValidatorUpdates := StartValidatorUpdates(ctx, def)

	if err := StartRemaining(ctx, def.Testnet.Testnet, def.Infra); err != nil {
		return err
	}

	if err := Wait(ctx, def.Testnet.Testnet, 5); err != nil { // allow some txs to go through
		return err
	}

	if def.Testnet.HasPerturbations() {
		if err := perturb(ctx, def.Testnet); err != nil {
			return err
		}
	}

	if def.Testnet.Evidence > 0 {
		return errors.New("evidence injection not supported yet")
	}

	if err := stopValidatorUpdates(); err != nil {
		return errors.Wrap(err, "stop validator updates")
	}

	// Start unit tests.
	if err := Test(ctx, def, false); err != nil {
		return err
	}

	if cfg.Preserve {
		log.Warn(ctx, "Docker containers not stopped, --preserve=true", nil)
	} else if err := CleanInfra(ctx, def); err != nil {
		return err
	}

	return nil
}

// Upgrade generates all local artifacts, but only copies the docker-compose file to the VMs.
// It them calls docker-compose up.
func Upgrade(ctx context.Context, def Definition, cfg DeployConfig, upgradeCfg types.UpgradeConfig) error {
	if err := Setup(ctx, def, cfg); err != nil {
		return err
	}

	return def.Infra.Upgrade(ctx, upgradeCfg)
}
