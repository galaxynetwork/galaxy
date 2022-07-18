package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// "galaxy/brand/create-brand" registers the necessary x/brand concrete types
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateBrand{}, "galaxy/brand/create-brand", nil)
	cdc.RegisterConcrete(&MsgEditBrand{}, "galaxy/brand/edit-brand", nil)
	cdc.RegisterConcrete(&MsgTransferOwnershipBrand{}, "galaxy/brand/transfer-ownership-brand", nil)
}

// RegisterInterfaces registers the x/brand interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateBrand{},
		&MsgEditBrand{},
		&MsgTransferOwnershipBrand{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)
}
