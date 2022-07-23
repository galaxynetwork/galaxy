package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		err := test.msg.ValidateBasic()
		if test.expectPass {
			require.NoError(t, err, "test for index: %d", index)
		} else {
			require.Error(t, err, "test for index: %d", index)

		}
	}
}
