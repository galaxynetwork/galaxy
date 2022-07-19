package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	flag "github.com/spf13/pflag"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/spf13/cobra"
)

const (
	FlagName          = "name"
	FlagDetails       = "details"
	FlagBrandImageUri = "brand-image-uri"
)

// NewTxCmd returns a root CLI command handler for all x/staking transaction commands.
func NewTxCmd() *cobra.Command {
	brandTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Brand transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	brandTxCmd.AddCommand(
		NewCreateBrandTxCmd(),
		NewEditBrandTxCmd(),
		NewTransferOwnershipBrandTxCmd(),
	)

	return brandTxCmd
}

func NewCreateBrandTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-brand [brand_id] --name [text] --details [text,optional] --brand-image-uri [text,optional]",
		Short: "create a new brand",
		Long: `"Create brand with brandID which is a trademark
Note, the '--from' flag is owner of brand.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewBuildCreateBrandMsg(args[0], clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagName, "", "The brand's name")
	cmd.Flags().String(FlagDetails, "", "The brand's (optional) details")
	cmd.Flags().String(FlagBrandImageUri, "", "The brand's (optional) image uri")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagName)

	return cmd
}

func NewEditBrandTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-brand [brand_id] --name [text] --details [text,optional] --brand-image-uri [text,optional]",
		Short: "edit an existing brand",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewBuildEditBrandMsg(args[0], clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagName, types.DoNotModify, "The brand's name")
	cmd.Flags().String(FlagDetails, types.DoNotModify, "The brand's (optional) details")
	cmd.Flags().String(FlagBrandImageUri, types.DoNotModify, "The brand's (optional) image uri")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagName)

	return cmd
}

func NewTransferOwnershipBrandTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [brand_id] [destowner_addr_bech32]",
		Short: "transfer ownership an existing brand",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			msg := types.NewMsgTransferOwnershipBrand(args[0], clientCtx.GetFromAddress().String(), args[1])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func NewBuildEditBrandMsg(brandID string, clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	name, _ := fs.GetString(FlagName)
	details, _ := fs.GetString(FlagDetails)
	brandImageUri, _ := fs.GetString(FlagBrandImageUri)

	return txf, types.NewMsgEditBrand(brandID, clientCtx.GetFromAddress().String(), types.NewBrandDescription(name, details, brandImageUri)), nil
}

func NewBuildCreateBrandMsg(brandID string, clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	name, _ := fs.GetString(FlagName)
	details, _ := fs.GetString(FlagDetails)
	brandImageUri, _ := fs.GetString(FlagBrandImageUri)

	return txf, types.NewMsgCreateBrand(brandID, clientCtx.GetFromAddress().String(), types.NewBrandDescription(name, details, brandImageUri)), nil
}
