package datu_examples

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"

	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"

	accountabilitykeeper "govchain/x/accountabilityscores/keeper"
	accountabilitytypes "govchain/x/accountabilityscores/types"
	budgetkeeper "govchain/x/budgetregistry/keeper"
	budgettypes "govchain/x/budgetregistry/types"
	disbursementkeeper "govchain/x/disbursementtracker/keeper"
	disbursementtypes "govchain/x/disbursementtracker/types"
	governancekeeper "govchain/x/governancevoting/keeper"
	governancetypes "govchain/x/governancevoting/types"
	procurementkeeper "govchain/x/procurementledger/keeper"
	procurementtypes "govchain/x/procurementledger/types"
	wasmkeeper "govchain/x/wasm/keeper"
	wasmtypes "govchain/x/wasm/types"
)

type budgetKeeperBridge struct {
	keeper *budgetkeeper.Keeper
}

func (b budgetKeeperBridge) GetBudget(ctx context.Context, id uint64) (budgettypes.Budget, error) {
	return b.keeper.GetBudget(ctx, id)
}

type procurementKeeperBridge struct {
	keeper *procurementkeeper.Keeper
}

func (p procurementKeeperBridge) GetProcurement(ctx context.Context, id uint64) (procurementtypes.Procurement, error) {
	return p.keeper.GetProcurement(ctx, id)
}

func newAddressCodec() address.Codec {
	bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	return addresscodec.NewBech32Codec(bech32Prefix)
}

func deriveRawAddress(seed string) []byte {
	hash := sha256.Sum256([]byte(seed))
	raw := make([]byte, 20)
	copy(raw, hash[:20])
	return raw
}

func mustBech32(codec address.Codec, raw []byte) string {
	addr, err := codec.BytesToString(raw)
	if err != nil {
		panic(fmt.Errorf("bech32 conversion failed: %w", err))
	}
	return addr
}

