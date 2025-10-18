package keeper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"govchain/x/datasets/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto/tmhash"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateEntry(ctx context.Context, msg *types.MsgCreateEntry) (*types.MsgCreateEntryResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to load params")
	}

	if strings.TrimSpace(msg.LiabilityContract) == "" {
		return nil, errorsmod.Wrap(types.ErrMissingLiabilityContract, "liability contract must be provided")
	}

	requiredApprovals := msg.RequiredApprovals
	if requiredApprovals == 0 {
		requiredApprovals = params.MinRequiredApprovals
	}
	if requiredApprovals < params.MinRequiredApprovals {
		return nil, errorsmod.Wrapf(types.ErrInsufficientApprovals, "required approvals %d below minimum %d", requiredApprovals, params.MinRequiredApprovals)
	}

	// Get SDK context to access transaction information
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	nextId, err := k.EntrySeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get next id")
	}

	// Extract transaction hash
	txBytes := sdkCtx.TxBytes()
	txHash := fmt.Sprintf("%X", tmhash.Sum(txBytes))

	entry := types.Entry{
		Id:                nextId,
		Creator:           msg.Creator,
		Title:             msg.Title,
		Description:       msg.Description,
		IpfsCid:           msg.IpfsCid,
		MimeType:          msg.MimeType,
		FileName:          msg.FileName,
		FileUrl:           msg.FileUrl,
		FallbackUrl:       msg.FallbackUrl,
		FileSize:          msg.FileSize,
		ChecksumSha_256:   msg.ChecksumSha_256,
		Agency:            msg.Agency,
		Category:          msg.Category,
		Submitter:         msg.Submitter,
		Timestamp:         msg.Timestamp,
		PinCount:          msg.PinCount,
		TxHash:            txHash,
		Status:            types.EntryStatus_ENTRY_STATUS_PENDING,
		RequiredApprovals: requiredApprovals,
		Approvals:         []*types.Approval{},
		LiabilityContract: msg.LiabilityContract,
	}

	if err := entry.Validate(params); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if err = k.Entry.Set(
		ctx,
		nextId,
		entry,
	); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to set entry")
	}

	return &types.MsgCreateEntryResponse{
		Id: nextId,
	}, nil
}

func (k msgServer) UpdateEntry(ctx context.Context, msg *types.MsgUpdateEntry) (*types.MsgUpdateEntryResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Checks that the element exists
	val, err := k.Entry.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get entry")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to load params")
	}

	val.Title = msg.Title
	val.Description = msg.Description
	val.IpfsCid = msg.IpfsCid
	val.MimeType = msg.MimeType
	val.FileName = msg.FileName
	val.FileUrl = msg.FileUrl
	val.FallbackUrl = msg.FallbackUrl
	val.FileSize = msg.FileSize
	val.ChecksumSha_256 = msg.ChecksumSha_256
	val.Agency = msg.Agency
	val.Category = msg.Category
	val.Submitter = msg.Submitter
	val.Timestamp = msg.Timestamp
	val.PinCount = msg.PinCount

	if err := val.Validate(params); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if err := k.Entry.Set(ctx, msg.Id, val); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update entry")
	}

	return &types.MsgUpdateEntryResponse{}, nil
}

func (k msgServer) DeleteEntry(ctx context.Context, msg *types.MsgDeleteEntry) (*types.MsgDeleteEntryResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Checks that the element exists
	val, err := k.Entry.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get entry")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if val.Status == types.EntryStatus_ENTRY_STATUS_APPROVED {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "cannot delete an approved entry")
	}

	if err := k.Entry.Remove(ctx, msg.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete entry")
	}

	return &types.MsgDeleteEntryResponse{}, nil
}

func (k msgServer) ApproveEntry(ctx context.Context, msg *types.MsgApproveEntry) (*types.MsgApproveEntryResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Approver); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}
	if strings.TrimSpace(msg.LiabilitySignature) == "" {
		return nil, errorsmod.Wrap(types.ErrMissingLiabilitySignature, "liability signature is required")
	}
	if strings.TrimSpace(msg.EvidenceUri) == "" {
		return nil, errorsmod.Wrap(types.ErrMissingEvidenceURI, "evidence URI is required")
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to load params")
	}

	entry, err := k.Entry.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("entry %d not found", msg.Id))
		}
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to load entry")
	}

	if entry.Status == types.EntryStatus_ENTRY_STATUS_APPROVED {
		return nil, errorsmod.Wrap(types.ErrAlreadyApproved, "entry already approved")
	}

	for _, approval := range entry.Approvals {
		if approval != nil && approval.Approver == msg.Approver {
			return nil, errorsmod.Wrap(types.ErrDuplicateApproval, "approver has already endorsed this entry")
		}
	}

	scorecard, err := k.accountabilityKeeper.GetScorecardBySubjectMetric(ctx, msg.Approver, params.ApprovalMetric)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrAccountabilityLookup, err.Error())
	}
	if scorecard.Score < params.MinAccountabilityScore {
		return nil, errorsmod.Wrapf(types.ErrAccountabilityThreshold, "score %d below required %d", scorecard.Score, params.MinAccountabilityScore)
	}

	approval := &types.Approval{
		Approver:            msg.Approver,
		AccountabilityScore: scorecard.Score,
		EvidenceUri:         msg.EvidenceUri,
		LiabilitySignature:  msg.LiabilitySignature,
		Timestamp:           time.Now().UTC().Format(time.RFC3339),
	}

	entry.Approvals = append(entry.Approvals, approval)
	if uint32(len(entry.Approvals)) >= entry.RequiredApprovals {
		entry.Status = types.EntryStatus_ENTRY_STATUS_APPROVED
	}

	if err := entry.Validate(params); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if err := k.Entry.Set(ctx, entry.Id, entry); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to persist approval")
	}

	return &types.MsgApproveEntryResponse{}, nil
}
