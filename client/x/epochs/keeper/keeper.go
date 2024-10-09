package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/piplabs/story/client/x/epochs/types"
)

type Keeper struct {
	cdc   codec.BinaryCodec
	hooks types.EpochHooks

	// NOTE(Narangde): Add storeService and EventService manually instead of Environment
	storeService store.KVStoreService
	EventService event.Service

	Schema    collections.Schema
	EpochInfo collections.Map[string, types.EpochInfo]
}

// NewKeeper returns a new keeper by codec and storeKey inputs.
func NewKeeper(storeService store.KVStoreService, eventService event.Service, cdc codec.BinaryCodec) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		storeService: storeService,
		EventService: eventService,
		EpochInfo:    collections.NewMap(sb, types.KeyPrefixEpoch, "epoch_info", collections.StringKey, codec.CollValue[types.EpochInfo](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return &k
}

// Set hooks.
func (k *Keeper) SetHooks(eh types.EpochHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set epochs hooks twice")
	}

	k.hooks = eh

	return k
}
