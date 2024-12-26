package utils

import (
	cosmosk1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/piplabs/story/lib/errors"
)

func EvmAddressToBech32AccAddress(evmAddress string) (sdk.AccAddress, error) {
	addressBytes, err := hexutil.Decode(evmAddress)
	if err != nil {
		return nil, err
	}

	return sdk.AccAddress(addressBytes), nil
}

func EvmAddressToBech32ValAddress(evmAddress string) (sdk.ValAddress, error) {
	addressBytes, err := hexutil.Decode(evmAddress)
	if err != nil {
		return nil, err
	}

	return sdk.ValAddress(addressBytes), nil
}

func CmpPubKeyToBech32ConsAddress(hexCmpPubKey string) (sdk.ConsAddress, error) {
	cmpPubKeyBytes, err := hexutil.Decode(hexCmpPubKey)
	if err != nil {
		return nil, err
	}

	cmpPubKey := &cosmosk1.PubKey{Key: cmpPubKeyBytes}

	return sdk.GetConsAddress(cmpPubKey), nil
}

func Bech32DelegatorAddressToEvmAddress(bech32DelegatorAddress string) (string, error) {
	delegatorAccAddress, err := sdk.AccAddressFromBech32(bech32DelegatorAddress)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert bech32 delegator address to evm address")
	}

	return hexutil.Encode(delegatorAccAddress.Bytes()), nil
}

func Bech32ValidatorAddressToEvmAddress(bech32ValidatorAddress string) (string, error) {
	validatorAccAddress, err := sdk.ValAddressFromBech32(bech32ValidatorAddress)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert bech32 validator address to evm address")
	}

	return hexutil.Encode(validatorAccAddress.Bytes()), nil
}
