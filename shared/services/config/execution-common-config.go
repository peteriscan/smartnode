package config

import (
	"github.com/rocket-pool/smartnode/shared/types/config"
)

const (
	// Param IDs
	ecHttpPortID     string = "httpPort"
	ecWsPortID       string = "wsPort"
	ecOpenRpcPortsID string = "openRpcPorts"

	// Defaults
	defaultEcP2pPort     uint16 = 30303
	defaultEcHttpPort    uint16 = 8545
	defaultEcWsPort      uint16 = 8546
	defaultOpenEcApiPort bool   = false
)

// Configuration for the Execution client
type ExecutionCommonConfig struct {
	Title string `yaml:"-"`

	// The HTTP API port
	HttpPort config.Parameter `yaml:"httpPort,omitempty"`

	// The Websocket API port
	WsPort config.Parameter `yaml:"wsPort,omitempty"`

	// Toggle for forwarding the HTTP and Websocket API ports outside of Docker
	OpenRpcPorts config.Parameter `yaml:"openRpcPorts,omitempty"`

	// P2P traffic port
	P2pPort config.Parameter `yaml:"p2pPort,omitempty"`

	// Label for Ethstats
	EthstatsLabel config.Parameter `yaml:"ethstatsLabel,omitempty"`

	// Login info for Ethstats
	EthstatsLogin config.Parameter `yaml:"ethstatsLogin,omitempty"`
}

// Create a new ExecutionCommonConfig struct
func NewExecutionCommonConfig(cfg *RocketPoolConfig, isFallback bool) *ExecutionCommonConfig {

	prefix := ""
	if isFallback {
		prefix = "FALLBACK_"
	}

	title := "Common Execution Client Settings"
	if isFallback {
		title = "Common Fallback Execution Client Settings"
	}

	return &ExecutionCommonConfig{
		Title: title,

		HttpPort: config.Parameter{
			ID:                   ecHttpPortID,
			Name:                 "HTTP Port",
			Description:          "The port your Execution client should use for its HTTP RPC endpoint.",
			Type:                 config.ParameterType_Uint16,
			Default:              map[config.Network]interface{}{config.Network_All: defaultEcHttpPort},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Api, config.ContainerID_Node, config.ContainerID_Watchtower, config.ContainerID_Eth1, config.ContainerID_Eth2},
			EnvironmentVariables: []string{prefix + "EC_HTTP_PORT"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		WsPort: config.Parameter{
			ID:                   ecWsPortID,
			Name:                 "Websocket Port",
			Description:          "The port your Execution client should use for its Websocket RPC endpoint.",
			Type:                 config.ParameterType_Uint16,
			Default:              map[config.Network]interface{}{config.Network_All: defaultEcWsPort},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1, config.ContainerID_Eth2},
			EnvironmentVariables: []string{prefix + "EC_WS_PORT"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		OpenRpcPorts: config.Parameter{
			ID:                   ecOpenRpcPortsID,
			Name:                 "Expose RPC Ports",
			Description:          "Expose the HTTP and Websocket RPC ports to your local network, so other local machines can access your Execution Client's RPC endpoint.",
			Type:                 config.ParameterType_Bool,
			Default:              map[config.Network]interface{}{config.Network_All: defaultOpenEcApiPort},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		P2pPort: config.Parameter{
			ID:                   "p2pPort",
			Name:                 "P2P Port",
			Description:          "The port Geth should use for P2P (blockchain) traffic to communicate with other nodes.",
			Type:                 config.ParameterType_Uint16,
			Default:              map[config.Network]interface{}{config.Network_All: defaultEcP2pPort},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "EC_P2P_PORT"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		EthstatsLabel: config.Parameter{
			ID:                   "ethstatsLabel",
			Name:                 "ETHStats Label",
			Description:          "If you would like to report your Execution client statistics to https://ethstats.net/, enter the label you want to use here.",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "ETHSTATS_LABEL"},
			CanBeBlank:           true,
			OverwriteOnUpgrade:   false,
		},

		EthstatsLogin: config.Parameter{
			ID:                   "ethstatsLogin",
			Name:                 "ETHStats Login",
			Description:          "If you would like to report your Execution client statistics to https://ethstats.net/, enter the login you want to use here.",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "ETHSTATS_LOGIN"},
			CanBeBlank:           true,
			OverwriteOnUpgrade:   false,
		},
	}
}

// Get the parameters for this config
func (cfg *ExecutionCommonConfig) GetParameters() []*config.Parameter {
	return []*config.Parameter{
		&cfg.HttpPort,
		&cfg.WsPort,
		&cfg.OpenRpcPorts,
		&cfg.P2pPort,
		&cfg.EthstatsLabel,
		&cfg.EthstatsLogin,
	}
}

// The the title for the config
func (cfg *ExecutionCommonConfig) GetConfigTitle() string {
	return cfg.Title
}
