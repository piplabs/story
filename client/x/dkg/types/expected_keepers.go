package types

import (
	"context"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type AccountKeeper interface {
	AddressCodec() address.Codec
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type StakingKeeper interface {
	GetAllValidators(ctx context.Context) (validators []stakingtypes.Validator, err error)
}

type BankKeeper interface {
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type DistributionKeeper interface {
	GetUbiBalanceByDenom(ctx context.Context, denom string) (math.Int, error)
	WithdrawUbiByDenomToModule(ctx context.Context, denom string, recipientModule string) (sdk.Coin, error)
}

// DKGContractClient defines the interface for interacting with DKG and CDR smart contracts.
type DKGContractClient interface {
	Register(ctx context.Context, round uint32, enclaveType [32]byte, startBlockHeight uint64, startBlockHash []byte, dkgPubKey []byte, commPubKey []byte, enclaveReport []byte) (*ethtypes.Receipt, error)
	Finalize(ctx context.Context, round uint32, enclaveType [32]byte, participantsRoot []byte, globalPubKey []byte, publicCoeffs [][]byte, pubKeyShare []byte, signature []byte) (*ethtypes.Receipt, error)
	SubmitEncryptedPartialDecryption(ctx context.Context, round uint32, pid uint32, encryptedPartial []byte, ephemeralPubKey []byte, pubShare []byte, requesterPubKey []byte, ciphertext []byte, uuid uint32, signature []byte) (*ethtypes.Receipt, error)
}
