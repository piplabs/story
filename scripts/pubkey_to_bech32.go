package scripts

import (
	"encoding/base64"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/piplabs/story/lib/k1util"
)

//nolint:unused // main for go run
func main() {
	Bech32HRP := "story"
	accountPubKeyPrefix := Bech32HRP + "pub"
	validatorAddressPrefix := Bech32HRP + "valoper"
	validatorPubKeyPrefix := Bech32HRP + "valoperpub"
	consNodeAddressPrefix := Bech32HRP + "valcons"
	consNodePubKeyPrefix := Bech32HRP + "valconspub"

	// Set and seal config
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(Bech32HRP, accountPubKeyPrefix)
	cfg.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	cfg.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	cfg.Seal()

	// Input the base64 encoded public key (secp256k1)
	b64e := "A+m9C7PAkf55wktLV5JfSLir1Hl2u7rKx1Vuv5ZZ53Cl" // pragma: allowlist secret
	b64d, _ := base64.StdEncoding.DecodeString(b64e)

	pubkey, err := k1util.PubKeyBytesToCosmos(b64d)
	if err != nil {
		println(err.Error())
		return
	}
	println("pubkey", pubkey.String())

	evmAddr, err := k1util.CosmosPubkeyToEVMAddress(pubkey.Bytes())
	if err != nil {
		println(err.Error())
		return
	}

	accAddr := sdk.AccAddress(evmAddr.Bytes())
	valAddr := sdk.ValAddress(evmAddr.Bytes())
	println("accAddr", accAddr.String())
	println("valAddr", valAddr.String())

	println("evmAddr", evmAddr.String())
}
