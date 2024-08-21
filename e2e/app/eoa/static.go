package eoa

import (
	"github.com/piplabs/story/lib/anvil"
	"github.com/piplabs/story/lib/netconf"
)

const (
	// fbFunder is the address of the fireblocks "funder" account.
	fbFunder  = "0xf63316AA39fEc9D2109AB0D9c7B1eE3a6F60AEA4"
	fbDev     = "0x7a6cF389082dc698285474976d7C75CAdE08ab7e"
	ZeroXDead = "0x000000000000000000000000000000000000dead"
)

//nolint:gochecknoglobals // Static mappings.
var statics = map[netconf.ID][]Account{
	netconf.Devnet: flatten(
		wellKnown(anvil.DevPrivateKey0(), RoleTester, RoleCreate3Deployer, RoleDeployer),
	),
	netconf.Staging: flatten(
		remote("0xC8103859Ac7CB547d70307EdeF1A2319FC305fdC", RoleCreate3Deployer),
		remote("0x274c4B3e5d27A65196d63964532366872F81D261", RoleDeployer),
	),
	netconf.Testnet: flatten(
		remote("0xeC5134556da0797A5C5cD51DD622b689Cac97Fe9", RoleCreate3Deployer),
		remote("0x0CdCc644158b7D03f40197f55454dc7a11Bd92c1", RoleDeployer),
	),
	netconf.Mainnet: flatten(
		dummy(RoleAdmin, RoleCreate3Deployer, RoleDeployer, RoleTester),
	),
}
