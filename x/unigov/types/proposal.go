package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeLendingMarket string = "Lending-Market"
	MaxDescriptionLength      int    = 1000
	MaxTitleLength            int    = 140
)

var (
	_ govtypes.Content = &LendingMarketProposal{}
)

//Register Compound Proposal type as a valid proposal type in goveranance module
func init() {
	govtypes.RegisterProposalType(ProposalTypeLendingMarket)
	govtypes.RegisterProposalTypeCodec(&LendingMarketProposal{}, "unigov/LendingMarketProposal")
}

func NewLendingMarketProposal(title, description string, accounts [][]byte, propId uint64,
	values []uint64, calldatas [][]byte,
	signatures []string) govtypes.Content {
	return &LendingMarketProposal{
		Title:       title,
		Description: description,
		Account:     accounts,
		PropId:      propId,
		Values:      values,
		Calldatas:   calldatas,
		Signatures:  signatures,
	}
}

func (*LendingMarketProposal) ProposalRoute() string { return RouterKey }

func (*LendingMarketProposal) ProposalType() string {
	return ProposalTypeLendingMarket
}

func (lm *LendingMarketProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(lm); err != nil {
		return err
	}

	cd, vals, sigs := len(lm.Calldatas), len(lm.Values), len(lm.Signatures)

	if cd != vals {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal array arguments must be same length")
	}

	if vals != sigs {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal array arguments must be same length")
	}
	return nil
}
