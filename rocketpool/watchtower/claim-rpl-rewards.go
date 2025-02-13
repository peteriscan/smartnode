package watchtower

import (
	"fmt"
	"math/big"

	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/rewards"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/urfave/cli"

	"github.com/rocket-pool/smartnode/shared/services"
	"github.com/rocket-pool/smartnode/shared/services/config"
	"github.com/rocket-pool/smartnode/shared/services/wallet"
	"github.com/rocket-pool/smartnode/shared/utils/api"
	"github.com/rocket-pool/smartnode/shared/utils/log"
	"github.com/rocket-pool/smartnode/shared/utils/math"
)

// Claim RPL rewards task
type claimRplRewards struct {
	c       *cli.Context
	log     log.ColorLogger
	cfg     *config.RocketPoolConfig
	w       *wallet.Wallet
	rp      *rocketpool.RocketPool
	enabled bool
}

// Create claim RPL rewards task
func newClaimRplRewards(c *cli.Context, logger log.ColorLogger) (*claimRplRewards, error) {

	// Get services
	cfg, err := services.GetConfig(c)
	if err != nil {
		return nil, err
	}
	w, err := services.GetWallet(c)
	if err != nil {
		return nil, err
	}
	rp, err := services.GetRocketPool(c)
	if err != nil {
		return nil, err
	}

	// Check if auto-claiming is disabled
	isEnabled := true
	gasThreshold := cfg.Smartnode.RplClaimGasThreshold.Value.(float64)
	if gasThreshold == 0 {
		logger.Println("RPL claim gas threshold is set to 0, automatic claims will be disabled.")
		isEnabled = false
	}

	// Return task
	return &claimRplRewards{
		c:       c,
		log:     logger,
		cfg:     cfg,
		w:       w,
		rp:      rp,
		enabled: isEnabled,
	}, nil

}

// Claim RPL rewards
func (t *claimRplRewards) run() error {

	// Check to see if autoclaim is disabled
	if !t.enabled {
		return nil
	}

	// Wait for eth client to sync
	if err := services.WaitEthClientSynced(t.c, true); err != nil {
		return err
	}

	// Get node account
	nodeAccount, err := t.w.GetNodeAccount()
	if err != nil {
		return err
	}

	// Check node trusted status
	nodeTrusted, err := trustednode.GetMemberExists(t.rp, nodeAccount.Address, nil)
	if err != nil {
		return err
	}
	if !nodeTrusted {
		return nil
	}

	// Log
	t.log.Println("Checking for RPL rewards to claim...")

	// Check for rewards
	rewardsAmountWei, err := rewards.GetTrustedNodeClaimRewardsAmount(t.rp, nodeAccount.Address, nil)
	if err != nil {
		return err
	}
	if rewardsAmountWei.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	// Log
	t.log.Printlnf("%.6f RPL is available to claim...", math.RoundDown(eth.WeiToEth(rewardsAmountWei), 6))

	// Get transactor
	opts, err := t.w.GetNodeAccountTransactor()
	if err != nil {
		return err
	}

	// Get the gas limit
	gasInfo, err := rewards.EstimateClaimTrustedNodeRewardsGas(t.rp, opts)
	if err != nil {
		return fmt.Errorf("Could not estimate the gas required to claim RPL: %w", err)
	}

	// Print the gas info
	maxFee := eth.GweiToWei(WatchtowerMaxFee)
	if !api.PrintAndCheckGasInfo(gasInfo, false, 0, t.log, maxFee, 0) {
		return nil
	}

	// Set the gas settings
	opts.GasFeeCap = maxFee
	opts.GasTipCap = eth.GweiToWei(WatchtowerMaxPriorityFee)
	opts.GasLimit = gasInfo.SafeGasLimit

	// Claim rewards
	hash, err := rewards.ClaimTrustedNodeRewards(t.rp, opts)
	if err != nil {
		return err
	}

	// Print TX info and wait for it to be mined
	err = api.PrintAndWaitForTransaction(t.cfg, hash, t.rp.Client, t.log)
	if err != nil {
		return err
	}

	// Log & return
	t.log.Printlnf("Successfully claimed %.6f RPL in rewards.", math.RoundDown(eth.WeiToEth(rewardsAmountWei), 6))
	return nil

}
