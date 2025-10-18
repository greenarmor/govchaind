package keeper_test

import (
	"context"
	"testing"

	address "cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"govchain/x/budgetregistry/keeper"
	"govchain/x/budgetregistry/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig()
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_budget")).Ctx

	k := keeper.NewKeeper(storeService, encCfg.Codec, addressCodec)

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
	}
}

func (f *fixture) newAddress(t *testing.T, raw string) string {
	t.Helper()
	addr, err := f.addressCodec.BytesToString([]byte(raw))
	if err != nil {
		t.Fatalf("failed to convert address: %v", err)
	}
	return addr
}

func TestRegisterBudget(t *testing.T) {
	f := initFixture(t)

	creator := f.newAddress(t, "budgetCreatorAddr__________")
	amount := sdk.NewInt64Coin("php", 1_000_000).Amount

	registered, err := f.keeper.RegisterBudget(f.ctx, types.Budget{
		Agency:      "DOF",
		FiscalYear:  "2025",
		Title:       "General Appropriations",
		Amount:      amount,
		Currency:    "php",
		MetadataURI: "ipfs://cid",
		CreatedBy:   creator,
	})
	if err != nil {
		t.Fatalf("register budget: %v", err)
	}
	if registered.Id == 0 {
		t.Fatalf("expected budget id to be set")
	}
	if registered.Status != types.BudgetStatusDraft {
		t.Fatalf("unexpected status %s", registered.Status)
	}

	stored, err := f.keeper.GetBudget(f.ctx, registered.Id)
	if err != nil {
		t.Fatalf("get budget: %v", err)
	}
	if stored.Title != "General Appropriations" {
		t.Fatalf("unexpected title %s", stored.Title)
	}
}

func TestUpdateBudgetStatus(t *testing.T) {
	f := initFixture(t)
	creator := f.newAddress(t, "budgetCreatorAddr__________")
	amount := sdk.NewInt64Coin("php", 5_000_000).Amount

	registered, err := f.keeper.RegisterBudget(f.ctx, types.Budget{
		Agency:      "DBM",
		FiscalYear:  "2024",
		Title:       "Budget Execution",
		Amount:      amount,
		Currency:    "php",
		MetadataURI: "ipfs://cid",
		CreatedBy:   creator,
	})
	if err != nil {
		t.Fatalf("register budget: %v", err)
	}

	updated, err := f.keeper.UpdateBudgetStatus(f.ctx, registered.Id, types.BudgetStatusLegislated)
	if err != nil {
		t.Fatalf("update status: %v", err)
	}
	if updated.Status != types.BudgetStatusLegislated {
		t.Fatalf("unexpected status %s", updated.Status)
	}

	if _, err := f.keeper.UpdateBudgetStatus(f.ctx, registered.Id, types.BudgetStatusDraft); err == nil {
		t.Fatalf("expected invalid transition error")
	}
}
