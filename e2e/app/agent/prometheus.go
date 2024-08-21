package agent

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/storyprotocol/iliad/e2e/types"
	"github.com/storyprotocol/iliad/lib/errors"
	"github.com/storyprotocol/iliad/lib/log"
	"github.com/storyprotocol/iliad/lib/netconf"

	_ "embed"
)

type Secrets struct {
	URL  string
	User string
	Pass string
}

const promPort = 26660 // Default metrics port for all iliad apps (from cometBFT)
const gethPromPort = 6060

//go:embed prometheus.yml.tmpl
var promConfigTmpl []byte

func WriteConfig(ctx context.Context, testnet types.Testnet, secrets Secrets) error {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	bz, err := genPromConfig(ctx, testnet, secrets, hostname)
	if err != nil {
		return errors.Wrap(err, "generating prometheus config")
	}

	promFile := filepath.Join(testnet.Dir, "prometheus", "prometheus.yml")
	if err := os.MkdirAll(filepath.Dir(promFile), 0755); err != nil {
		return errors.Wrap(err, "creating prometheus dir")
	}

	if err := os.WriteFile(promFile, bz, 0644); err != nil {
		return errors.Wrap(err, "writing prometheus config")
	}

	return nil
}

func genPromConfig(ctx context.Context, testnet types.Testnet, secrets Secrets, hostname string) ([]byte, error) {
	var nodeTargets []string
	for _, node := range testnet.Nodes {
		// Prometheus is always inside the same docker-compose, so use service names.
		nodeTargets = append(nodeTargets, fmt.Sprintf("%s:%d", node.Name, promPort))
	}

	var evmTargets []string
	for _, iliadEVM := range testnet.IliadEVMs {
		evmTargets = append(evmTargets, fmt.Sprintf("%s:%d", iliadEVM.InstanceName, gethPromPort))
	}

	network := string(testnet.Network)
	if testnet.Network == netconf.Devnet {
		network = fmt.Sprintf("%s-%s", testnet.Name, hostname)
	}

	if secrets.URL == "" {
		log.Warn(ctx, "Prometheus remote URL not set, metrics not being pushed to Grafana cloud", nil)
	} else {
		log.Info(ctx, "Prometheus metrics pushed to Grafana cloud", "network", network)
	}

	data := promTmplData{
		Network:        network,
		Host:           hostname,
		RemoteURL:      secrets.URL,
		RemoteUsername: secrets.User,
		RemotePassword: secrets.Pass,
		ScrapeConfigs: []promScrapConfig{
			{
				JobName:     "iliad",
				MetricsPath: "/metrics",
				targets:     nodeTargets,
			},
			{
				JobName:     "geth",
				MetricsPath: "/debug/metrics/prometheus",
				targets:     evmTargets,
			},
		},
	}

	t, err := template.New("").Parse(string(promConfigTmpl))
	if err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}

	var bz bytes.Buffer
	if err := t.Execute(&bz, data); err != nil {
		return nil, errors.Wrap(err, "executing template")
	}

	return bz.Bytes(), nil
}

type promTmplData struct {
	Network        string            // Used a "network" label to all metrics
	Host           string            // Hostname of the docker host machine
	RemoteURL      string            // URL to the Grafana cloud server
	RemoteUsername string            // Username to the Grafana cloud server
	RemotePassword string            // Password to the Grafana cloud server
	ScrapeConfigs  []promScrapConfig // List of scrape configs
}

type promScrapConfig struct {
	JobName     string
	MetricsPath string
	targets     []string
}

func (c promScrapConfig) Targets() string {
	return strings.Join(c.targets, ",")
}

// ConfigForHost returns a new prometheus agent config with the given host and iliad targets.
//
//	It replaces the iliad targets with provided.
//	It replaces the geth targets with provided.
//	It replaces the host label.
func ConfigForHost(bz []byte, newHost string, iliads []string, geths []string) []byte {
	var iliadTargets []string
	for _, iliad := range iliads {
		iliadTargets = append(iliadTargets, fmt.Sprintf(`"%s:%d"`, iliad, promPort))
	}
	replace := fmt.Sprintf(`[%s] # iliad targets`, strings.Join(iliadTargets, ","))
	bz = regexp.MustCompile(`(?m)\[.*\] # iliad targets$`).
		ReplaceAll(bz, []byte(replace))

	var gethTargets []string
	for _, geth := range geths {
		gethTargets = append(gethTargets, fmt.Sprintf(`"%s:%d"`, geth, gethPromPort))
	}
	replace = fmt.Sprintf(`[%s] # geth targets`, strings.Join(gethTargets, ","))
	bz = regexp.MustCompile(`(?m)\[.*\] # geth targets$`).
		ReplaceAll(bz, []byte(replace))

	bz = regexp.MustCompile(`(?m)host: '.*'$`).
		ReplaceAll(bz, []byte(fmt.Sprintf(`host: '%s'`, newHost)))

	return bz
}
