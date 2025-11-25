package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/iwoov/snell-master/backend/agent/internal/client"
	"github.com/iwoov/snell-master/backend/agent/internal/manager"
	"github.com/iwoov/snell-master/backend/pkg/logger"
)

const (
	nftBinary         = "nft"
	defaultNFTTimeout = 3 * time.Second
	defaultNFTFamily  = "inet"
	defaultNFTTable   = "snell"
	defaultNFTChain   = "traffic"
)

// CommandRunner 用于在测试中注入自定义命令执行器。
type CommandRunner interface {
	Run(ctx context.Context, name string, args ...string) ([]byte, error)
}

// TrafficMonitor 负责读取 nftables 统计信息.
type TrafficMonitor struct {
	runner    CommandRunner
	family    string
	table     string
	chain     string
	lastStats map[uint]*PortTraffic
}

// PortTraffic 记录实例上次统计的累计流量。
type PortTraffic struct {
	Port       int
	BytesTotal int64
	Timestamp  time.Time
}

// NFTablesOutput 对应 nft -j 输出。
type NFTablesOutput struct {
	Nftables []NFTableObject `json:"nftables"`
}

type NFTableObject struct {
	Rule *NFTableRule `json:"rule,omitempty"`
}

// NFTableRule 描述一条 nft 规则。
type NFTableRule struct {
	Family  string        `json:"family"`
	Table   string        `json:"table"`
	Chain   string        `json:"chain"`
	Handle  int64         `json:"handle,omitempty"`
	Comment string        `json:"comment,omitempty"`
	Expr    []interface{} `json:"expr"`
}

// NewTrafficMonitor 创建 TrafficMonitor。
func NewTrafficMonitor(runner CommandRunner) *TrafficMonitor {
	if runner == nil {
		runner = &execRunner{}
	}
	tm := &TrafficMonitor{
		runner:    runner,
		family:    defaultNFTFamily,
		table:     defaultNFTTable,
		chain:     defaultNFTChain,
		lastStats: make(map[uint]*PortTraffic),
	}
	if err := tm.EnsureChain(context.Background()); err != nil {
		logger.WithModule("monitor").Warnf("init nftables chain failed: %v", err)
	}
	return tm
}

// EnsureChain 确保存放统计的 nftables 表与链存在。
func (m *TrafficMonitor) EnsureChain(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), defaultNFTTimeout)
		defer cancel()
	}
	if _, err := m.runner.Run(ctx, nftBinary, "list", "table", m.family, m.table); err == nil {
		return nil
	}
	if _, err := m.runner.Run(ctx, nftBinary, "add", "table", m.family, m.table); err != nil {
		return fmt.Errorf("create nft table: %w", err)
	}
	chainSpec := "{ type filter hook input priority 0; }"
	if _, err := m.runner.Run(ctx, nftBinary, "add", "chain", m.family, m.table, m.chain, chainSpec); err != nil {
		return fmt.Errorf("create nft chain: %w", err)
	}
	return nil
}

// AddInstanceRules 为实例端口添加统计规则。
func (m *TrafficMonitor) AddInstanceRules(ctx context.Context, inst *manager.Instance) error {
	if inst == nil {
		return fmt.Errorf("instance is nil")
	}
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), defaultNFTTimeout)
		defer cancel()
	}
	port := strconv.Itoa(inst.Port)
	comment := inst.Username

	for _, proto := range []string{"tcp", "udp"} {
		if _, err := m.runner.Run(ctx, nftBinary, "add", "rule", m.family, m.table, m.chain,
			proto, "dport", port, "counter", "comment", comment); err != nil {
			return fmt.Errorf("add %s rule: %w", proto, err)
		}
	}
	return nil
}

