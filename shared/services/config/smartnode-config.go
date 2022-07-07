package config

import (
	"path/filepath"

	"github.com/rocket-pool/smartnode/shared"
	"github.com/rocket-pool/smartnode/shared/types/config"
)

// Constants
const (
	smartnodeTag        string = "rocketpool/smartnode:v" + shared.RocketPoolVersion
	powProxyTag         string = "rocketpool/smartnode-pow-proxy:v" + shared.RocketPoolVersion
	pruneProvisionerTag string = "rocketpool/eth1-prune-provision:v0.0.1"
	ecMigratorTag       string = "rocketpool/ec-migrator:v1.0.0"
	NetworkID           string = "network"
	ProjectNameID       string = "projectName"
	SnapshotID          string = "rocketpool-dao.eth"
)

// Defaults
const defaultProjectName string = "rocketpool"

// Configuration for the Smartnode
type SmartnodeConfig struct {
	Title string `yaml:"-"`

	// The parent config
	parent *RocketPoolConfig `yaml:"-"`

	////////////////////////////
	// User-editable settings //
	////////////////////////////

	// Docker container prefix
	ProjectName config.Parameter `yaml:"projectName,omitempty"`

	// The path of the data folder where everything is stored
	DataPath config.Parameter `yaml:"dataPath,omitempty"`

	// Which network we're on
	Network config.Parameter `yaml:"network,omitempty"`

	// Manual max fee override
	ManualMaxFee config.Parameter `yaml:"manualMaxFee,omitempty"`

	// Manual priority fee override
	PriorityFee config.Parameter `yaml:"priorityFee,omitempty"`

	// Threshold for auto RPL claims
	RplClaimGasThreshold config.Parameter `yaml:"rplClaimGasThreshold,omitempty"`

	// Threshold for auto minipool stakes
	MinipoolStakeGasThreshold config.Parameter `yaml:"minipoolStakeGasThreshold,omitempty"`

	///////////////////////////
	// Non-editable settings //
	///////////////////////////

	// The URL to provide the user so they can follow pending transactions
	txWatchUrl map[config.Network]string `yaml:"-"`

	// The URL to use for staking rETH
	stakeUrl map[config.Network]string `yaml:"-"`

	// The map of networks to execution chain IDs
	chainID map[config.Network]uint `yaml:"-"`

	// The path within the daemon Docker container of the wallet file
	walletPath string `yaml:"-"`

	// The path within the daemon Docker container of the wallet's password file
	passwordPath string `yaml:"-"`

	// The path within the daemon Docker container of the validator key folder
	validatorKeychainPath string `yaml:"-"`

	// The path that custom validator keys will be stored (ones for minipools that aren't derived from the node wallet)
	customKeyRecoverPath string `yaml:"-"`

	// The path of the password file for custom validator keys
	customKeyPasswordFilePath string `yaml:"-"`

	// The contract address of RocketStorage
	storageAddress map[config.Network]string `yaml:"-"`

	// The contract address of the 1inch oracle
	oneInchOracleAddress map[config.Network]string `yaml:"-"`

	// The contract address of the RPL token
	rplTokenAddress map[config.Network]string `yaml:"-"`

	// The contract address of the RPL faucet
	rplFaucetAddress map[config.Network]string `yaml:"-"`

	// The contract address for Snapshot delegation
	snapshotDelegationAddress map[config.Network]string `yaml:"-"`
}

