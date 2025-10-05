package keeper

import (
	"context"

	"govchain/x/datasets/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) EntriesByMimetype(ctx context.Context, req *types.QueryEntriesByMimetypeRequest) (*types.QueryEntriesByMimetypeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// TODO: Process the query

	return &types.QueryEntriesByMimetypeResponse{}, nil
}
