package app

import (
	"context"
	"crypto/ecdsa"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"

	"cosmossdk.io/math"

	e2e "github.com/cometbft/cometbft/test/e2e/pkg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"

	"github.com/piplabs/story/client/genutil/evm/predeploys"
	"github.com/piplabs/story/contracts/bindings"
	"github.com/piplabs/story/e2e/app/eoa"
	"github.com/piplabs/story/lib/errors"
	"github.com/piplabs/story/lib/ethclient"
	"github.com/piplabs/story/lib/ethclient/ethbackend"
	"github.com/piplabs/story/lib/k1util"
	"github.com/piplabs/story/lib/log"
	"github.com/piplabs/story/lib/txmgr"
)

// FundValidatorsForTesting funds validators in ephemeral networks: devnet and staging.
// This is required by load generation for periodic validator self-delegation.
func FundValidatorsForTesting(ctx context.Context, def Definition) error {
	if !def.Testnet.Network.IsEphemeral() {
		// Only fund validators in ephemeral networks, devnet and staging.
		return nil
	}

	log.Info(ctx, "Funding validators for testing", "count", len(def.Testnet.Nodes))

	network := networkFromDef(def)
	omniEVM, _ := network.IliadEVMChain()
	funder := eoa.MustAddress(network.ID, eoa.RoleTester) // Fund validators using tester eoa
	_, fundBackend, err := def.Backends().BindOpts(ctx, omniEVM.ID, funder)
	if err != nil {
		return errors.Wrap(err, "bind opts")
	}

	// Iterate over all nodes, since all maybe become validators.
	var eg errgroup.Group
	for _, node := range def.Testnet.Nodes {
		eg.Go(func() error {
			addr, _ := k1util.PubKeyToAddress(node.PrivvalKey.PubKey())
			tx, _, err := fundBackend.Send(ctx, funder, txmgr.TxCandidate{
				To:       &addr,
				GasLimit: 100_000,
				Value:    math.NewInt(1000).MulRaw(params.Ether).BigInt(),
			})
			if err != nil {
				return errors.Wrap(err, "send")
			}
			recp, err := fundBackend.WaitMined(ctx, tx)
			if err != nil {
				return errors.Wrap(err, "wait mined")
			}

			bal, err := fundBackend.EtherBalanceAt(ctx, addr)
			if err != nil {
				return err
			}

			log.Debug(ctx, "Funded validator address",
				"node", node.Name, "addr", addr,
				"balance", bal, "height", recp.BlockNumber.Uint64())

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "wait fund")
	}

	return nil
}

type valUpdate struct {
	Height int64
	Powers map[*e2e.Node]int64
}

func StartValidatorUpdates(ctx context.Context, def Definition) func() error {
	errChan := make(chan error, 1)
	returnErr := func(err error) {
		if err != nil {
			log.Error(ctx, "Validator updates failed", err)
		}
		select {
		case errChan <- err:
		default:
			log.Error(ctx, "Error channel full, dropping error", err)
		}
	}

	go func() {
		// Get all private keys
		var privkeys []*ecdsa.PrivateKey
		for _, node := range def.Testnet.Nodes {
			pk, err := k1util.StdPrivKeyFromComet(node.PrivvalKey)
			if err != nil {
				returnErr(err)
				return
			}

			privkeys = append(privkeys, pk)
		}

		// Get a sorted list of validator updates
		var updates []valUpdate
		for height, powers := range def.Testnet.ValidatorUpdates {
			updates = append(updates, valUpdate{
				Height: height,
				Powers: powers,
			})
		}
		sort.Slice(updates, func(i, j int) bool {
			return updates[i].Height < updates[j].Height
		})

		// Create a backend to trigger deposits from
		network := networkFromDef(def)
		endpoints := externalEndpoints(def)
		iliadEVM, _ := network.IliadEVMChain()
		rpc, found := endpoints[iliadEVM.Name]
		if !found {
			returnErr(errors.New("get rpc"))
			return
		}
		ethCl, err := ethclient.Dial(iliadEVM.Name, rpc)
		if err != nil {
			returnErr(errors.Wrap(err, "dial"))
			return
		}
		valBackend, err := ethbackend.NewBackend(iliadEVM.Name, iliadEVM.ID, iliadEVM.BlockPeriod, ethCl, privkeys...)
		if err != nil {
			returnErr(errors.Wrap(err, "new backend"))
			return
		}

		// Create the IPTokenStaking contract
		ipTokenStaking, err := bindings.NewIPTokenStaking(common.HexToAddress(predeploys.IPTokenStaking), valBackend)
		if err != nil {
			returnErr(errors.Wrap(err, "new iliad stake"))
			return
		}

		// Wait for each update, then submit self-delegations
		for _, update := range updates {
			log.Debug(ctx, "Waiting for next validator update", "wait_for_height", update.Height)
			_, _, err := waitForHeight(ctx, def.Testnet.Testnet, update.Height)
			if err != nil {
				returnErr(errors.Wrap(err, "wait for height"))
				return
			}

			for node, power := range update.Powers {
				pubkey := node.PrivvalKey.PubKey()
				addr, err := k1util.PubKeyToAddress(pubkey)
				if err != nil {
					returnErr(errors.Wrap(err, "pubkey to addr"))
					return
				}

				// Wait until we have enough balance.
				// FundValidatorsForTesting should ensure this, but this sometimes fails...?
				for range 10 {
					height, err := valBackend.BlockNumber(ctx)
					if err != nil {
						returnErr(errors.Wrap(err, "block height"))
						return
					}

					balance, err := valBackend.EtherBalanceAt(ctx, addr)
					if err != nil {
						returnErr(errors.Wrap(err, "balance at"))
						return
					}

					if balance > float64(power) {
						break // We have enough balance
					}

					log.Warn(ctx, "Cannot self-delegate, balance to low (will retry)", nil,
						"height", height, "balance", balance, "require", power,
						"node", node.Name, "addr", addr.Hex())
					time.Sleep(time.Second)
				}

				txOpts, err := valBackend.BindOpts(ctx, addr)
				if err != nil {
					returnErr(errors.Wrap(err, "bind opts"))
					return
				}
				txOpts.Value = math.NewInt(power).MulRaw(params.Ether).BigInt()

				// NOTE: We can use CreateValidator here, rather than Delegate (self-delegation)
				// because current e2e manifest validator_udpates are only used to create a new validator,
				// and not to self-delegate an existing one.
				// TODO: Use CreateValidator
				tx, err := ipTokenStaking.CreateValidatorOnBehalf(txOpts, pubkey.Bytes())
				if err != nil {
					returnErr(errors.Wrap(err, "deposit", "node", node.Name, "addr", addr.Hex()))
					return
				}
				rec, err := valBackend.WaitMined(ctx, tx)
				if err != nil {
					returnErr(errors.Wrap(err, "wait minded", "node", node.Name, "addr", addr.Hex()))
					return
				}

				log.Info(ctx, "Deposited stake",
					"validator", node.Name,
					"address", addr.Hex(),
					"power", power,
					"height", rec.BlockNumber.Uint64(),
				)
			}
		}

		returnErr(nil)
	}()

	return func() error {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "timeout")
		}
	}
}
