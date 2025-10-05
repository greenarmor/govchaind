package keeper

import (
	"context"

	"govchain/x/datasets/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) EntriesByAgency(ctx context.Context, req *types.QueryEntriesByAgencyRequest) (*types.QueryEntriesByAgencyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: Process the query

	return &types.QueryEntriesByAgencyResponse{}, nil
}
