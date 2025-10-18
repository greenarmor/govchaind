package types

import (
	"fmt"
	"strings"
)

// Validate enforces invariants on the entry against module parameters.
func (e Entry) Validate(params Params) error {
	if strings.TrimSpace(e.LiabilityContract) == "" {
		return fmt.Errorf("liability contract is required")
	}

	switch e.Status {
	case EntryStatus_ENTRY_STATUS_PENDING, EntryStatus_ENTRY_STATUS_APPROVED:
	default:
		return fmt.Errorf("invalid entry status: %s", e.Status.String())
	}

	if e.RequiredApprovals < params.MinRequiredApprovals {
		return fmt.Errorf("required approvals %d below minimum %d", e.RequiredApprovals, params.MinRequiredApprovals)
	}

	if uint32(len(e.Approvals)) > e.RequiredApprovals {
		return fmt.Errorf("stored approvals exceed required approvals")
	}

	seenApprovers := make(map[string]struct{})
	for _, approval := range e.Approvals {
		if strings.TrimSpace(approval.Approver) == "" {
			return fmt.Errorf("approval approver is required")
		}
		if _, exists := seenApprovers[approval.Approver]; exists {
			return fmt.Errorf("duplicate approval for %s", approval.Approver)
		}
		seenApprovers[approval.Approver] = struct{}{}

		if approval.AccountabilityScore < params.MinAccountabilityScore {
			return fmt.Errorf("approval from %s does not meet accountability threshold", approval.Approver)
		}
		if strings.TrimSpace(approval.LiabilitySignature) == "" {
			return fmt.Errorf("approval from %s is missing liability signature", approval.Approver)
		}
		if strings.TrimSpace(approval.EvidenceUri) == "" {
			return fmt.Errorf("approval from %s is missing evidence URI", approval.Approver)
		}
	}

	if e.Status == EntryStatus_ENTRY_STATUS_APPROVED && uint32(len(e.Approvals)) < e.RequiredApprovals {
		return fmt.Errorf("approved entry is missing required approvals")
	}

	if e.Status == EntryStatus_ENTRY_STATUS_PENDING && uint32(len(e.Approvals)) >= e.RequiredApprovals {
		return fmt.Errorf("pending entry already satisfies approval threshold")
	}

	return nil
}
