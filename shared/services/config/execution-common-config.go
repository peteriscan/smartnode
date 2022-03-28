package config

// Param IDs
const ecHttpPortID string = "httpPort"
const ecWsPortID string = "wsPort"
const ecOpenRpcPortsID string = "openRpcPorts"

// Defaults
const defaultEcHttpPort uint16 = 8545
const defaultEcWsPort uint16 = 8546
const defaultOpenEcApiPort bool = false

// Configuration for the Execution client
type ExecutionCommonConfig struct {
	Title string `yaml:"title,omitempty"`

	// The HTTP API port
	HttpPort Parameter `yaml:"httpPort,omitempty"`

	// The Websocket API port
	WsPort Parameter `yaml:"wsPort,omitempty"`

	// Toggle for forwarding the HTTP and Websocket API ports outside of Docker
	OpenRpcPorts Parameter `yaml:"openRpcPorts,omitempty"`
}

// Create a new ExecutionCommonConfig struct
func NewExecutionCommonConfig(config *RocketPoolConfig, isFallback bool) *ExecutionCommonConfig {

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

		HttpPort: Parameter{
			ID:                   ecHttpPortID,
			Name:                 "HTTP Port",
			Description:          "The port your Execution client should use for its HTTP RPC endpoint.",
			Type:                 ParameterType_Uint16,
			Default:              map[Network]interface{}{Network_All: defaultEcHttpPort},
			AffectsContainers:    []ContainerID{ContainerID_Api, ContainerID_Node, ContainerID_Watchtower, ContainerID_Eth1, ContainerID_Eth2},
			EnvironmentVariables: []string{prefix + "EC_HTTP_PORT"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		WsPort: Parameter{
			ID:                   ecWsPortID,
			Name:                 "Websocket Port",
			Description:          "The port your Execution client should use for its Websocket RPC endpoint.",
			Type:                 ParameterType_Uint16,
			Default:              map[Network]interface{}{Network_All: defaultEcWsPort},
			AffectsContainers:    []ContainerID{ContainerID_Eth1, ContainerID_Eth2},
			EnvironmentVariables: []string{prefix + "EC_WS_PORT"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		OpenRpcPorts: Parameter{
			ID:                   ecOpenRpcPortsID,
			Name:                 "Expose RPC Ports",
			Description:          "Expose the HTTP and Websocket RPC ports to your local network, so other local machines can access your Execution Client's RPC endpoint.",
			Type:                 ParameterType_Bool,
			Default:              map[Network]interface{}{Network_All: defaultOpenEcApiPort},
			AffectsContainers:    []ContainerID{ContainerID_Eth1},
			EnvironmentVariables: []string{},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},
	}
}

// Get the parameters for this config
func (config *ExecutionCommonConfig) GetParameters() []*Parameter {
	return []*Parameter{
		&config.HttpPort,
		&config.WsPort,
		&config.OpenRpcPorts,
	}
}

// The the title for the config
func (config *ExecutionCommonConfig) GetConfigTitle() string {
	return config.Title
}