// Generates a new Smartnode configuration
func NewSmartnodeConfig(cfg *RocketPoolConfig) *SmartnodeConfig {

	return &SmartnodeConfig{
		Title:  "Smartnode Settings",
		parent: cfg,

		ProjectName: config.Parameter{
			ID:                   ProjectNameID,
			Name:                 "Project Name",
			Description:          "This is the prefix that will be attached to all of the Docker containers managed by the Smartnode.",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: defaultProjectName},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Api, config.ContainerID_Node, config.ContainerID_Watchtower, config.ContainerID_Eth1, config.ContainerID_Eth2, config.ContainerID_Validator, config.ContainerID_Grafana, config.ContainerID_Prometheus, config.ContainerID_Exporter},
			EnvironmentVariables: []string{"COMPOSE_PROJECT_NAME"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		DataPath: config.Parameter{
			ID:                   "dataPath",
			Name:                 "Data Path",
			Description:          "The absolute path of the `data` folder that contains your node wallet's encrypted file, the password for your node wallet, and all of the validator keys for your minipools. You may use environment variables in this string.",
			Type:                 config.ParameterType_String,
			Default:              map[config.Network]interface{}{config.Network_All: getDefaultDataDir(cfg)},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Api, config.ContainerID_Node, config.ContainerID_Watchtower, config.ContainerID_Validator},
			EnvironmentVariables: []string{"ROCKETPOOL_DATA_FOLDER"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		Network: config.Parameter{
			ID:                   NetworkID,
			Name:                 "Network",
			Description:          "The Ethereum network you want to use - select Prater Testnet to practice with fake ETH, or Mainnet to stake on the real network using real ETH.",
			Type:                 config.ParameterType_Choice,
			Default:              map[config.Network]interface{}{config.Network_All: config.Network_Mainnet},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Api, config.ContainerID_Node, config.ContainerID_Watchtower, config.ContainerID_Eth1, config.ContainerID_Eth2, config.ContainerID_Validator},
			EnvironmentVariables: []string{"NETWORK"},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
			Options: []config.ParameterOption{{
				Name:        "Ethereum Mainnet",
				Description: "This is the real Ethereum main network, using real ETH and real RPL to make real validators.",
				Value:       config.Network_Mainnet,
			}, {
				Name:        "Prater Testnet",
				Description: "This is the Prater test network, using free fake ETH and free fake RPL to make fake validators.\nUse this if you want to practice running the Smartnode in a free, safe environment before moving to mainnet.",
				Value:       config.Network_Prater,
			}},
		},

		ManualMaxFee: config.Parameter{
			ID:                   "manualMaxFee",
			Name:                 "Manual Max Fee",
			Description:          "Set this if you want all of the Smartnode's transactions to use this specific max fee value (in gwei), which is the most you'd be willing to pay (*including the priority fee*).\n\nA value of 0 will show you the current suggested max fee based on the current network conditions and let you specify it each time you do a transaction.\n\nAny other value will ignore the recommended max fee and explicitly use this value instead.\n\nThis applies to automated transactions (such as claiming RPL and staking minipools) as well.",
			Type:                 config.ParameterType_Float,
			Default:              map[config.Network]interface{}{config.Network_All: float64(0)},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Node, config.ContainerID_Watchtower},
			EnvironmentVariables: []string{},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		PriorityFee: config.Parameter{
			ID:                   "priorityFee",
			Name:                 "Priority Fee",
			Description:          "The default value for the priority fee (in gwei) for all of your transactions. This describes how much you're willing to pay *above the network's current base fee* - the higher this is, the more ETH you give to the miners for including your transaction, which generally means it will be mined faster (as long as your max fee is sufficiently high to cover the current network conditions).\n\nMust be larger than 0.",
			Type:                 config.ParameterType_Float,
			Default:              map[config.Network]interface{}{config.Network_All: float64(2)},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Node, config.ContainerID_Watchtower},
			EnvironmentVariables: []string{},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		RplClaimGasThreshold: config.Parameter{
			ID:                   "rplClaimGasThreshold",
			Name:                 "RPL Claim Gas Threshold",
			Description:          "Automatic RPL rewards claims will use the `Rapid` suggestion from the gas estimator, based on current network conditions. This threshold is a limit (in gwei) you can put on that suggestion; your node will not try to claim RPL rewards automatically until the suggestion is below this limit.",
			Type:                 config.ParameterType_Float,
			Default:              map[config.Network]interface{}{config.Network_All: float64(150)},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Node, config.ContainerID_Watchtower},
			EnvironmentVariables: []string{},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		MinipoolStakeGasThreshold: config.Parameter{
			ID:   "minipoolStakeGasThreshold",
			Name: "Minipool Stake Gas Threshold",
			Description: "Once a newly created minipool passes the scrub check and is ready to perform its second 16 ETH deposit (the `stake` transaction), your node will try to do so automatically using the `Rapid` suggestion from the gas estimator as its max fee. This threshold is a limit (in gwei) you can put on that suggestion; your node will not `stake` the new minipool until the suggestion is below this limit.\n\n" +
				"Note that to ensure your minipool does not get dissolved, the node will ignore this limit and automatically execute the `stake` transaction at whatever the suggested fee happens to be once too much time has passed since its first deposit (currently 7 days).",
			Type:                 config.ParameterType_Float,
			Default:              map[config.Network]interface{}{config.Network_All: float64(150)},
			AffectsContainers:    []config.ContainerID{config.ContainerID_Node},
			EnvironmentVariables: []string{},
			CanBeBlank:           false,
			OverwriteOnUpgrade:   false,
		},

		txWatchUrl: map[config.Network]string{
			config.Network_Mainnet: "https://etherscan.io/tx",
			config.Network_Prater:  "https://goerli.etherscan.io/tx",
		},

		stakeUrl: map[config.Network]string{
			config.Network_Mainnet: "https://stake.rocketpool.net",
			config.Network_Prater:  "https://testnet.rocketpool.net",
		},

		chainID: map[config.Network]uint{
			config.Network_Mainnet: 1, // Mainnet
			config.Network_Prater:  5, // Goerli
		},

		walletPath: "/.rocketpool/data/wallet",

		passwordPath: "/.rocketpool/data/password",

		validatorKeychainPath: "/.rocketpool/data/validators",

		customKeyRecoverPath: "/.rocketpool/data/custom-keys",

		customKeyPasswordFilePath: "/.rocketpool/data/custom-key-passwords",

		storageAddress: map[config.Network]string{
			config.Network_Mainnet: "0x1d8f8f00cfa6758d7bE78336684788Fb0ee0Fa46",
			config.Network_Prater:  "0xd8Cd47263414aFEca62d6e2a3917d6600abDceB3",
		},

		oneInchOracleAddress: map[config.Network]string{
			config.Network_Mainnet: "0x07D91f5fb9Bf7798734C3f606dB065549F6893bb",
			config.Network_Prater:  "0x4eDC966Df24264C9C817295a0753804EcC46Dd22",
		},

		rplTokenAddress: map[config.Network]string{
			config.Network_Mainnet: "0xb4efd85c19999d84251304bda99e90b92300bd93",
			config.Network_Prater:  "0xb4efd85c19999d84251304bda99e90b92300bd93",
		},

		rplFaucetAddress: map[config.Network]string{
			config.Network_Mainnet: "",
			config.Network_Prater:  "0x95D6b8E2106E3B30a72fC87e2B56ce15E37853F9",
		},

		snapshotDelegationAddress: map[config.Network]string{
			config.Network_Mainnet: "0x469788fE6E9E9681C6ebF3bF78e7Fd26Fc015446",
			config.Network_Prater:  "0xD0897D68Cd66A710dDCecDe30F7557972181BEDc",
		},
	}

}

