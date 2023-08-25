package grpc

import (
	"context"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// GetBalance returns the balance for the given address.
// It uses the network's config's denom.
func (gqh *GrpcQueryHelper) GetBalance(address sdktypes.AccAddress, denom string) (*banktypes.QueryBalanceResponse, error) {
	bankClient := gqh.getBankClient()
	return bankClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: address.String(),
		Denom:   denom,
	})
}
