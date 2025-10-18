package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/datasets module sentinel errors
var (
	ErrInvalidSigner             = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrMissingLiabilityContract  = errors.Register(ModuleName, 1101, "liability contract required")
	ErrInsufficientApprovals     = errors.Register(ModuleName, 1102, "required approvals below module minimum")
	ErrMissingLiabilitySignature = errors.Register(ModuleName, 1103, "liability signature required")
	ErrMissingEvidenceURI        = errors.Register(ModuleName, 1104, "evidence URI required")
	ErrAlreadyApproved           = errors.Register(ModuleName, 1105, "entry already approved")
	ErrDuplicateApproval         = errors.Register(ModuleName, 1106, "duplicate approval")
	ErrAccountabilityLookup      = errors.Register(ModuleName, 1107, "failed to resolve accountability score")
	ErrAccountabilityThreshold   = errors.Register(ModuleName, 1108, "accountability score below threshold")
)
