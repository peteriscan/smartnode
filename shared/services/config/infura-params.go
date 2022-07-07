package config

import (
	"github.com/rocket-pool/smartnode/shared/types/config"
)

// Constants
const (
	infuraEventLogInterval int    = 25000
	powProxyStopSignal     string = "SIGTERM"
)

// Configuration for Infura
type InfuraConfig struct {
	Title string `yaml:"-"`

	// Common parameters that Infura doesn't support and should be hidden
	UnsupportedCommonParams []string `yaml:"-"`

	// Compatible consensus clients
	CompatibleConsensusClients []config.ConsensusClient `yaml:"-"`

	// The max number of events to query in a single event log query
	EventLogInterval int `yaml:"-"`

	// The Infura project ID
	ProjectID config.Parameter `yaml:"projectID,omitempty"`

	// The Docker Hub tag for Geth
	ContainerTag config.Parameter `yaml:"containerTag,omitempty"`

	// Custom command line flags
	AdditionalFlags config.Parameter `yaml:"additionalFlags,omitempty"`
}

// Generates a new Infura configuration
func NewInfuraConfig(cfg *RocketPoolConfig, isFallback bool) *InfuraConfig {

	prefix := ""
	if isFallback {
		prefix = "FALLBACK_"
	}

	title := "Infura Settings"
	if isFallback {
		title = "Fallback Infura Settings"
	}

	return &InfuraConfig{
		Title: title,

		CompatibleConsensusClients: []config.ConsensusClient{
			config.ConsensusClient_Lighthouse,
			config.ConsensusClient_Nimbus,
			config.ConsensusClient_Prysm,
			config.ConsensusClient_Teku,
		},

		EventLogInterval: infuraEventLogInterval,

		ProjectID: config.Parameter{
			ID:                   "projectID",
			Name:                 "Project ID",
			Description:          "The ID of your `Ethereum` project in Infura. Note: This is your Project ID, not your Project Secret!",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: ""},
			Regex:                "^[0-9a-fA-F]{32}$",
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "INFURA_PROJECT_ID"},
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
func (cfg *InfuraConfig) GetParameters() []*config.Parameter {
	return []*config.Parameter{
		&cfg.ProjectID,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// The the title for the config
func (cfg *InfuraConfig) GetConfigTitle() string {
	return cfg.Title
}
