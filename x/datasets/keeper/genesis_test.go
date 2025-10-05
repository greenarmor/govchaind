package keeper_test

import (
	"testing"

	"govchain/x/datasets/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:     types.DefaultParams(),
		EntryList:  []types.Entry{{Id: 0}, {Id: 1}},
		EntryCount: 2,
	}
	f := initFixture(t)
	err := f.keeper.InitGenesis(f.ctx, genesisState)
	require.NoError(t, err)
	got, err := f.keeper.ExportGenesis(f.ctx)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.EqualExportedValues(t, genesisState.Params, got.Params)
	require.EqualExportedValues(t, genesisState.EntryList, got.EntryList)
	require.Equal(t, genesisState.EntryCount, got.EntryCount)

}
