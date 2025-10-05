package keeper_test

import (
	"context"
	"strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"govchain/x/datasets/keeper"
	"govchain/x/datasets/types"
)

func createNEntry(keeper keeper.Keeper, ctx context.Context, n int) []types.Entry {
	items := make([]types.Entry, n)
	for i := range items {
		iu := uint64(i)
		items[i].Id = iu
		items[i].Title = strconv.Itoa(i)
		items[i].Description = strconv.Itoa(i)
		items[i].IpfsCid = strconv.Itoa(i)
		items[i].MimeType = strconv.Itoa(i)
		items[i].FileName = strconv.Itoa(i)
		items[i].FileUrl = strconv.Itoa(i)
		items[i].FallbackUrl = strconv.Itoa(i)
		items[i].FileSize = strconv.Itoa(i)
		items[i].ChecksumSha_256 = strconv.Itoa(i)
		items[i].Agency = strconv.Itoa(i)
		items[i].Category = strconv.Itoa(i)
		items[i].Submitter = strconv.Itoa(i)
		items[i].Timestamp = strconv.Itoa(i)
		items[i].PinCount = strconv.Itoa(i)
		_ = keeper.Entry.Set(ctx, iu, items[i])
		_ = keeper.EntrySeq.Set(ctx, iu)
	}
	return items
}

func TestEntryQuerySingle(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNEntry(f.keeper, f.ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetEntryRequest
		response *types.QueryGetEntryResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetEntryRequest{Id: msgs[0].Id},
			response: &types.QueryGetEntryResponse{Entry: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetEntryRequest{Id: msgs[1].Id},
			response: &types.QueryGetEntryResponse{Entry: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetEntryRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.GetEntry(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.EqualExportedValues(t, tc.response, response)
			}
		})
	}
}

func TestEntryQueryPaginated(t *testing.T) {
	f := initFixture(t)
	qs := keeper.NewQueryServerImpl(f.keeper)
	msgs := createNEntry(f.keeper, f.ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllEntryRequest {
		return &types.QueryAllEntryRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListEntry(f.ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Entry), step)
			require.Subset(t, msgs, resp.Entry)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := qs.ListEntry(f.ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Entry), step)
			require.Subset(t, msgs, resp.Entry)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.ListEntry(f.ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.EqualExportedValues(t, msgs, resp.Entry)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.ListEntry(f.ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
