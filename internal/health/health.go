package health

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Status represents the health state of a dependency.
type Status string

const (
	StatusOK      Status = "ok"
	StatusDegraded Status = "degraded"
)

// Check holds the result of a single health check.
type Check struct {
	Name    string        `json:"name"`
	Status  Status        `json:"status"`
	Latency time.Duration `json:"latency_ms"`
	Error   string        `json:"error,omitempty"`
}

// Report aggregates all checks.
type Report struct {
	Status string  `json:"status"`
	Checks []Check `json:"checks"`
}

// Checker performs a named health check.
type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

// Runner runs a set of Checkers and produces a Report.
type Runner struct {
	checkers []Checker
	timeout  time.Duration
}

// NewRunner creates a Runner with the given checkers and per-check timeout.
func NewRunner(timeout time.Duration, checkers ...Checker) *Runner {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Runner{checkers: checkers, timeout: timeout}
}

// Run executes all checks and returns a Report.
func (r *Runner) Run(ctx context.Context) Report {
	overall := string(StatusOK)
	checks := make([]Check, 0, len(r.checkers))
	for _, c := range r.checkers {
		tctx, cancel := context.WithTimeout(ctx, r.timeout)
		start := time.Now()
		err := c.Check(tctx)
		cancel()
		latency := time.Since(start)
		ch := Check{Name: c.Name(), Status: StatusOK, Latency: latency}
		if err != nil {
			ch.Status = StatusDegraded
			ch.Error = err.Error()
			overall = string(StatusDegraded)
		}
		checks = append(checks, ch)
	}
	return Report{Status: overall, Checks: checks}
}

// VaultChecker pings Vault by hitting its sys/health endpoint.
type VaultChecker struct {
	address string
	client  *http.Client
}

// NewVaultChecker creates a VaultChecker for the given Vault address.
func NewVaultChecker(address string) *VaultChecker {
	return &VaultChecker{address: address, client: &http.Client{Timeout: 3 * time.Second}}
}

func (v *VaultChecker) Name() string { return "vault" }

func (v *VaultChecker) Check(ctx context.Context) error {
	url := fmt.Sprintf("%s/v1/sys/health?standbyok=true", v.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return fmt.Errorf("vault returned status %d", resp.StatusCode)
	}
	return nil
}
