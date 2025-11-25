package utils

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestIsPortAvailable(t *testing.T) {
	t.Parallel()

	ln := listenTCP(t)
	port := ln.Addr().(*net.TCPAddr).Port

	if IsPortAvailable(port) {
		t.Fatalf("expected port %d to be unavailable", port)
	}

	_ = ln.Close()

	if !IsPortAvailable(port) {
		t.Fatalf("expected port %d to be available after close", port)
	}
}

func TestFindAvailablePort(t *testing.T) {
	t.Parallel()

	ln := listenTCP(t)
	_ = ln.Close()

	port, err := FindAvailablePort(15000, 15100)
	if err != nil {
		if strings.Contains(err.Error(), "no available port") {
			t.Skipf("unable to allocate test port: %v", err)
		}
		t.Fatalf("FindAvailablePort() error = %v", err)
	}

	if port < 15000 || port > 15100 {
		t.Fatalf("port out of expected range: %d", port)
	}
}

func TestPIDFileOperations(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	pidPath := filepath.Join(dir, "agent.pid")

	if err := WritePIDFile(pidPath, os.Getpid()); err != nil {
		t.Fatalf("WritePIDFile() error = %v", err)
	}

	pid, err := ReadPIDFile(pidPath)
	if err != nil {
		t.Fatalf("ReadPIDFile() error = %v", err)
	}

	if pid != os.Getpid() {
		t.Fatalf("expected pid %d, got %d", os.Getpid(), pid)
	}

	if err := RemovePIDFile(pidPath); err != nil {
		t.Fatalf("RemovePIDFile() error = %v", err)
	}

	if _, err := os.Stat(pidPath); !os.IsNotExist(err) {
		t.Fatalf("pid file should be removed")
	}
}

func TestIsProcessRunning(t *testing.T) {
	t.Parallel()

	if !IsProcessRunning(os.Getpid()) {
		t.Fatalf("expected current process to be running")
	}

	if IsProcessRunning(999999) {
		t.Fatalf("expected high pid to be not running")
	}
}

func TestExecuteCommand(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("skip shell command test on windows")
	}

	out, err := ExecuteCommand("/bin/echo", "hello")
	if err != nil {
		t.Fatalf("ExecuteCommand() error = %v", err)
	}

	if strings.TrimSpace(out) != "hello" {
		t.Fatalf("expected output hello, got %s", out)
	}
}

func TestGetPublicIP(t *testing.T) {
	server := newTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("1.2.3.4"))
	}))
	t.Cleanup(server.Close)

	t.Setenv(publicIPServiceEnv, server.URL)

	ip, err := GetPublicIP()
	if err != nil {
		t.Fatalf("GetPublicIP() error = %v", err)
	}

	if ip != "1.2.3.4" {
		t.Fatalf("expected IP 1.2.3.4, got %s", ip)
	}
}

func listenTCP(t *testing.T) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skipf("network operations are not permitted: %v", err)
		}
		t.Fatalf("listen: %v", err)
	}
	return ln
}

func newTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	var (
		server      *httptest.Server
		panicReason interface{}
	)

	func() {
		defer func() {
			panicReason = recover()
		}()
		server = httptest.NewServer(handler)
	}()

	if panicReason != nil {
		msg := fmt.Sprint(panicReason)
		if strings.Contains(msg, "operation not permitted") {
			t.Skipf("network operations are not permitted: %v", panicReason)
		}
		panic(panicReason)
	}

	return server
}
