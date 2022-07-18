package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateBrand(t *testing.T) {
	tests := []struct {
		c          string
		msg        *MsgCreateBrand
		expectPass bool
	}{
		{"valid", NewMsgCreateBrand("brandid", sdk.AccAddress("pubkey...").String(), NewBrandDescription("name", "", "")), true},
		{"invalid description", NewMsgCreateBrand("brandid", sdk.AccAddress("pubkey...").String(), NewBrandDescription("", "", "")), false},
		{"invalid owner", NewMsgCreateBrand("brandid", "invalidbech32", NewBrandDescription("name", "", "")), false},
		{"invalid id", NewMsgCreateBrand("_wdwd", sdk.AccAddress("pubkey...").String(), NewBrandDescription("name", "", "")), false},
	}

	for _, test := range tests {
		if test.expectPass {
			require.NoError(t, test.msg.ValidateBasic(), "test: %s", test.c)
		} else {
			require.Error(t, test.msg.ValidateBasic(), "test: %s", test.c)
		}
	}
}

func TestMsgEditBrand(t *testing.T) {
	tests := []struct {
		c          string
		msg        *MsgEditBrand
		expectPass bool
	}{
		{"valid", NewMsgEditBrand("brandid", sdk.AccAddress("pubkey...").String(), NewBrandDescription("name", "", "")), true},
		{"invalid description", NewMsgEditBrand("brandid", sdk.AccAddress("pubkey...").String(), NewBrandDescription("", "", "")), false},
		{"invalid owner", NewMsgEditBrand("brandid", "invalidbech32", NewBrandDescription("name", "", "")), false},
		{"invalid id", NewMsgEditBrand("_wdwd", sdk.AccAddress("pubkey...").String(), NewBrandDescription("name", "", "")), false},
	}

	for _, test := range tests {
		if test.expectPass {
			require.NoError(t, test.msg.ValidateBasic(), "test: %s", test.c)
		} else {
			require.Error(t, test.msg.ValidateBasic(), "test: %s", test.c)
		}
	}
}

func TestMsgTransferSwapBrand(t *testing.T) {
	tests := []struct {
		c          string
		msg        *MsgTransferOwnershipBrand
		expectPass bool
	}{
		{"valid a", NewMsgTransferOwnershipBrand("brandid", sdk.AccAddress("pubkey...").String(), sdk.AccAddress("pubkey.2..").String()), true},
		{"valid b", NewMsgTransferOwnershipBrand("brandid", sdk.AccAddress("pubkey2...").String(), sdk.AccAddress("pubkey..").String()), true},
		{"invalid id a", NewMsgTransferOwnershipBrand("_brandid", sdk.AccAddress("pubkey2...").String(), sdk.AccAddress("pubkey..").String()), false},
		{"invalid id b", NewMsgTransferOwnershipBrand("", sdk.AccAddress("pubkey2...").String(), sdk.AccAddress("pubkey..").String()), false},
		{"invalid owner a", NewMsgTransferOwnershipBrand("brandid", "bech32", sdk.AccAddress("pubkey..").String()), false},
		{"invalid owner b", NewMsgTransferOwnershipBrand("brandid", "", sdk.AccAddress("pubkey..").String()), false},
		{"invalid dest owner a", NewMsgTransferOwnershipBrand("brandid", sdk.AccAddress("pubkey..").String(), "bech32"), false},
		{"invalid dest owner b", NewMsgTransferOwnershipBrand("brandid", sdk.AccAddress("pubkey..").String(), ""), false},
		{"invalid all owner", NewMsgTransferOwnershipBrand("brandid", "", ""), false},
	}

	for _, test := range tests {
		if test.expectPass {
			require.NoError(t, test.msg.ValidateBasic(), "test: %s", test.c)
		} else {
			require.Error(t, test.msg.ValidateBasic(), "test: %s", test.c)
		}
	}
}
