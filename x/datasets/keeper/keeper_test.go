package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	accountabilitykeeper "govchain/x/accountabilityscores/keeper"
	accountabilitytypes "govchain/x/accountabilityscores/types"
	"govchain/x/datasets/keeper"
	module "govchain/x/datasets/module"
	"govchain/x/datasets/types"
)

type fixture struct {
	ctx                  context.Context
	keeper               keeper.Keeper
	addressCodec         address.Codec
	accountabilityKeeper accountabilitykeeper.Keeper
}

func (f *fixture) setScorecard(t *testing.T, addr, metric string, score uint32) {
	t.Helper()
	_, err := f.accountabilityKeeper.UpsertScorecard(f.ctx, accountabilitytypes.Scorecard{
		Subject:     addr,
		Metric:      metric,
		Score:       score,
		Weight:      1,
		EvidenceURI: "ipfs://evidence",
		UpdatedBy:   addr,
	})
	if err != nil {
		t.Fatalf("failed to set scorecard: %v", err)
	}
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	accountabilityKey := storetypes.NewKVStoreKey(accountabilitytypes.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	accountabilityStoreService := runtime.NewKVStoreService(accountabilityKey)
	testCtx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test"))
	// mount accountability store on the same multi-store
	testCtx.CMS.MountStoreWithDB(accountabilityKey, storetypes.StoreTypeIAVL, testCtx.DB)
	if err := testCtx.CMS.LoadLatestVersion(); err != nil {
		t.Fatalf("failed to load accountability store: %v", err)
	}
	ctx := testCtx.Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)

	accountabilityKeeper := accountabilitykeeper.NewKeeper(
		accountabilityStoreService,
		encCfg.Codec,
		addressCodec,
	)

	k := keeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		authority,
		accountabilityKeeper,
	)

	// Initialize params
	if err := k.Params.Set(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:                  ctx,
		keeper:               k,
		addressCodec:         addressCodec,
		accountabilityKeeper: accountabilityKeeper,
	}
}
