package keeper

import (
	"github.com/piplabs/story/client/x/dkg/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnqueueAndDequeueDeals(t *testing.T) {
	k, _ := setupDKGKeeper(t)

	deals := []types.Deal{
		{
			Index:          1,
			RecipientIndex: 10,
			Signature:      []byte("sig1"),
			Deal: types.EncryptedDeal{
				DhKey:     []byte("dh1"),
				Signature: []byte("sig1"),
				Nonce:     []byte("nonce1"),
				Cipher:    []byte("cipher1"),
			},
		},
		{
			Index:          2,
			RecipientIndex: 20,
			Signature:      []byte("sig2"),
			Deal: types.EncryptedDeal{
				DhKey:     []byte("dh2"),
				Signature: []byte("sig2"),
				Nonce:     []byte("nonce2"),
				Cipher:    []byte("cipher2"),
			},
		},
	}

	k.EnqueueDeals(deals)

	moreDeals := []types.Deal{
		{
			Index:          3,
			RecipientIndex: 30,
			Signature:      []byte("sig3"),
			Deal: types.EncryptedDeal{
				DhKey:     []byte("dh3"),
				Signature: []byte("sig3"),
				Nonce:     []byte("nonce3"),
				Cipher:    []byte("cipher3"),
			},
		},
	}
	k.EnqueueDeals(moreDeals)

	got := k.DequeueDeals(3)
	require.Len(t, got, 3)
	require.Equal(t, uint32(1), got[0].Index)
	require.Equal(t, uint32(2), got[1].Index)
	require.Equal(t, uint32(3), got[2].Index)

	got = k.DequeueDeals(3)
	require.Len(t, got, 0)
}

func TestEnqueueAndDequeueResponses(t *testing.T) {
	k, _ := setupDKGKeeper(t)

	responses := []types.Response{
		{
			Index: 1,
			VssResponse: &types.VSSResponse{
				SessionId: []byte("sess1"),
				Index:     1,
				Status:    true,
				Signature: []byte("sig1"),
			},
		},
		{
			Index: 2,
			VssResponse: &types.VSSResponse{
				SessionId: []byte("sess2"),
				Index:     2,
				Status:    false,
				Signature: []byte("sig2"),
			},
		},
	}

	k.EnqueueResponses(responses)

	moreResponses := []types.Response{
		{
			Index: 3,
			VssResponse: &types.VSSResponse{
				SessionId: []byte("sess3"),
				Index:     3,
				Status:    true,
				Signature: []byte("sig3"),
			},
		},
	}
	k.EnqueueResponses(moreResponses)

	got := k.DequeueResponses(3)
	require.Len(t, got, 3)
	require.Equal(t, uint32(1), got[0].Index)
	require.Equal(t, uint32(2), got[1].Index)
	require.Equal(t, uint32(3), got[2].Index)

	got = k.DequeueResponses(3)
	require.Len(t, got, 0)
}
