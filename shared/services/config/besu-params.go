package config

import (
	"fmt"
	"runtime"

	"github.com/pbnjay/memory"
	"github.com/rocket-pool/smartnode/shared/types/config"
)

// Constants
const (
	besuTagAmd64         string = "hyperledger/besu:22.4.3-openjdk-latest"
	besuTagArm64         string = "hyperledger/besu:22.4.3-openjdk-latest"
	besuEventLogInterval int    = 25000
	besuMaxPeers         uint16 = 25
	besuStopSignal       string = "SIGTERM"
)

// Configuration for Besu
type BesuConfig struct {
	Title string `yaml:"-"`

	// Common parameters that Besu doesn't support and should be hidden
	UnsupportedCommonParams []string `yaml:"-"`

	// Compatible consensus clients
	CompatibleConsensusClients []config.ConsensusClient `yaml:"-"`

	// Max number of P2P peers to connect to
	JvmHeapSize config.Parameter `yaml:"jvmHeapSize,omitempty"`

	// The max number of events to query in a single event log query
	EventLogInterval int `yaml:"-"`

	// Max number of P2P peers to connect to
	MaxPeers config.Parameter `yaml:"maxPeers,omitempty"`

	// Historical state block regeneration limit
	MaxBackLayers config.Parameter `yaml:"maxBackLayers,omitempty"`

	// The Docker Hub tag for Besu
	ContainerTag config.Parameter `yaml:"containerTag,omitempty"`

	// Custom command line flags
	AdditionalFlags config.Parameter `yaml:"additionalFlags,omitempty"`
}

// Generates a new Besu configuration
func NewBesuConfig(cfg *RocketPoolConfig, isFallback bool) *BesuConfig {

	prefix := ""
	if isFallback {
		prefix = "FALLBACK_"
	}

	title := "Besu Settings"
	if isFallback {
		title = "Fallback Besu Settings"
	}

	return &BesuConfig{
		Title: title,

		UnsupportedCommonParams: []string{},

		CompatibleConsensusClients: []config.ConsensusClient{
			config.ConsensusClient_Lighthouse,
			config.ConsensusClient_Nimbus,
			config.ConsensusClient_Prysm,
			config.ConsensusClient_Teku,
		},

		EventLogInterval: nethermindEventLogInterval,

		JvmHeapSize: config.Parameter{
			ID:                   "jvmHeapSize",
			Name:                 "JVM Heap Size",
			Description:          "The max amount of RAM, in MB, that Besu's JVM should limit itself to. Setting this lower will cause Besu to use less RAM, though it will always use more than this limit.\n\nUse 0 for automatic allocation.",
			Type:                 config.ParameterType_Uint,
			Default:              map[config.Network]interface{}{config.Network_All: getBesuHeapSize()},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "BESU_JVM_HEAP_SIZE"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		MaxPeers: config.Parameter{
			ID:                   "maxPeers",
			Name:                 "Max Peers",
			Description:          "The maximum number of peers Besu should connect to. This can be lowered to improve performance on low-power systems or constrained networks. We recommend keeping it at 12 or higher.",
			Type:                 config.ParameterType_Uint16,
			Default:              map[config.Network]interface{}{config.Network_All: besuMaxPeers},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "EC_MAX_PEERS"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		MaxBackLayers: config.Parameter{
			ID:                   "maxBackLayers",
			Name:                 "Historical Block Replay Limit",
			Description:          "Besu has the ability to revisit the state of any historical block on the chain by \"replaying\" all of the previous blocks to get back to the target. This limit controls how many blocks you can replay - in other words, how far back Besu can go in time. Normal Execution client processing will be paused while a replay is in progress.\n\n[orange]NOTE: If you try to replay a state from a long time ago, it may take Besu several minutes to rebuild the state!",
			Type:                 config.ParameterType_Uint,
			Default:              map[config.Network]interface{}{config.Network_All: uint64(512)},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "BESU_MAX_BACK_LAYERS"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		ContainerTag: config.Parameter{
			ID:                   "containerTag",
			Name:                 "Container Tag",
			Description:          "The tag name of the Besu container you want to use on Docker Hub.",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: getBesuTag()},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "EC_CONTAINER_TAG"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   true,
		},

		AdditionalFlags: config.Parameter{
			ID:                   "additionalFlags",
			Name:                 "Additional Flags",
			Description:          "Additional custom command line flags you want to pass to Besu, to take advantage of other settings that the Smartnode's configuration doesn't cover.",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: ""},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Eth1},
			EnvironmentVariables: []string{prefix + "EC_ADDITIONAL_FLAGS"},
			CanBeBlank:           true,
			OverwriteOnUpgrade:   false,
		},
	}
}

// Get the container tag for Besu based on the current architecture
func getBesuTag() string {
	if runtime.GOARCH == "arm64" {
		return besuTagArm64
	} else if runtime.GOARCH == "amd64" {
		return besuTagAmd64
	} else {
		panic(fmt.Sprintf("Besu doesn't support architecture %s", runtime.GOARCH))
	}
}

// Get the recommended heap size for Besu
func getBesuHeapSize() uint64 {
	totalMemoryGB := memory.TotalMemory() / 1024 / 1024 / 1024
	if totalMemoryGB < 9 {
		return 512
	}
	return 0
}

// Get the parameters for this config
func (cfg *BesuConfig) GetParameters() []*config.Parameter {
	return []*config.Parameter{
		&cfg.JvmHeapSize,
		&cfg.MaxPeers,
		&cfg.MaxBackLayers,
		&cfg.ContainerTag,
		&cfg.AdditionalFlags,
	}
}

// The the title for the config
func (cfg *BesuConfig) GetConfigTitle() string {
	return cfg.Title
}
