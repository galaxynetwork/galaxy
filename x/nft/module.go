package nft

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	codec "github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/nft/keeper"
	"github.com/galaxies-labs/galaxy/x/nft/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {
	cdc codec.Codec
}

func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState

	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return genState.Validate()
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {}

func (AppModuleBasic) GetTxCmd() *cobra.Command { return &cobra.Command{} }

func (AppModuleBasic) GetQueryCmd() *cobra.Command { return &cobra.Command{} }

type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

func NewAppModule(cdc codec.Codec, keeper keeper.Keeper) AppModule {
	return AppModule{AppModuleBasic: NewAppModuleBasic(cdc), keeper: keeper}
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, raw json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	cdc.MustUnmarshalJSON(raw, &genState)

	//am.keeper.InitGenesis(ctx, genState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	//genState := am.keeper.ExportGenesis(ctx)
	//	cdc.MustMarshalJSON(&genState)
	return nil
}

func (am AppModule) RegisterInvariants(sdk.InvariantRegistry) {}

func (am AppModule) Route() sdk.Route { return sdk.Route{} }

func (am AppModule) QuerierRoute() string { return types.QuerierRoute }

func (am AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier { return nil }

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	querier := keeper.Querier{Keeper: am.keeper}
	types.RegisterQueryServer(cfg.QueryServer(), querier)
}

func (am AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
