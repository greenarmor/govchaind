package keeper

import (
	"context"
	"errors"
	"fmt"

	"govchain/x/datasets/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateEntry(ctx context.Context, msg *types.MsgCreateEntry) (*types.MsgCreateEntryResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	nextId, err := k.EntrySeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get next id")
	}

	var entry = types.Entry{
		Id:              nextId,
		Creator:         msg.Creator,
		Title:           msg.Title,
		Description:     msg.Description,
		IpfsCid:         msg.IpfsCid,
		MimeType:        msg.MimeType,
		FileName:        msg.FileName,
		FileUrl:         msg.FileUrl,
		FallbackUrl:     msg.FallbackUrl,
		FileSize:        msg.FileSize,
		ChecksumSha_256: msg.ChecksumSha_256,
		Agency:          msg.Agency,
		Category:        msg.Category,
		Submitter:       msg.Submitter,
		Timestamp:       msg.Timestamp,
		PinCount:        msg.PinCount,
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

	var entry = types.Entry{
		Creator:         msg.Creator,
		Id:              msg.Id,
		Title:           msg.Title,
		Description:     msg.Description,
		IpfsCid:         msg.IpfsCid,
		MimeType:        msg.MimeType,
		FileName:        msg.FileName,
		FileUrl:         msg.FileUrl,
		FallbackUrl:     msg.FallbackUrl,
		FileSize:        msg.FileSize,
		ChecksumSha_256: msg.ChecksumSha_256,
		Agency:          msg.Agency,
		Category:        msg.Category,
		Submitter:       msg.Submitter,
		Timestamp:       msg.Timestamp,
		PinCount:        msg.PinCount,
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

	if err := k.Entry.Set(ctx, msg.Id, entry); err != nil {
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

	if err := k.Entry.Remove(ctx, msg.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete entry")
	}

	return &types.MsgDeleteEntryResponse{}, nil
}
