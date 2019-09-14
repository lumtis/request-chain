package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ltacker/request-chain/x/requestchain/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	requestchainQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the requestchain module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	requestchainQueryCmd.AddCommand(client.GetCommands(
		GetCmdGetBlock(storeKey, cdc),
	)...)
	return requestchainQueryCmd
}

// GetCmdResolveName queries information about a name
func GetCmdGetBlock(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [index]",
		Short: "get index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			index := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/getblock/%s", queryRoute, index), nil)
			if err != nil {
				fmt.Printf("could not get block - %s \nError: %s \n", index, err)
				return nil
			}

			var out types.QueryGetBlock
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
