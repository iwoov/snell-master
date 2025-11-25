package template

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/iwoov/snell-master/backend/master/internal/model"
)

// SurgeGenerator 将模板渲染为完整的 Surge 配置。
type SurgeGenerator struct {
	template string
}

// NewSurgeGenerator 创建生成器。
func NewSurgeGenerator(tpl string) *SurgeGenerator {
	return &SurgeGenerator{template: tpl}
}

// Generate 根据用户、节点和实例数据生成文本。
func (g *SurgeGenerator) Generate(user *model.User, nodes []model.Node, instances []model.SnellInstance) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user is required")
	}
	if g.template == "" {
		return "", fmt.Errorf("template content is empty")
	}

	nodeMap := make(map[uint]*model.Node)
	for i := range nodes {
		nodeMap[nodes[i].ID] = &nodes[i]
	}

	var (
		nodeListBuf  bytes.Buffer
		nodeNameList []string
	)

	var primaryInstance *model.SnellInstance
	var primaryNode *model.Node

	for _, inst := range instances {
		node := nodeMap[inst.NodeID]
		if node == nil {
			continue
		}
		if primaryInstance == nil {
			primaryInstance = &inst
			primaryNode = node
		}
		name := fmt.Sprintf("%s-%d", node.Name, inst.Port)
		if emoji := countryEmoji(node.CountryCode); emoji != "" {
			name = emoji + " " + name
		}
		nodeNameList = append(nodeNameList, name)
		nodeLine := fmt.Sprintf("%s = snell, %s, %d, psk=%s, version=%d", name, node.Endpoint, inst.Port, inst.PSK, inst.Version)
		if inst.Obfs != "" {
			nodeLine += ", obfs=" + inst.Obfs
		}
		nodeListBuf.WriteString(nodeLine)
		nodeListBuf.WriteString("\n")
	}

	country := ""
	emoji := ""
	server := ""
	port := ""
	psk := ""
	if primaryNode != nil {
		country = primaryNode.CountryCode
		emoji = countryEmoji(country)
		server = primaryNode.Endpoint
	}
	if primaryInstance != nil {
		port = fmt.Sprintf("%d", primaryInstance.Port)
		psk = primaryInstance.PSK
	}

	replacer := strings.NewReplacer(
		"{{username}}", user.Username,
		"{{email}}", user.Email,
		"{{node_list}}", strings.TrimSpace(nodeListBuf.String()),
		"{{node_names}}", strings.Join(nodeNameList, ", "),
		"{{country}}", country,
		"{{emoji}}", emoji,
		"{{node_name}}", func() string {
			if primaryNode != nil {
				return primaryNode.Name
			}
			return ""
		}(),
		"{{server}}", server,
		"{{port}}", port,
		"{{psk}}", psk,
	)

	return replacer.Replace(g.template), nil
}

func countryEmoji(code string) string {
	if len(code) != 2 {
		return ""
	}
	code = strings.ToUpper(code)
	runes := []rune{}
	for _, ch := range code {
		if ch < 'A' || ch > 'Z' {
			return ""
		}
		runes = append(runes, rune(0x1F1E6+(ch-'A')))
	}
	return string(runes)
}
