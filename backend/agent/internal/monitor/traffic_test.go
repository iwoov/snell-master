package monitor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/iwoov/snell-master/backend/agent/internal/manager"
)

type stubRunner struct {
	responses map[string][]byte
	err       error
}

func newStubRunner() *stubRunner {
	return &stubRunner{responses: make(map[string][]byte)}
}

func (s *stubRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	key := name + " " + strings.Join(args, " ")
	if data, ok := s.responses[key]; ok {
		return data, nil
	}
	if s.err != nil {
		return nil, s.err
	}
	return nil, errors.New("command not mocked")
}

func TestTrafficMonitorUpdateTraffic(t *testing.T) {
	runner := newStubRunner()
	jsonTemplate := `{"nftables":[{"rule":{"family":"inet","table":"snell","chain":"traffic","handle":1,"expr":[{"match":{"right":6}},{"match":{"right":%d}},{"counter":{"bytes":%d}}]}}]}`
	runner.responses["nft -j list table inet snell"] = []byte(sprintf(jsonTemplate, 15000, 1000))

	mon := NewTrafficMonitor(runner)

	inst := &manager.Instance{ID: 1, Port: 15000, Username: "user"}
	list := []*manager.Instance{inst}

	// first call primes cache
	stats, err := mon.UpdateTraffic(context.Background(), list)
	if err != nil {
		t.Fatalf("UpdateTraffic() error = %v", err)
	}
	if len(stats) != 0 {
		t.Fatalf("expected no stats on first run")
	}

	runner.responses["nft -j list table inet snell"] = []byte(sprintf(jsonTemplate, 15000, 1600))
	stats, err = mon.UpdateTraffic(context.Background(), list)
	if err != nil {
		t.Fatalf("UpdateTraffic() error = %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 record, got %d", len(stats))
	}
	if stats[0].InstanceID != 1 {
		t.Fatalf("unexpected instance id")
	}
	if stats[0].BytesUpload+stats[0].BytesDownload != 600 {
		t.Fatalf("expected delta 600, got %d", stats[0].BytesUpload+stats[0].BytesDownload)
	}
}

func TestParseRuleExpr(t *testing.T) {
	expr := []interface{}{
		map[string]interface{}{"match": map[string]interface{}{"right": float64(12345)}},
		map[string]interface{}{"counter": map[string]interface{}{"bytes": float64(2048)}},
	}
	port, bytes := parseRuleExpr(expr)
	if port != 12345 || bytes != 2048 {
		t.Fatalf("unexpected parse result port=%d bytes=%d", port, bytes)
	}
}

func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
