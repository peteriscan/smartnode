package config

import (
	"github.com/rocket-pool/smartnode/shared/types/config"
)

// Constants
const defaultPocketGatewayMainnet string = "lb/613bb4ae8c124d00353c40a1"
const defaultPocketGatewayPrater string = "lb/6126b4a783e49000343a3a47"
const pocketEventLogInterval int = 25000

// Configuration for Pocket
type PocketConfig struct {
	Title string `yaml:"-"`

	// Common parameters that Pocket doesn't support and should be hidden
	UnsupportedCommonParams []string `yaml:"-"`

	// Compatible consensus clients
	CompatibleConsensusClients []config.ConsensusClient `yaml:"-"`

	// The max number of events to query in a single event log query
	EventLogInterval int `yaml:"-"`

	// The Pocket gateway ID
	GatewayID config.Parameter `yaml:"gatewayID,omitempty"`

	// The Docker Hub tag for Geth
	ContainerTag config.Parameter `yaml:"containerTag,omitempty"`

	// Custom command line flags
	AdditionalFlags config.Parameter `yaml:"additionalFlags,omitempty"`
}

// Generates a new Pocket configuration
func NewPocketConfig(cfg *RocketPoolConfig, isFallback bool) *PocketConfig {

	prefix := ""
	if isFallback {
		prefix = "FALLBACK_"
	}

	title := "Pocket Settings"
	if isFallback {
		title = "Fallback Pocket Settings"
	}

	return &PocketConfig{
		Title: title,

		UnsupportedCommonParams: []string{ecWsPortID},

		CompatibleConsensusClients: []config.ConsensusClient{
			config.ConsensusClient_Lighthouse,
			config.ConsensusClient_Nimbus,
			config.ConsensusClient_Prysm,
			config.ConsensusClient_Teku,
		},

		EventLogInterval: pocketEventLogInterval,

		GatewayID: config.Parameter{
			ID:          "gatewayID",
			Name:        "Gateway ID",
			Description: "If you would like to use a custom gateway for Pocket instead of the default Rocket Pool gateway, enter it here.",
			Type:        config.ParameterType_String,
			Default: map[config.Network]interface{}{
				config.Network_Mainnet: defaultPocketGatewayMainnet,
				config.Network_Prater:  defaultPocketGatewayPrater,
			},
			Regex:                "(^$|^(lb\\/)?[0-9a-zA-Z]{24,}$)",
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "POCKET_GATEWAY_ID"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		ContainerTag: config.Parameter{
			ID:                   "containerTag",
			Name:                 "Container Tag",
			Description:          "The tag name of the Rocket Pool EC Proxy container you want to use on Docker Hub.\nYou should leave this as the default unless you have a good reason to change it.",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: powProxyTag},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "EC_CONTAINER_TAG"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   true,
		},

		AdditionalFlags: config.Parameter{
			ID:                   "additionalFlags",
			Name:                 "Additional Flags",
			Description:          "Additional custom command line flags you want to pass to the EC Proxy, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "EC_ADDITIONAL_FLAGS"},
			CanBeBlank:           true,
			OverwriteOnUpgrade:   false,
		},
	}
}

// Get the parameters for this config
func (cfg *PocketConfig) GetParameters() []*config.Parameter {
	return []*config.Parameter{
		&cfg.GatewayID,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// The the title for the config
func (cfg *PocketConfig) GetConfigTitle() string {
	return cfg.Title
}
