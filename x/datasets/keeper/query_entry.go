package keeper

import (
	"context"
	"errors"

	"govchain/x/datasets/types"

	"cosmossdk.io/collections"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ListEntry(ctx context.Context, req *types.QueryAllEntryRequest) (*types.QueryAllEntryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	entrys, pageRes, err := query.CollectionPaginate(
		ctx,
		q.k.Entry,
		req.Pagination,
		func(_ uint64, value types.Entry) (types.Entry, error) {
			return value, nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllEntryResponse{Entry: entrys, Pagination: pageRes}, nil
}

func (q queryServer) GetEntry(ctx context.Context, req *types.QueryGetEntryRequest) (*types.QueryGetEntryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	entry, err := q.k.Entry.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, sdkerrors.ErrKeyNotFound
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &types.QueryGetEntryResponse{Entry: entry}, nil
}
