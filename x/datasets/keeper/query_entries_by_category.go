package keeper

import (
	"context"

	"govchain/x/datasets/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) EntriesByCategory(ctx context.Context, req *types.QueryEntriesByCategoryRequest) (*types.QueryEntriesByCategoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: Process the query

	return &types.QueryEntriesByCategoryResponse{}, nil
}
