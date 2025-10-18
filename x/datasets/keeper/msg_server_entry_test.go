package keeper_test

import (
	"fmt"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"govchain/x/datasets/keeper"
	"govchain/x/datasets/types"
)

func TestEntryMsgServerCreate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
                resp, err := srv.CreateEntry(f.ctx, &types.MsgCreateEntry{
                        Creator:           creator,
                        Title:             fmt.Sprintf("dataset-%d", i),
                        LiabilityContract: "ipfs://contract",
                        RequiredApprovals: 2,
                })
                require.NoError(t, err)
                require.Equal(t, i, int(resp.Id))
        }
}

func TestEntryMsgServerUpdate(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	_, err = srv.CreateEntry(f.ctx, &types.MsgCreateEntry{
		Creator:           creator,
		LiabilityContract: "ipfs://contract",
		RequiredApprovals: 2,
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgUpdateEntry
		err     error
	}{
		{
			desc:    "invalid address",
			request: &types.MsgUpdateEntry{Creator: "invalid"},
			err:     sdkerrors.ErrInvalidAddress,
		},
		{
			desc:    "unauthorized",
			request: &types.MsgUpdateEntry{Creator: unauthorizedAddr},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "key not found",
			request: &types.MsgUpdateEntry{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc:    "completed",
			request: &types.MsgUpdateEntry{Creator: creator},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.UpdateEntry(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEntryMsgServerDelete(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	unauthorizedAddr, err := f.addressCodec.BytesToString([]byte("unauthorizedAddr___________"))
	require.NoError(t, err)

	_, err = srv.CreateEntry(f.ctx, &types.MsgCreateEntry{
		Creator:           creator,
		LiabilityContract: "ipfs://contract",
		RequiredApprovals: 2,
	})
	require.NoError(t, err)

	tests := []struct {
		desc    string
		request *types.MsgDeleteEntry
		err     error
	}{
		{
			desc:    "invalid address",
			request: &types.MsgDeleteEntry{Creator: "invalid"},
			err:     sdkerrors.ErrInvalidAddress,
		},
		{
			desc:    "unauthorized",
			request: &types.MsgDeleteEntry{Creator: unauthorizedAddr},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "key not found",
			request: &types.MsgDeleteEntry{Creator: creator, Id: 10},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc:    "completed",
			request: &types.MsgDeleteEntry{Creator: creator},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err = srv.DeleteEntry(f.ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEntryMsgServerApprove(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	creator, err := f.addressCodec.BytesToString([]byte("signerAddr__________________"))
	require.NoError(t, err)

	approverA, err := f.addressCodec.BytesToString([]byte("approverAddrA_______________"))
	require.NoError(t, err)
	approverB, err := f.addressCodec.BytesToString([]byte("approverAddrB_______________"))
	require.NoError(t, err)

	params, err := f.keeper.Params.Get(f.ctx)
	require.NoError(t, err)
	f.setScorecard(t, approverA, params.ApprovalMetric, params.MinAccountabilityScore)
	f.setScorecard(t, approverB, params.ApprovalMetric, params.MinAccountabilityScore+5)

	createResp, err := srv.CreateEntry(f.ctx, &types.MsgCreateEntry{
		Creator:           creator,
		Title:             "critical-dataset",
		LiabilityContract: "ipfs://contract",
		RequiredApprovals: params.MinRequiredApprovals,
	})
	require.NoError(t, err)

	_, err = srv.ApproveEntry(f.ctx, &types.MsgApproveEntry{
		Approver:           approverA,
		Id:                 createResp.Id,
		EvidenceUri:        "ipfs://evidence-a",
		LiabilitySignature: "sig-a",
	})
	require.NoError(t, err)

	entry, err := f.keeper.Entry.Get(f.ctx, createResp.Id)
	require.NoError(t, err)
	require.Equal(t, types.EntryStatus_ENTRY_STATUS_PENDING, entry.Status)
	require.Len(t, entry.Approvals, 1)

	_, err = srv.ApproveEntry(f.ctx, &types.MsgApproveEntry{
		Approver:           approverA,
		Id:                 createResp.Id,
		EvidenceUri:        "ipfs://evidence-a",
		LiabilitySignature: "sig-a",
	})
	require.ErrorIs(t, err, types.ErrDuplicateApproval)

	_, err = srv.ApproveEntry(f.ctx, &types.MsgApproveEntry{
		Approver:           approverB,
		Id:                 createResp.Id,
		EvidenceUri:        "ipfs://evidence-b",
		LiabilitySignature: "sig-b",
	})
	require.NoError(t, err)

	entry, err = f.keeper.Entry.Get(f.ctx, createResp.Id)
	require.NoError(t, err)
	require.Equal(t, types.EntryStatus_ENTRY_STATUS_APPROVED, entry.Status)
	require.Len(t, entry.Approvals, 2)
}
