package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// registers the necessary x/nft concrete types
func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateClass{}, "galaxy/nft/create-class", nil)
	cdc.RegisterConcrete(&MsgEditClass{}, "galaxy/nft/edit-class", nil)

	cdc.RegisterConcrete(&MsgMintToNFT{}, "galaxy/nft/mint-to-nft", nil)
	cdc.RegisterConcrete(&MsgBurnNFT{}, "galaxy/nft/burn-nft", nil)
	cdc.RegisterConcrete(&MsgTransferNFT{}, "galaxy/nft/transfer-nft", nil)
	cdc.RegisterConcrete(&MsgUpdateNFT{}, "galaxy/nft/update-nft", nil)

}

// RegisterInterfaces registers the x/nft interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateClass{},
		&MsgEditClass{},

		&MsgMintToNFT{},
		&MsgBurnNFT{},
		&MsgTransferNFT{},
		&MsgUpdateNFT{},
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
