package config

import (
	"fmt"
	"runtime"
)

// Constants
const (
	akulaTagAmd64         string = "akula:latest"
	akulaTagArm64         string = "akula:latest"
	akulaEventLogInterval int    = 25000
	akulaMaxPeers         uint16 = 100
	akulaStopSignal       string = "SIGTERM"
)

// Configuration for Akula
type AkulaConfig struct {
	Title string `yaml:"-"`

	// Common parameters that Akula doesn't support and should be hidden
	UnsupportedCommonParams []string `yaml:"-"`

	// Compatible consensus clients
	CompatibleConsensusClients []ConsensusClient `yaml:"-"`

	// The max number of events to query in a single event log query
	EventLogInterval int `yaml:"-"`

	// Max number of P2P peers to connect to
	MaxPeers Parameter `yaml:"maxPeers,omitempty"`

	// The Docker Hub tag for Akula
	ContainerTag Parameter `yaml:"containerTag,omitempty"`

	// Custom command line flags
	AdditionalFlags Parameter `yaml:"additionalFlags,omitempty"`
}

// Generates a new Akula configuration
func NewAkulaConfig(config *RocketPoolConfig, isFallback bool) *AkulaConfig {

	prefix := ""
	if isFallback {
		prefix = "FALLBACK_"
	}

	title := "Akula Settings"
	if isFallback {
		title = "Fallback Akula Settings"
	}

	return &AkulaConfig{
		Title: title,

		UnsupportedCommonParams: []string{},

		CompatibleConsensusClients: []ConsensusClient{
			ConsensusClient_Lighthouse,
			ConsensusClient_Nimbus,
			ConsensusClient_Prysm,
			ConsensusClient_Teku,
		},

		EventLogInterval: akulaEventLogInterval,

		MaxPeers: Parameter{
			ID:                   "maxPeers",
			Name:                 "Max Peers",
			Description:          "The maximum number of peers Akula should connect to. This can be lowered to improve performance on low-power systems or constrained networks. We recommend keeping it at 12 or higher.",
			Type:                 ParameterType_Uint16,
			Default:              map[Network]interface{}{Network_All: akulaMaxPeers},
			AffectsContainers:    []ContainerID{ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "EC_MAX_PEERS"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		ContainerTag: Parameter{
			ID:                   "containerTag",
			Name:                 "Container Tag",
			Description:          "The tag name of the Akula container you want to use on Docker Hub.",
			Type:                 ParameterType_String,
			Default:              map[Network]interface{}{Network_All: getAkulaTag()},
			AffectsContainers:    []ContainerID{ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "EC_CONTAINER_TAG"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   true,
		},

		AdditionalFlags: Parameter{
			ID:                   "additionalFlags",
			Name:                 "Additional Flags",
			Description:          "Additional custom command line flags you want to pass to Akula, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
			Type:                 ParameterType_String,
			Default:              map[Network]interface{}{Network_All: ""},
			AffectsContainers:    []ContainerID{ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "EC_ADDITIONAL_FLAGS"},
			CanBeBlank:           true,
			OverwriteOnUpgrade:   false,
		},
	}
}

// Get the container tag for Akula based on the current architecture
func getAkulaTag() string {
	if runtime.GOARCH == "arm64" {
		return akulaTagArm64
	} else if runtime.GOARCH == "amd64" {
		return akulaTagAmd64
	} else {
		panic(fmt.Sprintf("Akula doesn't support architecture %s", runtime.GOARCH))
	}
}

// Get the parameters for this config
func (config *AkulaConfig) GetParameters() []*Parameter {
	return []*Parameter{
		&config.MaxPeers,
		&config.ContainerTag,
		&config.AdditionalFlags,
	}
}

// The the title for the config
func (config *AkulaConfig) GetConfigTitle() string {
	return config.Title
}
