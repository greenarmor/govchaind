package types

import (
	"context"

	budgettypes "govchain/x/budgetregistry/types"
)

// BudgetKeeper defines the subset of the budget registry keeper used by the procurement module.
type BudgetKeeper interface {
	GetBudget(context.Context, uint64) (budgettypes.Budget, error)
}
