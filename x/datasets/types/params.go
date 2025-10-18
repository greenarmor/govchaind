package types

import (
	"fmt"
	"strings"
)

// NewParams creates a new Params instance.
func NewParams(minApprovals, minScore uint32, metric string) Params {
	return Params{
		MinRequiredApprovals:   minApprovals,
		MinAccountabilityScore: minScore,
		ApprovalMetric:         metric,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(2, 70, "dataset_review")
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.MinRequiredApprovals == 0 {
		return fmt.Errorf("min required approvals must be greater than zero")
	}
	if p.MinAccountabilityScore == 0 || p.MinAccountabilityScore > 100 {
		return fmt.Errorf("min accountability score must be between 1 and 100")
	}
	if strings.TrimSpace(p.ApprovalMetric) == "" {
		return fmt.Errorf("approval metric is required")
	}
	return nil
}
