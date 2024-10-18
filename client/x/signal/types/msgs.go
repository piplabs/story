package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
)

const (
	URLMsgScheduleUpgrade = "/client.x.signal.Msg/ScheduleUpgrade"
)

var (
	_ sdk.Msg            = &MsgScheduleUpgrade{}
	_ legacytx.LegacyMsg = &MsgScheduleUpgrade{}
)

var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

func NewMsgScheduleUpgrade(upgrade Upgrade) *MsgScheduleUpgrade {
	return &MsgScheduleUpgrade{
		Upgrade: &upgrade,
	}
}

// GetSignBytes implements legacytx.LegacyMsg.
func (msg *MsgScheduleUpgrade) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// Route implements legacytx.LegacyMsg.
func (msg *MsgScheduleUpgrade) Route() string {
	return RouterKey
}

// Type implements legacytx.LegacyMsg.
func (msg *MsgScheduleUpgrade) Type() string {
	return URLMsgScheduleUpgrade
}
