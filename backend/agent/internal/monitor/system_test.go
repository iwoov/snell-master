package monitor

import (
	"context"
	"errors"
	"testing"
)

func TestSystemMonitorUpdate(t *testing.T) {
	mon := NewSystemMonitor(
		WithCPUProvider(func(context.Context) (float64, error) { return 12.5, nil }),
		WithMemoryProvider(func(context.Context) (float64, error) { return 34.9, nil }),
		WithDiskProvider(func(context.Context) (float64, error) { return 56.1, nil }),
	)

	if err := mon.Update(context.Background()); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if mon.CPUUsage() != 13 {
		t.Fatalf("expected cpu usage 13, got %d", mon.CPUUsage())
	}
	if mon.MemoryUsage() != 35 {
		t.Fatalf("expected memory usage 35, got %d", mon.MemoryUsage())
	}
	if mon.DiskUsage() != 56 {
		t.Fatalf("expected disk usage 56, got %d", mon.DiskUsage())
	}
	if mon.LastUpdated().IsZero() {
		t.Fatalf("last updated should be set")
	}
}

func TestSystemMonitorUpdateError(t *testing.T) {
	wantErr := errors.New("boom")
	mon := NewSystemMonitor(WithCPUProvider(func(context.Context) (float64, error) {
		return 0, wantErr
	}))

	if err := mon.Update(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}
