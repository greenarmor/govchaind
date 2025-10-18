package types

import sdkerrors "cosmossdk.io/errors"

var (
	ErrBudgetNotFound           = sdkerrors.Register(ModuleName, 1, "budget not found")
	ErrInvalidBudgetTransition  = sdkerrors.Register(ModuleName, 2, "invalid budget status transition")
	ErrBudgetCreatorInvalidAddr = sdkerrors.Register(ModuleName, 3, "invalid budget creator address")
)
