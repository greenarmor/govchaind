package datasets

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"govchain/testutil/sample"
	datasetssimulation "govchain/x/datasets/simulation"
	"govchain/x/datasets/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	datasetsGenesis := types.GenesisState{
		Params:    types.DefaultParams(),
		EntryList: []types.Entry{{Id: 0, Creator: sample.AccAddress()}, {Id: 1, Creator: sample.AccAddress()}}, EntryCount: 2,
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&datasetsGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateEntry          = "op_weight_msg_datasets"
		defaultWeightMsgCreateEntry int = 100
	)

	var weightMsgCreateEntry int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateEntry, &weightMsgCreateEntry, nil,
		func(_ *rand.Rand) {
			weightMsgCreateEntry = defaultWeightMsgCreateEntry
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateEntry,
		datasetssimulation.SimulateMsgCreateEntry(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateEntry          = "op_weight_msg_datasets"
		defaultWeightMsgUpdateEntry int = 100
	)

	var weightMsgUpdateEntry int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateEntry, &weightMsgUpdateEntry, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateEntry = defaultWeightMsgUpdateEntry
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateEntry,
		datasetssimulation.SimulateMsgUpdateEntry(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteEntry          = "op_weight_msg_datasets"
		defaultWeightMsgDeleteEntry int = 100
	)

	var weightMsgDeleteEntry int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteEntry, &weightMsgDeleteEntry, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteEntry = defaultWeightMsgDeleteEntry
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteEntry,
		datasetssimulation.SimulateMsgDeleteEntry(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