// RemoveInstanceRules 删除实例关联的规则（基于 comment 匹配）。
func (m *TrafficMonitor) RemoveInstanceRules(ctx context.Context, inst *manager.Instance) error {
	if inst == nil {
		return fmt.Errorf("instance is nil")
	}
	rules, err := m.listRules(ctx)
	if err != nil {
		return err
	}
	var errs []string
	for _, rule := range rules {
		if rule.Comment == inst.Username {
			handle := strconv.FormatInt(rule.Handle, 10)
			if _, err := m.runner.Run(ctx, nftBinary, "delete", "rule", m.family, m.table, m.chain, "handle", handle); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("delete rule: %s", strings.Join(errs, "; "))
	}
	return nil
}

// UpdateTraffic 读取 nftables 数据并计算增量。
func (m *TrafficMonitor) UpdateTraffic(ctx context.Context, instances []*manager.Instance) ([]client.InstanceTraffic, error) {
	portBytes, err := m.readPortBytes(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	result := make([]client.InstanceTraffic, 0, len(instances))
	for _, inst := range instances {
		if inst == nil {
			continue
		}
		bytesTotal, ok := portBytes[inst.Port]
		if !ok {
			continue
		}
		last := m.lastStats[inst.ID]
		var delta int64
		if last == nil {
			delta = 0
		} else {
			delta = bytesTotal - last.BytesTotal
			if delta < 0 {
				delta = bytesTotal
			}
		}
		m.lastStats[inst.ID] = &PortTraffic{Port: inst.Port, BytesTotal: bytesTotal, Timestamp: now}
		if delta <= 0 {
			continue
		}
		traffic := client.InstanceTraffic{
			InstanceID:    inst.ID,
			BytesUpload:   delta / 2,
			BytesDownload: delta - (delta / 2),
		}
		result = append(result, traffic)
	}
	return result, nil
}

// CleanupInstance 删除缓存数据。
func (m *TrafficMonitor) CleanupInstance(instanceID uint) {
	delete(m.lastStats, instanceID)
}

func (m *TrafficMonitor) readPortBytes(ctx context.Context) (map[int]int64, error) {
	rules, err := m.listRules(ctx)
	if err != nil {
		return nil, err
	}
	traffic := make(map[int]int64, len(rules))
	for _, rule := range rules {
		port, bytes := parseRuleExpr(rule.Expr)
		if port > 0 {
			traffic[port] += bytes
		}
	}
	return traffic, nil
}

func (m *TrafficMonitor) listRules(ctx context.Context) ([]*NFTableRule, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), defaultNFTTimeout)
		defer cancel()
	}
	output, err := m.runner.Run(ctx, nftBinary, "-j", "list", "table", m.family, m.table)
	if err != nil {
		return nil, fmt.Errorf("list nft table: %w", err)
	}
	var parsed NFTablesOutput
	if err := json.Unmarshal(output, &parsed); err != nil {
		return nil, fmt.Errorf("parse nft output: %w", err)
	}
	rules := make([]*NFTableRule, 0, len(parsed.Nftables))
	for _, obj := range parsed.Nftables {
		if obj.Rule != nil && obj.Rule.Chain == m.chain {
			rules = append(rules, obj.Rule)
		}
	}
	return rules, nil
}

func parseRuleExpr(exprs []interface{}) (int, int64) {
	var port int
	var bytes int64
	for _, expr := range exprs {
		exprMap, ok := expr.(map[string]interface{})
		if !ok {
			continue
		}
		if match, ok := exprMap["match"].(map[string]interface{}); ok {
			if right, ok := match["right"]; ok {
				switch v := right.(type) {
				case float64:
					if v > 0 && v <= 65535 {
						port = int(v)
					}
				case map[string]interface{}:
					if payload, ok := v["payload"].(map[string]interface{}); ok {
						if value, ok := payload["value"].(float64); ok {
							port = int(value)
						}
					}
				}
			}
		}
		if counter, ok := exprMap["counter"].(map[string]interface{}); ok {
			if b, ok := counter["bytes"].(float64); ok {
				bytes = int64(b)
			}
		}
	}
	return port, bytes
}
