package template

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed install-agent.sh
var installScriptTemplate string

// InstallScriptData 部署脚本数据
type InstallScriptData struct {
	MasterURL        string // Master 服务器地址
	APIToken         string // 节点 API Token
	AgentVersion     string // Agent 版本号
	NodeName         string // 节点名称
	AgentDownloadURL string // Agent 下载地址模板（包含 {version} 和 {arch} 占位符）
	AgentBinaryURL   string // Agent 二进制下载地址（用于向后兼容）
}

// GenerateInstallScript 生成部署脚本
func GenerateInstallScript(data InstallScriptData) (string, error) {
	tmpl, err := template.New("install").Parse(installScriptTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
