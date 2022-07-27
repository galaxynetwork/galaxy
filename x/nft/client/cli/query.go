package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/galaxies-labs/galaxy/x/nft/types"
	"github.com/spf13/cobra"
)

const (
	FlagOwner   = "owner"
	FlagBrandID = "brand-id"
	FlagClassID = "class-id"
)

// NewTxCmd returns a root CLI command handler for all x/nft transaction commands.
func GetQueryCmd() *cobra.Command {
	nftQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the nft moudle",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nftQueryCmd.AddCommand(
		GetCmdQueryClasses(),
		GetCmdQueryClass(),
		GetCmdQueryNFTs(),
		GetCmdQueryNFT(),
		GetCmdQuerOwner(),
		GetCmdSupplyClass(),
	)

	return nftQueryCmd
}

func GetCmdQueryClasses() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "classes",
		Short: "Query classes with optional filters",
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

			brandID, err := cmd.Flags().GetString(FlagBrandID)
			if err != nil {
				return err
			}

			result, err := queryClient.Classes(cmd.Context(), &types.QueryClassesRequest{
				Pagination: pageReq,
				BrandId:    brandID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "classes")

	cmd.Flags().String(FlagBrandID, "", "(optional) filter classes by brandID")
	return cmd
}

func GetCmdQueryClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "class [brand-id] [class-id]",
		Short: "Query class based on it's brand and class id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			result, err := queryClient.Class(cmd.Context(), &types.QueryClassRequest{
				BrandId: args[0],
				ClassId: args[1],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQueryNFTs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nfts",
		Short: "Query nfts with optional filters",
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

			brandID, err := cmd.Flags().GetString(FlagBrandID)
			if err != nil {
				return err
			}

			classID, err := cmd.Flags().GetString(FlagClassID)
			if err != nil {
				return err
			}

			owner, err := cmd.Flags().GetString(FlagOwner)
			if err != nil {
				return err
			}

			result, err := queryClient.NFTs(cmd.Context(), &types.QueryNFTsRequest{
				Pagination: pageReq,
				BrandId:    brandID,
				ClassId:    classID,
				Owner:      owner,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "nfts")

	cmd.Flags().String(FlagBrandID, "", "(optional) filter nfts by brandID, requires a classID")
	cmd.Flags().String(FlagClassID, "", "(optional) filter nfts by classID, required a brandID")
	cmd.Flags().String(FlagOwner, "", "(optional) filter nfts by owner address, bech32_address")
	return cmd
}

func GetCmdQueryNFT() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft [brand-id] [class-id] [id]",
		Short: "Query nft based on it's brand and class and nft id",
		Long: `"Query nft based on it's brand and class and nft id
		Note, the 'id' argument is nftId of type uint64.
		`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			nftID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			result, err := queryClient.NFT(cmd.Context(), &types.QueryNFTRequest{
				BrandId: args[0],
				ClassId: args[1],
				Id:      nftID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdQuerOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner [brand-id] [class-id] [id]",
		Short: "Query owner of nft based on it's brand and class and nft id",
		Long: `"Query owner of nft based on it's brand and class and nft id
		Note, the 'id' argument is nftId of type uint64.
		`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			nftID, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			result, err := queryClient.Owner(cmd.Context(), &types.QueryOwnerRequest{
				BrandId: args[0],
				ClassId: args[1],
				Id:      nftID,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetCmdSupplyClass() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supply [brand-id] [class-id]",
		Short: "Query number of nfts and sequence from class based on it's brand and class id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			result, err := queryClient.Supply(cmd.Context(), &types.QuerySupplyRequest{
				BrandId: args[0],
				ClassId: args[1],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(result)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
