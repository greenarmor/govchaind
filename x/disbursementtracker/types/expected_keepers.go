package types

import (
	"context"

	procurementtypes "govchain/x/procurementledger/types"
)

// ProcurementKeeper defines the procurement interface consumed by the disbursement module.
type ProcurementKeeper interface {
	GetProcurement(context.Context, uint64) (procurementtypes.Procurement, error)
}
