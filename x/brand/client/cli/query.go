package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/spf13/cobra"
)

const (
	FlagOwner = "owner"
)

// NewTxCmd returns a root CLI command handler for all x/staking transaction commands.
func GetQueryCmd() *cobra.Command {
	brandQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the brand moudle",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	brandQueryCmd.AddCommand(
		GetCmdQueryBrands(),
		GetCmdQueryBrand(),
	)

	return brandQueryCmd
}

func GetCmdQueryBrands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brands",
		Short: "Query brands with optional filters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			owner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			result, err := queryClient.Brands(cmd.Context(), &types.QueryBrandsRequest{
				Pagination: pageReq,
				Owner:      owner,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "brands")

	cmd.Flags().String(FlagOwner, "", "(optional) filter brands by owner address, bech32_address")
	return cmd
}

func GetCmdQueryBrand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brand [brand-id]",
		Short: "Query details of a single brand",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			result, err := queryClient.Brand(cmd.Context(), &types.QueryBrandRequest{
				BrandId: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&result.Brand)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
