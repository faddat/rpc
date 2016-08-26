package database

import (
	// Stdlib
	"encoding/json"
	"strconv"
	"strings"

	// RPC
	"github.com/go-steem/rpc/types"
	"github.com/go-steem/rpc/types/simpletypes"
)

type Config struct {
	SteemitBlockchainHardforkVersion string `json:"STEEMIT_BLOCKCHAIN_HARDFORK_VERSION"`
	SteemitBlockchainVersion         string `json:"STEEMIT_BLOCKCHAIN_VERSION"`
	SteemitBlockInterval             uint   `json:"STEEMIT_BLOCK_INTERVAL"`
}

type DynamicGlobalProperties struct {
	Time                     *simpletypes.Time `json:"time"`
	TotalPow                 *simpletypes.Int  `json:"total_pow"`
	NumPowWitnesses          *simpletypes.Int  `json:"num_pow_witnesses"`
	CurrentReserveRatio      *simpletypes.Int  `json:"current_reserve_ratio"`
	Id                       *simpletypes.ID   `json:"id"`
	CurrentSupply            string            `json:"current_supply"`
	CurrentSBDSupply         string            `json:"current_sbd_supply"`
	MaximumBlockSize         *simpletypes.Int  `json:"maximum_block_size"`
	RecentSlotsFilled        *simpletypes.Int  `json:"recent_slots_filled"`
	CurrentWitness           string            `json:"current_witness"`
	TotalRewardShares2       *simpletypes.Int  `json:"total_reward_shares2"`
	AverageBlockSize         *simpletypes.Int  `json:"average_block_size"`
	CurrentAslot             *simpletypes.Int  `json:"current_aslot"`
	LastIrreversibleBlockNum uint32            `json:"last_irreversible_block_num"`
	TotalVestingShares       string            `json:"total_vesting_shares"`
	TotalVersingFundSteem    string            `json:"total_vesting_fund_steem"`
	HeadBlockId              string            `json:"head_block_id"`
	VirtualSupply            string            `json:"virtual_supply"`
	ConfidentialSupply       string            `json:"confidential_supply"`
	ConfidentialSBDSupply    string            `json:"confidential_sbd_supply"`
	TotalRewardFundSteem     string            `json:"total_reward_fund_steem"`
	TotalActivityFundSteem   string            `json:"total_activity_fund_steem"`
	TotalActivityFundShares  *simpletypes.Int  `json:"total_activity_fund_shares"`
	SBDInterestRate          *simpletypes.Int  `json:"sbd_interest_rate"`
	MaxVirtualBandwidth      *simpletypes.Int  `json:"max_virtual_bandwidth"`
	HeadBlockNumber          *simpletypes.Int  `json:"head_block_number"`
}

type Block struct {
	Number                uint32               `json:"-"`
	Timestamp             *simpletypes.Time    `json:"timestamp"`
	Witness               string               `json:"witness"`
	WitnessSignature      string               `json:"witness_signature"`
	TransactionMerkleRoot string               `json:"transaction_merkle_root"`
	Previous              string               `json:"previous"`
	Extensions            [][]interface{}      `json:"extensions"`
	Transactions          []*types.Transaction `json:"transactions"`
}

type Content struct {
	Id                      *simpletypes.ID   `json:"id"`
	RootTitle               string            `json:"root_title"`
	Active                  *simpletypes.Time `json:"active"`
	AbsRshares              *simpletypes.Int  `json:"abs_rshares"`
	PendingPayoutValue      string            `json:"pending_payout_value"`
	TotalPendingPayoutValue string            `json:"total_pending_payout_value"`
	Category                string            `json:"category"`
	Title                   string            `json:"title"`
	LastUpdate              *simpletypes.Time `json:"last_update"`
	Stats                   string            `json:"stats"`
	Body                    string            `json:"body"`
	Created                 *simpletypes.Time `json:"created"`
	Replies                 []*Content        `json:"replies"`
	Permlink                string            `json:"permlink"`
	JsonMetadata            *ContentMetadata  `json:"json_metadata"`
	Children                *simpletypes.Int  `json:"children"`
	NetRshares              *simpletypes.Int  `json:"net_rshares"`
	URL                     string            `json:"url"`
	ActiveVotes             []*VoteState      `json:"active_votes"`
	ParentPermlink          string            `json:"parent_permlink"`
	CashoutTime             *simpletypes.Time `json:"cashout_time"`
	TotalPayoutValue        string            `json:"total_payout_value"`
	ParentAuthor            string            `json:"parent_author"`
	ChildrenRshares2        *simpletypes.Int  `json:"children_rshares2"`
	Author                  string            `json:"author"`
	Depth                   *simpletypes.Int  `json:"depth"`
	TotalVoteWeight         *simpletypes.Int  `json:"total_vote_weight"`
}

func (content *Content) IsStory() bool {
	return content.ParentAuthor == ""
}

type ContentMetadata struct {
	Flag  bool
	Users []string
	Tags  []string
	Image []string
}

type ContentMetadataRaw struct {
	Users simpletypes.StringSlice `json:"users"`
	Tags  simpletypes.StringSlice `json:"tags"`
	Image simpletypes.StringSlice `json:"image"`
}

func (metadata *ContentMetadata) UnmarshalJSON(data []byte) error {
	unquoted, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	switch unquoted {
	case "true":
		metadata.Flag = true
		return nil
	case "false":
		metadata.Flag = false
		return nil
	}

	if len(unquoted) == 0 {
		var value ContentMetadata
		metadata = &value
		return nil
	}

	var raw ContentMetadataRaw
	if err := json.NewDecoder(strings.NewReader(unquoted)).Decode(&raw); err != nil {
		return err
	}

	metadata.Users = raw.Users
	metadata.Tags = raw.Tags
	metadata.Image = raw.Image

	return nil
}

type VoteState struct {
	Voter   string            `json:"voter"`
	Weight  *simpletypes.Int  `json:"weight"`
	Rshares *simpletypes.Int  `json:"rshares"`
	Percent *simpletypes.Int  `json:"percent"`
	Time    *simpletypes.Time `json:"time"`
}
