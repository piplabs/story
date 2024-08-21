package contracts

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/storyprotocol/iliad/e2e/app/eoa"
	"github.com/storyprotocol/iliad/lib/create3"
	"github.com/storyprotocol/iliad/lib/netconf"
)

//
// Create3Factory.
//

func MainnetCreate3Factory() common.Address {
	return crypto.CreateAddress(eoa.MustAddress(netconf.Mainnet, eoa.RoleCreate3Deployer), 0)
}

func TestnetCreate3Factory() common.Address {
	return crypto.CreateAddress(eoa.MustAddress(netconf.Testnet, eoa.RoleCreate3Deployer), 0)
}

func StagingCreate3Factory() common.Address {
	return crypto.CreateAddress(eoa.MustAddress(netconf.Staging, eoa.RoleCreate3Deployer), 0)
}

func DevnetCreate3Factory() common.Address {
	return crypto.CreateAddress(eoa.MustAddress(netconf.Devnet, eoa.RoleCreate3Deployer), 0)
}

//
// IPTokenStaking.
//

func MainnetIPTokenStaking() common.Address {
	return create3.Address(MainnetCreate3Factory(), IPTokenStakingSalt(netconf.Mainnet), eoa.MustAddress(netconf.Mainnet, eoa.RoleDeployer))
}

func TestnetIPTokenStaking() common.Address {
	return create3.Address(TestnetCreate3Factory(), IPTokenStakingSalt(netconf.Testnet), eoa.MustAddress(netconf.Testnet, eoa.RoleDeployer))
}

func StagingIPTokenStaking() common.Address {
	return create3.Address(StagingCreate3Factory(), IPTokenStakingSalt(netconf.Staging), eoa.MustAddress(netconf.Staging, eoa.RoleDeployer))
}

func DevnetIPTokenStaking() common.Address {
	return create3.Address(DevnetCreate3Factory(), IPTokenStakingSalt(netconf.Devnet), eoa.MustAddress(netconf.Devnet, eoa.RoleDeployer))
}

//
// Salts.
//

func IPTokenStakingSalt(network netconf.ID) string {
	// only portal salts are versioned
	return salt(network, "iptokenstaking-"+network.Version())
}

//
// Utils.
//

// salt generates a salt for a contract deployment. For ephemeral networks,
// the salt includes a random per-run suffix. For persistent networks, the
// sale is static.
func salt(network netconf.ID, contract string) string {
	return string(network) + "-" + contract
}
