package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/galaxies-labs/galaxy/x/brand/types"
	"github.com/spf13/cobra"
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
		GetCmdQueryBrandsByOwner(),
	)

	return brandQueryCmd
}

func GetCmdQueryBrands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brands",
		Short: "Query for all brands",
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

			result, err := queryClient.Brands(cmd.Context(), &types.QueryBrandsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "brands")

	return cmd
}

func GetCmdQueryBrand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brand [brand_id]",
		Short: "Query a brand",
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

func GetCmdQueryBrandsByOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brands-by-owner [owner_addr_bech32]",
		Short: "Query for all brands by owner",
		Args:  cobra.ExactArgs(1),
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

			result, err := queryClient.BrandsByOwner(cmd.Context(), &types.QueryBrandsByOwnerRequest{
				Owner:      args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "brands-by-owner")

	return cmd
}
