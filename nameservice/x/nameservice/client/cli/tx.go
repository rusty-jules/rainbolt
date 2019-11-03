package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/sdk-tutorials/nameservice/x/nameservice/internal/types"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	nameserviceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Nameservice transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nameserviceTxCmd.AddCommand(client.PostCommands(
		GetCmdCreateOrder(cdc),
	)...)

	return nameserviceTxCmd
}

func GetCmdCreateOrder(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-order [channel_state] [channel_token] [amount]",
		Short: "create an order and place funds in escrow",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			// var blsPubKey, commitPubKey types.Bls12381PubKey // [96]byte
			// copy(blsPubKey[:], []byte(args[0]))
			// copy(commitPubKey[:], []byte(args[1]))

			msg := types.NewMsgCreateOrder(cliCtx.GetFromAddress(), args[0], args[1], coins)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// func GetCmdFillOrder(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use: "fill-order "
// 	},
// }
