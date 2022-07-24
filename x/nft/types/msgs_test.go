package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/galaxies-labs/galaxy/internal/util"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateClass(t *testing.T) {
	tests := []struct {
		expectPass bool
		msg        *MsgCreateClass
	}{
		{true, NewMsgCreateClass("brandid", "classid", sdk.AccAddress("owner1...").String(), 10_000, NewClassDescription("name", "", "", ""))},
		{true, NewMsgCreateClass("brandid", "class-id123", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},
		{true, NewMsgCreateClass("brandid", "cla", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},

		//invliad brandID
		{false, NewMsgCreateClass("-", "class-id123", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},
		//invliad classID
		{false, NewMsgCreateClass("brandid", "", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},
		{false, NewMsgCreateClass("brandid", "-classid", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},
		{false, NewMsgCreateClass("brandid", "ci", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},
		//invliad ownerAddress
		{false, NewMsgCreateClass("brandid", "classid", "owner1...", 1, NewClassDescription("name", "", "", ""))},
		//invliad feeBasisPoints
		{false, NewMsgCreateClass("brandid", "classid", sdk.AccAddress("owner1...").String(), 0, NewClassDescription("name", "", "", ""))},
		{false, NewMsgCreateClass("brandid", "classid", sdk.AccAddress("owner1...").String(), MaxFeeBasisPoints+1, NewClassDescription("name", "", "", ""))},
	}

	for index, test := range tests {
		require.Equal(t, test.msg.Type(), TypeMsgCreateClass)
		err := test.msg.ValidateBasic()
		if test.expectPass {
			require.NoError(t, err, "test for index: %d", index)
		} else {
			require.Error(t, err, "test for index: %d", index)

		}
	}
}

func TestMsgEditClass(t *testing.T) {
	tests := []struct {
		expectPass bool
		msg        *MsgEditClass
	}{
		{true, NewMsgEditClass("brandid", "classid", sdk.AccAddress("owner1...").String(), 10_000, NewClassDescription("name", "", "", ""))},
		{true, NewMsgEditClass("brandid", "class-id123", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},
		{true, NewMsgEditClass("brandid", "cla", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},

		//invliad brandID
		{false, NewMsgEditClass("-", "class-id123", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},
		//invliad classID
		{false, NewMsgEditClass("brandid", "", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},
		{false, NewMsgEditClass("brandid", "-classid", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},
		{false, NewMsgEditClass("brandid", "ci", sdk.AccAddress("owner1...").String(), 1, NewClassDescription("name", "", "", ""))},
		//invliad ownerAddress
		{false, NewMsgEditClass("brandid", "classid", "owner1...", 1, NewClassDescription("name", "", "", ""))},
		//invliad feeBasisPoints
		{false, NewMsgEditClass("brandid", "classid", sdk.AccAddress("owner1...").String(), 0, NewClassDescription("name", "", "", ""))},
		{false, NewMsgEditClass("brandid", "classid", sdk.AccAddress("owner1...").String(), MaxFeeBasisPoints+1, NewClassDescription("name", "", "", ""))},
	}

	for index, test := range tests {
		require.Equal(t, test.msg.Type(), TypeMsgEditClass)
		err := test.msg.ValidateBasic()
		if test.expectPass {
			require.NoError(t, err, "test for index: %d", index)
		} else {
			require.Error(t, err, "test for index: %d", index)

		}
	}
}

func TestMsgMintNFT(t *testing.T) {
	tests := []struct {
		expectPass bool
		msg        *MsgMintNFT
	}{
		{true, NewMsgMintNFT("brandid", "classid", "ipfs://hash", "https://nft.json",
			sdk.AccAddress("minter").String(), sdk.AccAddress("recipient").String())},
		{true, NewMsgMintNFT("brandid", "classid", "ipfs://hash", "",
			sdk.AccAddress("minter").String(), sdk.AccAddress("recipient").String())},
		{true, NewMsgMintNFT("brandid", "classid", util.GenStringWithLength(MaxUriLength), util.GenStringWithLength(MaxUriLength),
			sdk.AccAddress("minter").String(), sdk.AccAddress("recipient").String())},

		{false, NewMsgMintNFT("", "classid", "ipfs://hash", "",
			sdk.AccAddress("minter").String(), sdk.AccAddress("recipient").String())},
		{false, NewMsgMintNFT("brandid", "", "ipfs://hash", "",
			sdk.AccAddress("minter").String(), sdk.AccAddress("recipient").String())},
		{false, NewMsgMintNFT("brandid", "classid", "", "",
			sdk.AccAddress("minter").String(), sdk.AccAddress("recipient").String())},
		{false, NewMsgMintNFT("brandid", "classid", "ipfs://hash", "",
			"minter", sdk.AccAddress("recipient").String())},
		{false, NewMsgMintNFT("brandid", "classid", "ipfs://hash", "",
			sdk.AccAddress("minter").String(), "recipient")},
	}

	for index, test := range tests {
		require.Equal(t, test.msg.Type(), TypeMsgMintNFT)
		err := test.msg.ValidateBasic()
		if test.expectPass {
			require.NoError(t, err, "test for index: %d", index)
		} else {
			require.Error(t, err, "test for index: %d", index)

		}
	}
}

func TestMsgBurnNFT(t *testing.T) {
	tests := []struct {
		expectPass bool
		msg        *MsgBurnNFT
	}{
		{true, NewMsgBurnNFT("brandid", "classid", 1, sdk.AccAddress("sender").String())},
		{true, NewMsgBurnNFT("brandid", "classid", 2, sdk.AccAddress("sender2").String())},

		{false, NewMsgBurnNFT("brandid", "classid", 0, sdk.AccAddress("sender").String())},
		{false, NewMsgBurnNFT("", "classid", 1, sdk.AccAddress("sender").String())},
		{false, NewMsgBurnNFT("brandid", "", 1, sdk.AccAddress("sender").String())},
		{false, NewMsgBurnNFT("brandid", "classid", 1, "sender")},
	}

	for index, test := range tests {
		require.Equal(t, test.msg.Type(), TypeMsgBurnNFT)
		err := test.msg.ValidateBasic()
		if test.expectPass {
			require.NoError(t, err, "test for index: %d", index)
		} else {
			require.Error(t, err, "test for index: %d", index)

		}
	}
}

func TestMsgTransferNFT(t *testing.T) {
	tests := []struct {
		expectPass bool
		msg        *MsgTransferNFT
	}{
		{true, NewMsgTransferNFT("brandid", "classid", 1,
			sdk.AccAddress("sender").String(), sdk.AccAddress("recipient").String())},

		{false, NewMsgTransferNFT("brandid", "classid", 0,
			sdk.AccAddress("sender").String(), sdk.AccAddress("recipient").String())},
		{false, NewMsgTransferNFT("", "classid", 1,
			sdk.AccAddress("sender").String(), sdk.AccAddress("recipient").String())},
		{false, NewMsgTransferNFT("brandid", "", 1,
			sdk.AccAddress("sender").String(), sdk.AccAddress("recipient").String())},
		{false, NewMsgTransferNFT("brandid", "classid", 1,
			"sender", sdk.AccAddress("recipient").String())},
		{false, NewMsgTransferNFT("brandid", "classid", 1,
			sdk.AccAddress("sender").String(), "recipient")},
	}

	for index, test := range tests {
		require.Equal(t, test.msg.Type(), TypeMsgTransferNFT)
		err := test.msg.ValidateBasic()
		if test.expectPass {
			require.NoError(t, err, "test for index: %d", index)
		} else {
			require.Error(t, err, "test for index: %d", index)

		}
	}
}

func TestMsgUpdateNFT(t *testing.T) {
	tests := []struct {
		expectPass bool
		msg        *MsgUpdateNFT
	}{
		{true, NewMsgUpdateNFT("brandid", "classid", 1,
			"", sdk.AccAddress("sender").String())},
		{true, NewMsgUpdateNFT("brandid", "class-id", 1,
			"https://nft.json", sdk.AccAddress("sender2").String())},

		{false, NewMsgUpdateNFT("brandid", "classid", 0,
			"", sdk.AccAddress("sender").String())},
		{false, NewMsgUpdateNFT("", "classid", 1,
			"", sdk.AccAddress("sender").String())},
		{false, NewMsgUpdateNFT("brandid", "", 1,
			"", sdk.AccAddress("sender").String())},
		{false, NewMsgUpdateNFT("brandid", "classid", 1,
			"", "sender")},
	}

	for index, test := range tests {
		require.Equal(t, test.msg.Type(), TypeMsgUpdateNFT)
		err := test.msg.ValidateBasic()
		if test.expectPass {
			require.NoError(t, err, "test for index: %d", index)
		} else {
			require.Error(t, err, "test for index: %d", index)

		}
	}
}
