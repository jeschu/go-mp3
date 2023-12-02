package health

import (
	health "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"time"
)

func RegisterHealthChecks() (h health.Health, err error) {
	h = health.New()
	err = h.RegisterCheck(
		checks.NewHostResolveCheck("google.com", 1),
		health.ExecutionTimeout(200*time.Millisecond),
		health.ExecutionPeriod(10*time.Second),
	)
	return
}
