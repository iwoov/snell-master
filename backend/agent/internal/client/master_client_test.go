package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMasterClientGet(t *testing.T) {
	t.Parallel()

	server := newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/ping" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Fatalf("unexpected auth header: %s", got)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	t.Cleanup(server.Close)

	client := NewMasterClient(server.URL, "token")
	data, err := client.Get("/api/ping")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if string(data) != `{"ok":true}` {
		t.Fatalf("unexpected data: %s", data)
	}
}

func TestMasterClientPostHTTPError(t *testing.T) {
	t.Parallel()

	server := newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"bad"}`))
	}))
	t.Cleanup(server.Close)

	client := NewMasterClient(server.URL, "token")
	_, err := client.Post("/api/test", map[string]string{"foo": "bar"})
	if err == nil {
		t.Fatal("expected error")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected HTTPError, got %T", err)
	}
	if httpErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", httpErr.StatusCode)
	}
}

func TestFetchConfig(t *testing.T) {
	t.Parallel()

	server := newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/agent/config" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code":0,"message":"ok","data":{"instances":[{"id":1,"user_id":10,"username":"user","port":1234,"psk":"psk","version":1}]}}`))
	}))
	t.Cleanup(server.Close)

	client := NewMasterClient(server.URL, "token")
	configs, err := client.FetchConfig()
	if err != nil {
		t.Fatalf("FetchConfig() error = %v", err)
	}
	if len(configs) != 1 || configs[0].ID != 1 {
		t.Fatalf("unexpected configs: %#v", configs)
	}
}

func TestFetchConfigErrorCode(t *testing.T) {
	t.Parallel()

	server := newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":100,"message":"failed","data":{"instances":[]}}`))
	}))
	t.Cleanup(server.Close)

	client := NewMasterClient(server.URL, "token")
	if _, err := client.FetchConfig(); err == nil {
		t.Fatal("expected error")
	}
}

func TestReportHeartbeat(t *testing.T) {
	t.Parallel()

	var body HeartbeatRequest
	server := newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_, _ = w.Write([]byte(`{"code":0,"message":"ok"}`))
	}))
	t.Cleanup(server.Close)

	client := NewMasterClient(server.URL, "token")
	if err := client.ReportHeartbeat(10, 20, 3, "v1.0"); err != nil {
		t.Fatalf("ReportHeartbeat() error = %v", err)
	}
	if body.CPUUsage != 10 || body.MemoryUsage != 20 || body.InstanceCount != 3 || body.Version != "v1.0" {
		t.Fatalf("unexpected body: %#v", body)
	}
}

func TestReportTraffic(t *testing.T) {
	t.Parallel()

	var body TrafficReportRequest
	server := newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_, _ = w.Write([]byte(`{"code":0,"message":"ok"}`))
	}))
	t.Cleanup(server.Close)

	client := NewMasterClient(server.URL, "token")
	traffic := []InstanceTraffic{{InstanceID: 1, BytesUpload: 100, BytesDownload: 200}}
	if err := client.ReportTraffic(traffic); err != nil {
		t.Fatalf("ReportTraffic() error = %v", err)
	}
	if len(body.Traffic) != 1 || body.Traffic[0].InstanceID != 1 {
		t.Fatalf("unexpected traffic body: %#v", body)
	}
}

func TestReportStatus(t *testing.T) {
	t.Parallel()

	var body StatusReportRequest
	server := newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(data, &body)
		_, _ = w.Write([]byte(`{"code":0,"message":"ok"}`))
	}))
	t.Cleanup(server.Close)

	client := NewMasterClient(server.URL, "token")
	statuses := []InstanceStatus{{InstanceID: 2, Status: 1}}
	if err := client.ReportStatus(statuses); err != nil {
		t.Fatalf("ReportStatus() error = %v", err)
	}
	if len(body.Statuses) != 1 || body.Statuses[0].InstanceID != 2 {
		t.Fatalf("unexpected status body: %#v", body)
	}
}

func TestGetSnellConfig(t *testing.T) {
	t.Parallel()

	server := newTestHTTPServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/agent/snell-config" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code":0,"message":"ok","data":{"version":"5.0.1","base_url":"http://example","download_urls":{"amd64":"http://example/amd64.zip"}}}`))
	}))
	t.Cleanup(server.Close)

	client := NewMasterClient(server.URL, "token")
	cfg, err := client.GetSnellConfig()
	if err != nil {
		t.Fatalf("GetSnellConfig() error = %v", err)
	}
	if cfg.Version != "5.0.1" || cfg.DownloadURLs["amd64"] == "" {
		t.Fatalf("unexpected snell config: %#v", cfg)
	}
}

func newTestHTTPServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()

	var (
		server      *httptest.Server
		panicReason any
	)

	func() {
		defer func() { panicReason = recover() }()
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
