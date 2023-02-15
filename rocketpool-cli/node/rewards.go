package node

import (
	"fmt"
	"time"

	"github.com/urfave/cli"

	rprewards "github.com/rocket-pool/smartnode/shared/services/rewards"
	"github.com/rocket-pool/smartnode/shared/services/rocketpool"
	cliutils "github.com/rocket-pool/smartnode/shared/utils/cli"
)

func getRewards(c *cli.Context) error {

	// Get RP client
	rp, err := rocketpool.NewClientFromCtx(c)
	if err != nil {
		return err
	}
	defer rp.Close()

	// Check and assign the EC status
	err = cliutils.CheckClientStatus(rp)
	if err != nil {
		return err
	}

	// Get eligible intervals
	rewardsInfoResponse, err := rp.GetRewardsInfo()
	if err != nil {
		return fmt.Errorf("error getting rewards info: %w", err)
	}

	if !rewardsInfoResponse.Registered {
		fmt.Printf("This node is not currently registered.\n")
		return nil
	}

	// Check for missing Merkle trees with rewards available
	missingIntervals := []rprewards.IntervalInfo{}
	invalidIntervals := []rprewards.IntervalInfo{}
	for _, intervalInfo := range rewardsInfoResponse.InvalidIntervals {
		if !intervalInfo.TreeFileExists {
			fmt.Printf("You are missing the rewards tree file for interval %d.\n", intervalInfo.Index)
			missingIntervals = append(missingIntervals, intervalInfo)
		} else if !intervalInfo.MerkleRootValid {
			fmt.Printf("Your local copy of the rewards tree file for interval %d does not match the canonical one.\n", intervalInfo.Index)
			invalidIntervals = append(invalidIntervals, intervalInfo)
		}
	}

	// Download the Merkle trees for all unclaimed intervals that don't exist
	if len(missingIntervals) > 0 || len(invalidIntervals) > 0 {
		fmt.Println()
		fmt.Printf("%sNOTE: If you would like to regenerate these tree files manually, please answer `n` to the prompt below and run `rocketpool network generate-rewards-tree` before claiming your rewards.%s\n", colorBlue, colorReset)
		if !cliutils.Confirm("Would you like to download all missing rewards tree files now?") {
			fmt.Println("Cancelled.")
			return nil
		}

		// Load the config file
		cfg, _, err := rp.LoadConfig()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		// Download the files
		for _, missingInterval := range missingIntervals {
			fmt.Printf("Downloading interval %d file... ", missingInterval.Index)
			err := rprewards.DownloadRewardsFile(cfg, missingInterval.Index, missingInterval.CID, false)
			if err != nil {
				fmt.Println()
				return err
			}
			fmt.Println("done!")
		}
		for _, invalidInterval := range invalidIntervals {
			fmt.Printf("Downloading interval %d file... ", invalidInterval.Index)
			err := rprewards.DownloadRewardsFile(cfg, invalidInterval.Index, invalidInterval.CID, false)
			if err != nil {
				fmt.Println()
				return err
			}
			fmt.Println("done!")
		}
		fmt.Println()

		// Reload rewards now that the files are in place
		rewardsInfoResponse, err = rp.GetRewardsInfo()
		if err != nil {
			return fmt.Errorf("error getting rewards info: %w", err)
		}
	}

	// Get node RPL rewards status
	rewards, err := rp.NodeRewards()
	if err != nil {
		return err
	}

	fmt.Printf("%sNOTE: Legacy rewards from pre-Redstone are temporarily not being included in the below figures. They will be added back in a future release. We apologize for the inconvenience!%s\n\n", colorYellow, colorReset)

	fmt.Println("=== ETH ===")
	fmt.Printf("You have earned %.4f ETH from the Beacon Chain (including your commissions) so far.\n", rewards.BeaconRewards)
	fmt.Printf("You have claimed %.4f ETH from the Smoothing Pool.\n", rewards.CumulativeEthRewards)
	fmt.Printf("You still have %.4f ETH in unclaimed Smoothing Pool rewards.\n", rewards.UnclaimedEthRewards)

	nextRewardsTime := rewards.LastCheckpoint.Add(rewards.RewardsInterval)
	nextRewardsTimeString := cliutils.GetDateTimeString(uint64(nextRewardsTime.Unix()))
	timeToCheckpointString := time.Until(nextRewardsTime).Round(time.Second).String()

	// Assume 365 days in a year, 24 hours per day
	rplApr := rewards.EstimatedRewards / rewards.TotalRplStake / rewards.RewardsInterval.Hours() * (24 * 365) * 100

	fmt.Println("\n=== RPL ===")
	fmt.Printf("The current rewards cycle started on %s.\n", cliutils.GetDateTimeString(uint64(rewards.LastCheckpoint.Unix())))
	fmt.Printf("It will end on %s (%s from now).\n", nextRewardsTimeString, timeToCheckpointString)

	if rewards.UnclaimedRplRewards > 0 {
		fmt.Printf("You currently have %f unclaimed RPL from staking rewards.\n", rewards.UnclaimedRplRewards)
	}
	if rewards.UnclaimedTrustedRplRewards > 0 {
		fmt.Printf("You currently have %f unclaimed RPL from Oracle DAO duties.\n", rewards.UnclaimedTrustedRplRewards)
	}

	fmt.Println()
	fmt.Printf("Your estimated RPL staking rewards for this cycle: %f RPL (this may change based on network activity).\n", rewards.EstimatedRewards)
	fmt.Printf("Based on your current total stake of %f RPL, this is approximately %.2f%% APR.\n", rewards.TotalRplStake, rplApr)
	fmt.Printf("Your node has received %f RPL staking rewards in total.\n", rewards.CumulativeRplRewards)

	if rewards.Trusted {
		rplTrustedApr := rewards.EstimatedTrustedRplRewards / rewards.TrustedRplBond / rewards.RewardsInterval.Hours() * (24 * 365) * 100

		fmt.Println()
		fmt.Printf("You will receive an estimated %f RPL in rewards for Oracle DAO duties (this may change based on network activity).\n", rewards.EstimatedTrustedRplRewards)
		fmt.Printf("Based on your bond of %f RPL, this is approximately %.2f%% APR.\n", rewards.TrustedRplBond, rplTrustedApr)
		fmt.Printf("Your node has received %f RPL Oracle DAO rewards in total.\n", rewards.CumulativeTrustedRplRewards)
	}

	fmt.Println()
	fmt.Println("You may claim these rewards at any time. You no longer need to claim them within this interval.")

	// Return
	return nil

}
