package docker

import (
	"fmt"
	"net"

	e2e "github.com/cometbft/cometbft/test/e2e/pkg"

	"github.com/piplabs/story/e2e/types"
	"github.com/piplabs/story/lib/errors"
)

const (
	ipPrefix      = "10.186.73." // See github.com/cometbft/cometbft/test/e2e/pkg for reference
	startIPSuffix = 100
	startPort     = 8000
)

var localhost = net.ParseIP("127.0.0.1") //nolint:gochecknoglobals // Static IP

// NewInfraData returns a new InfrastructureData for the given manifest.
// In addition to normal.
func NewInfraData(manifest types.Manifest) (types.InfrastructureData, error) {
	infd, err := e2e.NewDockerInfrastructureData(manifest.Manifest)
	if err != nil {
		return types.InfrastructureData{}, errors.Wrap(err, "creating docker infrastructure data")
	}

	// IP generator
	ipSuffix := startIPSuffix
	nextInternalIP := func() net.IP {
		defer func() { ipSuffix++ }()
		return net.ParseIP(fmt.Sprintf(ipPrefix+"%d", ipSuffix))
	}

	// Port generator
	port := startPort
	nextPort := func() uint32 {
		defer func() { port++ }()
		return uint32(port)
	}

	for name := range manifest.IliadEVMs() {
		infd.Instances[name] = e2e.InstanceData{
			IPAddress:    nextInternalIP(),
			ExtIPAddress: localhost,
			Port:         nextPort(),
		}
	}

	for _, name := range manifest.AnvilChains {
		infd.Instances[name] = e2e.InstanceData{
			IPAddress:    nextInternalIP(),
			ExtIPAddress: localhost,
			Port:         nextPort(),
		}
	}

	return types.InfrastructureData{
		InfrastructureData: infd,
	}, nil
}
