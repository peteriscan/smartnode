package api

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/settings/protocol"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/smartnode/shared/services/config"
	cfgtypes "github.com/rocket-pool/smartnode/shared/types/config"
	"github.com/rocket-pool/smartnode/shared/utils/log"
	"github.com/rocket-pool/smartnode/shared/utils/math"
)

// The fraction of the timeout period to trigger overdue transactions
const TimeoutSafetyFactor int = 2

// Print the gas price and cost of a TX
func PrintAndCheckGasInfo(gasInfo rocketpool.GasInfo, checkThreshold bool, gasThresholdGwei float64, logger log.ColorLogger, maxFeeWei *big.Int, gasLimit uint64) bool {

	// Check the gas threshold if requested
	if checkThreshold {
		gasThresholdWei := math.RoundUp(gasThresholdGwei*eth.WeiPerGwei, 0)
		gasThreshold := new(big.Int).SetUint64(uint64(gasThresholdWei))
		if maxFeeWei.Cmp(gasThreshold) != -1 {
			logger.Printlnf("Current network gas price is %.2f Gwei, which is not lower than the set threshold of %.2f Gwei. "+
				"Aborting the transaction.", eth.WeiToGwei(maxFeeWei), gasThresholdGwei)
			return false
		}
	} else {
		logger.Println("This transaction does not check the gas threshold limit, continuing...")
	}

	// Print the total TX cost
	var gas *big.Int
	var safeGas *big.Int
	if gasLimit != 0 {
		gas = new(big.Int).SetUint64(gasLimit)
		safeGas = gas
	} else {
		gas = new(big.Int).SetUint64(gasInfo.EstGasLimit)
		safeGas = new(big.Int).SetUint64(gasInfo.SafeGasLimit)
	}
	totalGasWei := new(big.Int).Mul(maxFeeWei, gas)
	totalSafeGasWei := new(big.Int).Mul(maxFeeWei, safeGas)
	logger.Printlnf("This transaction will use a max fee of %.6f Gwei, for a total of up to %.6f - %.6f ETH.",
		eth.WeiToGwei(maxFeeWei),
		math.RoundDown(eth.WeiToEth(totalGasWei), 6),
		math.RoundDown(eth.WeiToEth(totalSafeGasWei), 6))

	return true
}

// Print a TX's details to the logger and waits for it to be mined.
func PrintAndWaitForTransaction(cfg *config.RocketPoolConfig, hash common.Hash, ec rocketpool.ExecutionClient, logger log.ColorLogger) error {

	txWatchUrl := cfg.Smartnode.GetTxWatchUrl()
	hashString := hash.String()

	logger.Printlnf("Transaction has been submitted with hash %s.", hashString)
	if txWatchUrl != "" {
		logger.Printlnf("You may follow its progress by visiting:")
		logger.Printlnf("%s/%s\n", txWatchUrl, hashString)
	}
	logger.Println("Waiting for the transaction to be mined...")

	// Wait for the TX to be mined
	if _, err := utils.WaitForTransaction(ec, hash); err != nil {
		return fmt.Errorf("Error mining transaction: %w", err)
	}

	return nil

}

// Gets the event log interval supported by the selected eth1 client
func GetEventLogInterval(cfg *config.RocketPoolConfig) (*big.Int, error) {

	// Get event log interval
	var eventLogInterval *big.Int = nil
	ecMode := cfg.ExecutionClientMode.Value

	switch ecMode {
	case cfgtypes.Mode_External:
		// Use the Geth limit for external clients
		eventLogInterval = big.NewInt(int64(cfg.Geth.EventLogInterval))

	case cfgtypes.Mode_Local:
		switch cfg.ExecutionClient.Value {
		case cfgtypes.ExecutionClient_Geth:
			eventLogInterval = big.NewInt(int64(cfg.Geth.EventLogInterval))
		case cfgtypes.ExecutionClient_Nethermind:
			eventLogInterval = big.NewInt(int64(cfg.Nethermind.EventLogInterval))
		case cfgtypes.ExecutionClient_Besu:
			eventLogInterval = big.NewInt(int64(cfg.Besu.EventLogInterval))
		case cfgtypes.ExecutionClient_Infura:
			eventLogInterval = big.NewInt(int64(cfg.Infura.EventLogInterval))
		case cfgtypes.ExecutionClient_Pocket:
			eventLogInterval = big.NewInt(int64(cfg.Pocket.EventLogInterval))
		default:
			return nil, fmt.Errorf("unknown execution client selected: %v", cfg.ExecutionClient.Value)
		}

	default:
		return nil, fmt.Errorf("unknown execution client mode selected: %v", ecMode)
	}

	return eventLogInterval, nil

}

// True if a transaction is due and needs to bypass the gas threshold
func IsTransactionDue(rp *rocketpool.RocketPool, startTime time.Time) (bool, time.Duration, error) {

	// Get the dissolve timeout
	timeout, err := protocol.GetMinipoolLaunchTimeout(rp, nil)
	if err != nil {
		return false, 0, err
	}

	dueTime := timeout / time.Duration(TimeoutSafetyFactor)
	isDue := time.Since(startTime) > dueTime
	timeUntilDue := time.Until(startTime.Add(dueTime))
	return isDue, timeUntilDue, nil

}
