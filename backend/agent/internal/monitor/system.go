package monitor

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

const metricTimeout = 2 * time.Second

var errMetricUnsupported = errors.New("metric unsupported")

type metricFunc func(ctx context.Context) (float64, error)

type cpuSample struct {
	idle  uint64
	total uint64
}

// SystemMonitor 负责采集节点资源使用情况。
type SystemMonitor struct {
	cpuUsage       int
	memoryUsage    int
	diskUsage      int
	lastUpdated    time.Time
	lastCPUSample  cpuSample
	cpuSampleValid bool

	cpuFn  metricFunc
	memFn  metricFunc
	diskFn metricFunc
}

// NewSystemMonitor 创建 SystemMonitor，允许通过选项注入自定义采集实现。
func NewSystemMonitor(opts ...SystemOption) *SystemMonitor {
	monitor := &SystemMonitor{}
	monitor.cpuFn = monitor.defaultCPUUsage
	monitor.memFn = defaultMemoryUsage
	monitor.diskFn = defaultDiskUsage
	for _, opt := range opts {
		opt(monitor)
	}
	return monitor
}

// SystemOption 允许测试注入采集函数。
type SystemOption func(*SystemMonitor)

func WithCPUProvider(fn metricFunc) SystemOption    { return func(m *SystemMonitor) { m.cpuFn = fn } }
func WithMemoryProvider(fn metricFunc) SystemOption { return func(m *SystemMonitor) { m.memFn = fn } }
func WithDiskProvider(fn metricFunc) SystemOption   { return func(m *SystemMonitor) { m.diskFn = fn } }

// Update 刷新所有指标。
func (m *SystemMonitor) Update(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), metricTimeout)
		defer cancel()
	}

	cpuPercent, err := m.cpuFn(ctx)
	if err != nil && !errors.Is(err, errMetricUnsupported) {
		return fmt.Errorf("cpu usage: %w", err)
	}
	memPercent, err := m.memFn(ctx)
	if err != nil && !errors.Is(err, errMetricUnsupported) {
		return fmt.Errorf("memory usage: %w", err)
	}
	diskPercent, err := m.diskFn(ctx)
	if err != nil && !errors.Is(err, errMetricUnsupported) {
		return fmt.Errorf("disk usage: %w", err)
	}

	m.cpuUsage = clampPercent(cpuPercent)
	m.memoryUsage = clampPercent(memPercent)
	m.diskUsage = clampPercent(diskPercent)
	m.lastUpdated = time.Now()
	return nil
}

func (m *SystemMonitor) CPUUsage() int          { return m.cpuUsage }
func (m *SystemMonitor) MemoryUsage() int       { return m.memoryUsage }
func (m *SystemMonitor) DiskUsage() int         { return m.diskUsage }
func (m *SystemMonitor) LastUpdated() time.Time { return m.lastUpdated }
func (m *SystemMonitor) GetCPUUsage() int       { return m.cpuUsage }
func (m *SystemMonitor) GetMemoryUsage() int    { return m.memoryUsage }
func (m *SystemMonitor) GetDiskUsage() int      { return m.diskUsage }

func (m *SystemMonitor) defaultCPUUsage(context.Context) (float64, error) {
	if runtime.GOOS != "linux" {
		return 0, errMetricUnsupported
	}
	sample, err := readLinuxCPUSample()
	if err != nil {
		return 0, err
	}
	if !m.cpuSampleValid {
		m.lastCPUSample = sample
		m.cpuSampleValid = true
		return 0, nil
	}
	deltaIdle := float64(sample.idle - m.lastCPUSample.idle)
	deltaTotal := float64(sample.total - m.lastCPUSample.total)
	m.lastCPUSample = sample
	if deltaTotal <= 0 {
		return 0, nil
	}
	return (1 - deltaIdle/deltaTotal) * 100, nil
}

func defaultMemoryUsage(context.Context) (float64, error) {
	if runtime.GOOS != "linux" {
		return 0, errMetricUnsupported
	}
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var total, available float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			total = parseMeminfoValue(line)
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			available = parseMeminfoValue(line)
		}
	}
	if total == 0 {
		return 0, fmt.Errorf("meminfo missing total")
	}
	used := total - available
	if used < 0 {
		used = 0
	}
	return (used / total) * 100, nil
}

func defaultDiskUsage(context.Context) (float64, error) {
	var stat unix.Statfs_t
	if err := unix.Statfs("/", &stat); err != nil {
		return 0, err
	}
	total := float64(stat.Blocks) * float64(stat.Bsize)
	free := float64(stat.Bavail) * float64(stat.Bsize)
	if total == 0 {
		return 0, fmt.Errorf("statfs total is zero")
	}
	used := total - free
	if used < 0 {
		used = 0
	}
	return (used / total) * 100, nil
}

func clampPercent(value float64) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return int(value + 0.5)
}

func parseMeminfoValue(line string) float64 {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0
	}
	value, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return 0
	}
	return value
}

func readLinuxCPUSample() (cpuSample, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return cpuSample{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return cpuSample{}, fmt.Errorf("empty /proc/stat")
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) < 5 {
		return cpuSample{}, fmt.Errorf("invalid cpu line")
	}
	var values []uint64
	for _, field := range fields[1:] {
		v, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return cpuSample{}, err
		}
		values = append(values, v)
	}
	var total uint64
	for _, v := range values {
		total += v
	}
	idle := values[3]
	return cpuSample{idle: idle, total: total}, nil
}
