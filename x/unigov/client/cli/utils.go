package cli

import (
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"github.com/Canto-Network/canto/v3/x/unigov/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

// PARSING METADATA ACCORDING TO PROPOSAL STRUCT IN GOVTYPES TYPE IN UNIGOV

// ParseRegisterCoinProposal reads and parses a ParseRegisterCoinProposal from a file.
func ParseLendingMarketMetadata(cdc codec.JSONCodec, metadataFile string) (types.LendingMarketMetadata, error) {
	propMetaData := types.LendingMarketMetadata{}

	contents, err := ioutil.ReadFile(filepath.Clean(metadataFile))
	if err != nil {
		return propMetaData, err
	}

	// if err = cdc.UnmarshalJSON(contents, &propMetaData); err != nil {
	// 	return propMetaData, err
	// }

	if err = json.Unmarshal(contents, &propMetaData); err != nil {
		return types.LendingMarketMetadata{}, err 
	}
	
	return propMetaData, nil
}

func ParseTreasuryMetadata(cdc codec.JSONCodec, metadataFile string) (types.TreasuryProposalMetadata, error) {
	propMetaData := types.TreasuryProposalMetadata{}

	contents, err := ioutil.ReadFile(filepath.Clean(metadataFile))
	if err != nil {
		return propMetaData, err
	}

	// if err = cdc.UnmarshalJSON(contents, &propMetaData); err != nil {
	// 	return propMetaData, err
	// }

	if err = json.Unmarshal(contents, &propMetaData); err != nil {
		return types.TreasuryProposalMetadata{}, err 
	}
	
	return propMetaData, nil
}