func TestDATUModuleOperationsExample(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig()
	addressCodec := newAddressCodec()

	keys := map[string]*storetypes.KVStoreKey{
		accountabilitytypes.StoreKey: storetypes.NewKVStoreKey(accountabilitytypes.StoreKey),
		budgettypes.StoreKey:         storetypes.NewKVStoreKey(budgettypes.StoreKey),
		disbursementtypes.StoreKey:   storetypes.NewKVStoreKey(disbursementtypes.StoreKey),
		governancetypes.StoreKey:     storetypes.NewKVStoreKey(governancetypes.StoreKey),
		procurementtypes.StoreKey:    storetypes.NewKVStoreKey(procurementtypes.StoreKey),
		wasmtypes.StoreKey:           storetypes.NewKVStoreKey(wasmtypes.StoreKey),
	}

	sdkCtx := testutil.DefaultContextWithKeys(keys, nil, nil)
	ctx := sdkCtx

	wasmAuthorityRaw := deriveRawAddress("datu-wasm-authority")
	_ = mustBech32(addressCodec, wasmAuthorityRaw)

	wasmKeeperInstance := wasmkeeper.NewKeeper(runtime.NewKVStoreService(keys[wasmtypes.StoreKey]), encCfg.Codec, addressCodec, wasmAuthorityRaw)
	budgetKeeperInstance := budgetkeeper.NewKeeper(runtime.NewKVStoreService(keys[budgettypes.StoreKey]), encCfg.Codec, addressCodec)
	budgetBridge := budgetKeeperBridge{keeper: &budgetKeeperInstance}

	procurementKeeperInstance := procurementkeeper.NewKeeper(runtime.NewKVStoreService(keys[procurementtypes.StoreKey]), encCfg.Codec, addressCodec, budgetBridge)
	procurementBridge := procurementKeeperBridge{keeper: &procurementKeeperInstance}

	disbursementKeeperInstance := disbursementkeeper.NewKeeper(runtime.NewKVStoreService(keys[disbursementtypes.StoreKey]), encCfg.Codec, addressCodec, procurementBridge)
	accountabilityKeeperInstance := accountabilitykeeper.NewKeeper(runtime.NewKVStoreService(keys[accountabilitytypes.StoreKey]), encCfg.Codec, addressCodec)
	governanceKeeperInstance := governancekeeper.NewKeeper(runtime.NewKVStoreService(keys[governancetypes.StoreKey]), encCfg.Codec, addressCodec)

	contractAddress := mustBech32(addressCodec, deriveRawAddress("datu-contract-address"))
	contractCreator := mustBech32(addressCodec, deriveRawAddress("datu-contract-creator"))
	contractAdmin := mustBech32(addressCodec, deriveRawAddress("datu-contract-admin"))

	contract, err := wasmKeeperInstance.RegisterContract(ctx, wasmtypes.Contract{
		Address:        contractAddress,
		Creator:        contractCreator,
		Admin:          contractAdmin,
		Label:          "DATU Treasury Coordinator",
		CodeIdentifier: "datu-integrated-v1",
	})
	require.NoError(t, err)
	require.NotZero(t, contract.Id)

	storedContract, err := wasmKeeperInstance.GetContract(ctx, contract.Id)
	require.NoError(t, err)
	require.Equal(t, contract.Label, storedContract.Label)

	var contractLabels []string
	err = wasmKeeperInstance.WalkContracts(ctx, func(c wasmtypes.Contract) (bool, error) {
		contractLabels = append(contractLabels, c.Label)
		return false, nil
	})
	require.NoError(t, err)
	require.Equal(t, []string{"DATU Treasury Coordinator"}, contractLabels)

	budgetCreator := mustBech32(addressCodec, deriveRawAddress("datu-budget-owner"))
	budget, err := budgetKeeperInstance.RegisterBudget(ctx, budgettypes.Budget{
		Agency:      "Department of Budget and Management",
		FiscalYear:  "2025",
		Title:       "Digital Transparency Fund",
		Amount:      sdk.NewInt64Coin("php", 50_000_000).Amount,
		Currency:    "php",
		MetadataURI: fmt.Sprintf("wasm://%s/budgets/digital-transparency", contract.Address),
		CreatedBy:   budgetCreator,
	})
	require.NoError(t, err)
	require.NotZero(t, budget.Id)
	require.Equal(t, budgettypes.BudgetStatusDraft, budget.Status)

	budget, err = budgetKeeperInstance.UpdateBudgetStatus(ctx, budget.Id, budgettypes.BudgetStatusLegislated)
	require.NoError(t, err)
	require.Equal(t, budgettypes.BudgetStatusLegislated, budget.Status)

	loadedBudget, err := budgetKeeperInstance.GetBudget(ctx, budget.Id)
	require.NoError(t, err)
	require.Equal(t, "Digital Transparency Fund", loadedBudget.Title)

	exists, err := budgetKeeperInstance.HasBudget(ctx, budget.Id)
	require.NoError(t, err)
	require.True(t, exists)

	var budgetTitles []string
	err = budgetKeeperInstance.WalkBudgets(ctx, func(b budgettypes.Budget) (bool, error) {
		budgetTitles = append(budgetTitles, b.Title)
		return false, nil
	})
	require.NoError(t, err)
	require.Equal(t, []string{"Digital Transparency Fund"}, budgetTitles)

	procurementOfficer := mustBech32(addressCodec, deriveRawAddress("datu-procurement-officer"))
	procurement, err := procurementKeeperInstance.RegisterProcurement(ctx, procurementtypes.Procurement{
		BudgetId:    budget.Id,
		Reference:   "PROC-2025-001",
		Agency:      "DBM",
		Title:       "Open Transparency Platform",
		Category:    "ICT",
		Amount:      sdk.NewInt64Coin("php", 12_500_000).Amount,
		Currency:    "php",
		Supplier:    "GovTech PH",
		MetadataURI: fmt.Sprintf("wasm://%s/procurements/open-transparency-platform", contract.Address),
		Officer:     procurementOfficer,
	})
	require.NoError(t, err)
	require.NotZero(t, procurement.Id)
	require.Equal(t, procurementtypes.ProcurementStatusPlanning, procurement.Status)

	procurement, err = procurementKeeperInstance.UpdateProcurementStatus(ctx, procurement.Id, procurementtypes.ProcurementStatusTendering)
	require.NoError(t, err)
	require.Equal(t, procurementtypes.ProcurementStatusTendering, procurement.Status)

	fetchedProcurement, err := procurementKeeperInstance.GetProcurement(ctx, procurement.Id)
	require.NoError(t, err)
	require.Equal(t, "Open Transparency Platform", fetchedProcurement.Title)

	var procurementRefs []string
	err = procurementKeeperInstance.WalkProcurements(ctx, func(p procurementtypes.Procurement) (bool, error) {
		procurementRefs = append(procurementRefs, p.Reference)
		return false, nil
	})
	require.NoError(t, err)
	require.Equal(t, []string{"PROC-2025-001"}, procurementRefs)

	disbursementPerformer := mustBech32(addressCodec, deriveRawAddress("datu-disbursement-performer"))
	disbursement, err := disbursementKeeperInstance.RegisterDisbursement(ctx, disbursementtypes.Disbursement{
		ProcurementId: procurement.Id,
		Amount:        sdk.NewInt64Coin("php", 6_000_000).Amount,
		Currency:      "php",
		EvidenceURI:   fmt.Sprintf("wasm://%s/disbursements/open-transparency-platform/initial", contract.Address),
		Notes:         "Initial milestone release",
		PerformedBy:   disbursementPerformer,
	})
	require.NoError(t, err)
	require.NotZero(t, disbursement.Id)
	require.Equal(t, disbursementtypes.DisbursementStatusScheduled, disbursement.Status)

	disbursement, err = disbursementKeeperInstance.UpdateDisbursementStatus(ctx, disbursement.Id, disbursementtypes.DisbursementStatusReleased)
	require.NoError(t, err)
	require.Equal(t, disbursementtypes.DisbursementStatusReleased, disbursement.Status)

	fetchedDisbursement, err := disbursementKeeperInstance.GetDisbursement(ctx, disbursement.Id)
	require.NoError(t, err)
	require.Equal(t, "Initial milestone release", fetchedDisbursement.Notes)

	totalDisbursed, err := disbursementKeeperInstance.GetDisbursedTotal(ctx, procurement.Id)
	require.NoError(t, err)
	require.True(t, totalDisbursed.Equal(sdk.NewInt64Coin("php", 6_000_000).Amount))

	var disbursementIDs []uint64
	err = disbursementKeeperInstance.WalkDisbursements(ctx, func(d disbursementtypes.Disbursement) (bool, error) {
		disbursementIDs = append(disbursementIDs, d.Id)
		return false, nil
	})
	require.NoError(t, err)
	require.Equal(t, []uint64{disbursement.Id}, disbursementIDs)

	scoreUpdater := mustBech32(addressCodec, deriveRawAddress("datu-score-updater"))
	scorecard, err := accountabilityKeeperInstance.UpsertScorecard(ctx, accountabilitytypes.Scorecard{
		Subject:     contract.Address,
		Metric:      "on_time_disbursement",
		Score:       88,
		Weight:      10,
		EvidenceURI: fmt.Sprintf("wasm://%s/scorecards/disbursement-timeliness", contract.Address),
		UpdatedBy:   scoreUpdater,
	})
	require.NoError(t, err)
	require.NotZero(t, scorecard.Id)

	scorecard, err = accountabilityKeeperInstance.UpsertScorecard(ctx, accountabilitytypes.Scorecard{
		Id:          scorecard.Id,
		Subject:     contract.Address,
		Metric:      "on_time_disbursement",
		Score:       92,
		Weight:      10,
		EvidenceURI: fmt.Sprintf("wasm://%s/scorecards/disbursement-timeliness", contract.Address),
		UpdatedBy:   scoreUpdater,
	})
	require.NoError(t, err)
	require.Equal(t, uint32(92), scorecard.Score)

	fetchedScorecard, err := accountabilityKeeperInstance.GetScorecard(ctx, scorecard.Id)
	require.NoError(t, err)
	require.Equal(t, scorecard.Id, fetchedScorecard.Id)

	scopedScorecard, err := accountabilityKeeperInstance.GetScorecardBySubjectMetric(ctx, contract.Address, "on_time_disbursement")
	require.NoError(t, err)
	require.Equal(t, scorecard.Id, scopedScorecard.Id)

	var scorecardIDs []uint64
	err = accountabilityKeeperInstance.WalkScorecards(ctx, func(s accountabilitytypes.Scorecard) (bool, error) {
		scorecardIDs = append(scorecardIDs, s.Id)
		return false, nil
	})
	require.NoError(t, err)
	require.Equal(t, []uint64{scorecard.Id}, scorecardIDs)

	delegator := mustBech32(addressCodec, deriveRawAddress("datu-delegator"))
	delegatee := mustBech32(addressCodec, deriveRawAddress("datu-delegatee"))

	delegation, err := governanceKeeperInstance.RegisterDelegation(ctx, governancetypes.Delegation{
		Delegator:   delegator,
		Delegatee:   delegatee,
		Scope:       fmt.Sprintf("budget:%s", contract.Address),
		ExpiresAt:   1_730_000_000,
		Active:      true,
		MetadataURI: fmt.Sprintf("wasm://%s/delegations/%s", contract.Address, delegator),
	})
	require.NoError(t, err)
	require.NotZero(t, delegation.Id)
	require.True(t, delegation.Active)

	delegation, err = governanceKeeperInstance.DeactivateDelegation(ctx, delegation.Id)
	require.NoError(t, err)
	require.False(t, delegation.Active)

	fetchedDelegation, err := governanceKeeperInstance.GetDelegation(ctx, delegation.Id)
	require.NoError(t, err)
	require.Equal(t, delegation.Id, fetchedDelegation.Id)

	scopedDelegation, err := governanceKeeperInstance.GetDelegationByScope(ctx, delegator, fmt.Sprintf("budget:%s", contract.Address))
	require.NoError(t, err)
	require.Equal(t, delegation.Id, scopedDelegation.Id)

	var delegationIDs []uint64
	err = governanceKeeperInstance.WalkDelegations(ctx, func(d governancetypes.Delegation) (bool, error) {
		delegationIDs = append(delegationIDs, d.Id)
		return false, nil
	})
	require.NoError(t, err)
	require.Equal(t, []uint64{delegation.Id}, delegationIDs)
}
