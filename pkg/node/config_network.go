package node

import (
	"errors"
	"fmt"
	"time"

	"github.com/bacalhau-project/bacalhau/pkg/lib/validate"
	"github.com/bacalhau-project/bacalhau/pkg/models"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/samber/lo"
)

var supportedNetworks = []string{
	models.NetworkTypeLibp2p,
	models.NetworkTypeNATS,
}

type NetworkConfig struct {
	Type           string
	Libp2pHost     host.Host // only set if using libp2p transport, nil otherwise
	ReconnectDelay time.Duration

	// NATS config for requesters to be reachable by compute nodes
	Port              int
	AdvertisedAddress string
	Orchestrators     []string

	// Storage directory for NATS features that require it
	StoreDir string

	// AuthSecret is a secret string that clients must use to connect. NATS servers
	// must supply this config, while clients can also supply it as the user part
	// of their Orchestrator URL.
	AuthSecret string

	// NATS config for requester nodes to connect with each other
	ClusterName              string
	ClusterPort              int
	ClusterAdvertisedAddress string

	// When using NATS, never set this value unless you are connecting multiple requester
	// nodes together. This should never reference this current running instance (e.g.
	// don't use localhost).
	ClusterPeers []string
}

func (c *NetworkConfig) Validate() error {
	var mErr error
	if validate.IsBlank(c.Type) {
		mErr = errors.Join(mErr, errors.New("missing network type"))
	} else if !lo.Contains(supportedNetworks, c.Type) {
		mErr = errors.Join(mErr, fmt.Errorf("network type %s not in supported values %s", c.Type, supportedNetworks))
	}
	return mErr
}
