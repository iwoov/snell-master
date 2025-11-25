package manager

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/iwoov/snell-master/backend/agent/internal/client"
	agentutils "github.com/iwoov/snell-master/backend/pkg/utils"
)

func TestGenerateConfig(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	mgr := NewInstanceManager(dir, "/bin/true", 10000, 20000)
	inst := &Instance{ID: 1, Port: 12345, PSK: "secret"}

	path, err := mgr.generateConfig(inst)
	if err != nil {
		t.Fatalf("generateConfig() error = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if string(data) == "" {
		t.Fatalf("config should not be empty")
	}
}

func TestStartStopInstance(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip process test on windows")
	}

	dir := t.TempDir()
	binary := writeFakeSnellBinary(t, dir)
	port := pickTestPort(t, 15000, 16000)

	mgr := NewInstanceManager(dir, binary, 15000, 20000)
	inst := &Instance{ID: 1, Port: port, PSK: "secret"}

	if err := mgr.StartInstance(inst); err != nil {
		t.Fatalf("StartInstance() error = %v", err)
	}
	if inst.PID == 0 {
		t.Fatalf("expected PID to be set")
	}

	if err := mgr.StopInstance(inst); err != nil {
		t.Fatalf("StopInstance() error = %v", err)
	}
	if inst.PID != 0 {
		t.Fatalf("expected PID reset to 0")
	}
}

func TestSyncInstances(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip process test on windows")
	}

	dir := t.TempDir()
	binary := writeFakeSnellBinary(t, dir)
	mgr := NewInstanceManager(dir, binary, 10000, 20000)

	port := pickTestPort(t, 20000, 21000)
	remote := []client.InstanceConfig{{ID: 1, UserID: 1, Username: "user", Port: port, PSK: "a", Version: 1}}
	if err := mgr.SyncInstances(remote); err != nil {
		t.Fatalf("SyncInstances() create error = %v", err)
	}
	if len(mgr.GetAllInstances()) != 1 {
		t.Fatalf("expected 1 instance after sync")
	}

	remote[0].PSK = "b"
	if err := mgr.SyncInstances(remote); err != nil {
		t.Fatalf("SyncInstances() update error = %v", err)
	}

	if err := mgr.SyncInstances(nil); err != nil {
		t.Fatalf("SyncInstances() delete error = %v", err)
	}
	if len(mgr.GetAllInstances()) != 0 {
		t.Fatalf("expected 0 instance after delete")
	}
}

func TestPortManagement(t *testing.T) {
	mgr := NewInstanceManager(t.TempDir(), "/bin/true", 10000, 20000)
	mgr.setInstance(&Instance{ID: 1, Port: 10001})
	mgr.setInstance(&Instance{ID: 2, Port: 10002})

	ports := mgr.GetUsedPorts()
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(ports))
	}
	if !mgr.IsPortUsed(10002) {
		t.Fatalf("expected port reported as used")
	}
	if mgr.IsPortUsed(9999) {
		t.Fatalf("port 9999 should not be used")
	}
}

func writeFakeSnellBinary(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "snell.sh")
	script := "#!/bin/sh\ntrap 'exit 0' TERM\nwhile true; do sleep 1; done\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}
	return path
}

func pickTestPort(t *testing.T, start, end int) int {
	t.Helper()
	port, err := agentutils.FindAvailablePort(start, end)
	if err != nil {
		if strings.Contains(err.Error(), "no available port") || strings.Contains(err.Error(), "operation not permitted") {
			t.Skipf("skipping port-dependent test: %v", err)
		}
		t.Fatalf("FindAvailablePort() error = %v", err)
	}
	return port
}
