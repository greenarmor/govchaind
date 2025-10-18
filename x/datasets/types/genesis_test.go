package types_test

import (
	"testing"

	"govchain/x/datasets/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				EntryList: []types.Entry{
					{
						Id:                0,
						Status:            types.EntryStatus_ENTRY_STATUS_PENDING,
						RequiredApprovals: types.DefaultParams().MinRequiredApprovals,
						LiabilityContract: "ipfs://contract",
					},
					{
						Id:                1,
						Status:            types.EntryStatus_ENTRY_STATUS_PENDING,
						RequiredApprovals: types.DefaultParams().MinRequiredApprovals,
						LiabilityContract: "ipfs://contract",
					},
				},
				EntryCount: 2,
			},
			valid: true,
		}, {
			desc: "duplicated entry",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				EntryList: []types.Entry{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		}, {
			desc: "invalid entry count",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				EntryList: []types.Entry{
					{
						Id: 1,
					},
				},
				EntryCount: 0,
			},
			valid: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
