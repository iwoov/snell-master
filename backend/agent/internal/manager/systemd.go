package manager

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

const serviceTemplate = `[Unit]
Description=Snell Server Instance {{.ID}}
After=network.target

[Service]
Type=simple
User=root
ExecStart={{.BinaryPath}} -c {{.ConfigPath}}
Restart=on-failure
RestartSec=5
StandardOutput=append:{{.LogPath}}
StandardError=append:{{.LogPath}}

[Install]
WantedBy=multi-user.target
`

type serviceData struct {
	ID         uint
	BinaryPath string
	ConfigPath string
	LogPath    string
}

func (m *InstanceManager) generateServiceFile(instance *Instance) (string, error) {
	configPath, logPath := m.generateFilePaths(instance.ID)

	data := serviceData{
		ID:         instance.ID,
		BinaryPath: m.snellBinary,
		ConfigPath: configPath,
		LogPath:    logPath,
	}

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	servicePath := fmt.Sprintf("/etc/systemd/system/snell-instance-%d.service", instance.ID)
	if err := os.WriteFile(servicePath, buf.Bytes(), 0644); err != nil {
		return "", err
	}

	return servicePath, nil
}

func (m *InstanceManager) systemctl(args ...string) error {
	cmd := exec.Command("systemctl", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl %v failed: %w, output: %s", args, err, string(output))
	}
	return nil
}

func (m *InstanceManager) enableAndStartService(id uint) error {
	serviceName := fmt.Sprintf("snell-instance-%d.service", id)
	if err := m.systemctl("daemon-reload"); err != nil {
		return err
	}
	if err := m.systemctl("enable", "--now", serviceName); err != nil {
		return err
	}
	return nil
}

func (m *InstanceManager) stopAndDisableService(id uint) error {
	serviceName := fmt.Sprintf("snell-instance-%d.service", id)
	// Ignore errors if service not found
	_ = m.systemctl("disable", "--now", serviceName)

	servicePath := fmt.Sprintf("/etc/systemd/system/%s", serviceName)
	if err := os.Remove(servicePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return m.systemctl("daemon-reload")
}

func (m *InstanceManager) isServiceActive(id uint) bool {
	serviceName := fmt.Sprintf("snell-instance-%d.service", id)
	cmd := exec.Command("systemctl", "is-active", "--quiet", serviceName)
	return cmd.Run() == nil
}
