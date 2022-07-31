package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/galaxies-labs/galaxy/x/nft/types"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	FlagName           = "name"
	FlagDetails        = "details"
	FlagFeeBasisPoints = "fee-basis-points"
	FlagExternalUrl    = "external-url"
	FlagImageUri       = "image-uri"
	FlagUri            = "uri"
	FlagVarUri         = "var-uri"
)

func NewTxCmd() *cobra.Command {
	nftTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "NFT transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nftTxCmd.AddCommand(
		NewCreateClassTxCmd(),
		NewEditClassTxCmd(),
		NewMintNFTTxCmd(),
		NewUpdateNFTTxCmd(),
		NewTransferNFTTxCmd(),
		NewBurnNFTTxCmd(),
	)

	return nftTxCmd
}

func NewCreateClassTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-class [brand-id] [id]  --fee-basis-points [uint32,optional] --name [text,optional] --details [text,optional] --external-url [text,optional] --image-uri [text,optional]",
		Short: "Create a new class within brand",
		Long: `"Create a new class within brand"
Note, the 'brand-id' argument is id of the brand to create class.
the 'id' argument is id of the class to be created.
the '--fee-basis-points' flag is commission to be received nft of a class is traded. (0~10,000)
the '--from' flag is owner of brand.
the '--name' flag is name of class.
the '--details' flag is details of class.
the '--image-uri' flag is representative image uri of class.
the '--external-url' flag is website url of class.
`, Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewBuildCreateClassMsg(args[0], args[1], clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().Uint32(FlagFeeBasisPoints, 0, "The class (optional) fee basis points")
	cmd.Flags().String(FlagName, "", "The class (optional) name")
	cmd.Flags().String(FlagDetails, "", "The class (optional) details")
	cmd.Flags().String(FlagExternalUrl, "", "The class (optional) external url")
	cmd.Flags().String(FlagImageUri, "", "The class (optional) image uri")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}

func NewEditClassTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-class [brand-id] [id]  --fee-basis-points [uint32,optional] --name [text,optional] --details [text,optional] --external-url [text,optional] --image-uri [text,optional]",
		Short: "Edit an existing class within brand",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewBuildEditClassMsg(args[0], args[1], clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().Uint32(FlagFeeBasisPoints, types.DoNotModifyFeeBasisPoints, "The class (optional) fee basis points")
	cmd.Flags().String(FlagName, types.DoNotModifyDesc, "The class (optional) name")
	cmd.Flags().String(FlagDetails, types.DoNotModifyDesc, "The class (optional) details")
	cmd.Flags().String(FlagExternalUrl, types.DoNotModifyDesc, "The class (optional) external url")
	cmd.Flags().String(FlagImageUri, types.DoNotModifyDesc, "The class (optional) image uri")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	return cmd
}

func NewMintNFTTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [brand-id] [class-id] [recipient,bech32_address] --uri [text] --var-uri [text,optional]",
		Short: "Mint a new nft within class",
		Long: `"Mint a new nft within class"
Note, the '--from' flag is owner of brand (minter).
the '--uri' flag is url of metadata stored off chain.
the '--var-uri' flag is url of data wanted by the owner stored off chain.
`, Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewBuildMintNFTMsg(args[0], args[1], args[2], clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagUri, "", "Uri of nft metadata stored off chain")
	cmd.Flags().String(FlagVarUri, "", "(optional) Uri of nft data wanted by the owner stored off chain")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagUri)
	return cmd
}

func NewUpdateNFTTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [brand-id] [class-id] [nft-id] --var-uri [text] --from [sender]",
		Short: "Update var-uri of an existing nft",
		Long: `"Update var-uri of an existing nft"
		Note, the '--from' flag is owner of nft.
		the '--var-uri' flag is url of data wanted by the owner stored off chain.
		`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			txf, msg, err := NewBuildUpdateNFTMsg(args[0], args[1], args[2], clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().String(FlagVarUri, "", "Uri of nft data wanted by the owner stored off chain")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(FlagVarUri)
	return cmd
}

func NewTransferNFTTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [brand-id] [class-id] [nft-id] [recipient,bech32_address] --from [sender]",
		Short: "Transfer ownership of existing nft",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			nftID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferNFT(args[0], args[1], nftID, clientCtx.GetFromAddress().String(), args[3])
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

func NewBurnNFTTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [brand-id] [class-id] [nft-id] --from [sender]",
		Short: "Burn existing nft",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(clientCtx, cmd.Flags()).WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

			nftID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgBurnNFT(args[0], args[1], nftID, clientCtx.GetFromAddress().String())
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

func NewBuildCreateClassMsg(brandID string, classID string, clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	name, _ := fs.GetString(FlagName)
	details, _ := fs.GetString(FlagDetails)
	feeBasisPoints, _ := fs.GetUint32(FlagFeeBasisPoints)
	externalUrl, _ := fs.GetString(FlagExternalUrl)
	imageUri, _ := fs.GetString(FlagImageUri)

	return txf,
		types.NewMsgCreateClass(brandID, classID, clientCtx.GetFromAddress().String(),
			feeBasisPoints, types.NewClassDescription(name, details, externalUrl, imageUri)), nil
}

func NewBuildEditClassMsg(brandID string, classID string, clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	name, _ := fs.GetString(FlagName)
	details, _ := fs.GetString(FlagDetails)
	feeBasisPoints, _ := fs.GetUint32(FlagFeeBasisPoints)
	externalUrl, _ := fs.GetString(FlagExternalUrl)
	imageUri, _ := fs.GetString(FlagImageUri)

	return txf,
		types.NewMsgEditClass(brandID, classID, clientCtx.GetFromAddress().String(),
			feeBasisPoints, types.NewClassDescription(name, details, externalUrl, imageUri)), nil
}

func NewBuildMintNFTMsg(brandID string, classID string, recipient string, clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	uri, _ := fs.GetString(FlagUri)
	varUri, _ := fs.GetString(FlagVarUri)

	return txf,
		types.NewMsgMintNFT(brandID, classID, uri, varUri, clientCtx.GetFromAddress().String(), recipient), nil
}

func NewBuildUpdateNFTMsg(brandID string, classID string, id string, clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, sdk.Msg, error) {
	varUri, _ := fs.GetString(FlagVarUri)

	nftID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return txf, nil, err
	}
	return txf,
		types.NewMsgUpdateNFT(brandID, classID, nftID, varUri, clientCtx.GetFromAddress().String()), nil
}