// Get the parameters for this config
func (cfg *SmartnodeConfig) GetParameters() []*config.Parameter {
	return []*config.Parameter{
		&cfg.Network,
		&cfg.ProjectName,
		&cfg.DataPath,
		&cfg.ManualMaxFee,
		&cfg.PriorityFee,
		&cfg.RplClaimGasThreshold,
		&cfg.MinipoolStakeGasThreshold,
	}
}

// Getters for the non-editable parameters

func (cfg *SmartnodeConfig) GetTxWatchUrl() string {
	return cfg.txWatchUrl[cfg.Network.Value.(config.Network)]
}

func (cfg *SmartnodeConfig) GetStakeUrl() string {
	return cfg.stakeUrl[cfg.Network.Value.(config.Network)]
}

func (cfg *SmartnodeConfig) GetChainID() uint {
	return cfg.chainID[cfg.Network.Value.(config.Network)]
}

func (cfg *SmartnodeConfig) GetWalletPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), "wallet")
	} else {
		return cfg.walletPath
	}
}

func (cfg *SmartnodeConfig) GetPasswordPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), "password")
	} else {
		return cfg.passwordPath
	}
}

func (cfg *SmartnodeConfig) GetValidatorKeychainPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), "validators")
	} else {
		return cfg.validatorKeychainPath
	}
}

func (cfg *SmartnodeConfig) GetCustomKeyPath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), "custom-keys")
	} else {
		return cfg.customKeyRecoverPath
	}
}

func (cfg *SmartnodeConfig) GetCustomKeyPasswordFilePath() string {
	if cfg.parent.IsNativeMode {
		return filepath.Join(cfg.DataPath.Value.(string), "custom-key-passwords")
	} else {
		return cfg.customKeyPasswordFilePath
	}
}

func (cfg *SmartnodeConfig) GetStorageAddress() string {
	return cfg.storageAddress[cfg.Network.Value.(config.Network)]
}

func (cfg *SmartnodeConfig) GetOneInchOracleAddress() string {
	return cfg.oneInchOracleAddress[cfg.Network.Value.(config.Network)]
}

func (cfg *SmartnodeConfig) GetRplTokenAddress() string {
	return cfg.rplTokenAddress[cfg.Network.Value.(config.Network)]
}

func (cfg *SmartnodeConfig) GetRplFaucetAddress() string {
	return cfg.rplFaucetAddress[cfg.Network.Value.(config.Network)]
}

func (cfg *SmartnodeConfig) GetSnapshotDelegationAddress() string {
	return cfg.snapshotDelegationAddress[cfg.Network.Value.(config.Network)]
}

func (cfg *SmartnodeConfig) GetSmartnodeContainerTag() string {
	return smartnodeTag
}

func (cfg *SmartnodeConfig) GetPowProxyContainerTag() string {
	return powProxyTag
}

func (cfg *SmartnodeConfig) GetPruneProvisionerContainerTag() string {
	return pruneProvisionerTag
}

func (cfg *SmartnodeConfig) GetEcMigratorContainerTag() string {
	return ecMigratorTag
}

func (cfg *SmartnodeConfig) GetVotingSnapshotID() [32]byte {
	// So the contract wants a Keccak'd hash of the voting ID, but Snapshot's service wants ASCII so it can display the ID in plain text; we have to do this to make it play nicely with Snapshot
	buffer := [32]byte{}
	idBytes := []byte(SnapshotID)
	copy(buffer[0:], idBytes)
	return buffer
}

// The the title for the config
func (cfg *SmartnodeConfig) GetConfigTitle() string {
	return cfg.Title
}

func getDefaultDataDir(config *RocketPoolConfig) string {
	return filepath.Join(config.RocketPoolDirectory, "data")
}
