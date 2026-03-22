package keeper

import (
	"context"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	dkgtestutil "github.com/piplabs/story/client/x/dkg/testutil"
	"github.com/piplabs/story/client/x/dkg/types"
)

// TestHandleDecryptRequest_FullPath exercises the complete decrypt request flow
// with mocked kernel and contract clients.
func TestHandleDecryptRequest_FullPath(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("decrypt-cc")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)
	mockContract := dkgtestutil.NewMockDKGContractClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		kernelRouter:   router,
		contractClient: mockContract,
	}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	label := make([]byte, 32)
	label[31] = 42 // UUID = 42

	req := types.DecryptRequest{
		Ciphertext:      []byte("ciphertext"),
		Label:           label,
		RequesterPubKey: []byte("requester-pub"),
	}

	// Kernel PartialDecryptTDH2 returns partial decrypt data
	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).Return(
		&types.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("partial"),
			EphemeralPubKey:            []byte("eph-pub"),
			PubShare:                   []byte("pub-share"),
			Signature:                  []byte("sig"),
		}, nil,
	)

	// Contract SubmitEncryptedPartialDecryption
	mockContract.EXPECT().SubmitEncryptedPartialDecryption(
		gomock.Any(),
		uint32(1), // round
		uint32(2), // pid
		[]byte("partial"),
		[]byte("eph-pub"),
		[]byte("pub-share"),
		[]byte("requester-pub"),
		[]byte("ciphertext"),
		uint32(42), // uuid from label
		[]byte("sig"),
	).Return(&ethtypes.Receipt{Status: ethtypes.ReceiptStatusSuccessful}, nil)

	err := k.handleDecryptRequest(ctx, session, req)
	require.NoError(t, err)
}

// TestHandleDecryptRequest_InvalidLabel exercises the label-too-short error.
func TestHandleDecryptRequest_InvalidLabel(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctrl := gomock.NewController(t)

	cc := []byte("decrypt-cc-short")
	mockKernel := dkgtestutil.NewMockKernelServiceClient(ctrl)

	router := NewKernelRouter(nil, nil)
	router.RegisterClient(cc, mockKernel)

	k := &Keeper{
		kernelRouter:   router,
		contractClient: nil, // won't be called
	}

	session := &types.DKGSession{
		Round:          1,
		Index:          2,
		GlobalPubKey:   []byte("global-pub"),
		CodeCommitment: cc,
	}

	// Label too short
	req := types.DecryptRequest{
		Ciphertext:      []byte("ciphertext"),
		Label:           make([]byte, 20), // < 32 bytes
		RequesterPubKey: []byte("requester-pub"),
	}

	// Kernel PartialDecryptTDH2 returns success, but labelToUUID will fail
	mockKernel.EXPECT().PartialDecryptTDH2(gomock.Any(), gomock.Any()).Return(
		&types.PartialDecryptTDH2Response{
			EncryptedPartialDecryption: []byte("partial"),
			EphemeralPubKey:            []byte("eph-pub"),
			PubShare:                   []byte("pub-share"),
			Signature:                  []byte("sig"),
		}, nil,
	)

	err := k.handleDecryptRequest(ctx, session, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid decrypt request label")
}
