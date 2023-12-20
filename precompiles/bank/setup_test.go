package bank_test

import (
	"testing"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v16/precompiles/bank"
	"github.com/evmos/evmos/v16/testutil/integration/evmos/factory"
	"github.com/evmos/evmos/v16/testutil/integration/evmos/grpc"
	testkeyring "github.com/evmos/evmos/v16/testutil/integration/evmos/keyring"
	"github.com/evmos/evmos/v16/testutil/integration/evmos/network"
	"github.com/stretchr/testify/suite"
)

var s *PrecompileTestSuite

// PrecompileTestSuite is the implementation of the TestSuite interface for ERC20 precompile
// unit tests.
type PrecompileTestSuite struct {
	suite.Suite

	bondDenom, tokenDenom string
	evmosAddr, xmplAddr   common.Address

	// tokenDenom is the specific token denomination used in testing the ERC20 precompile.
	// This denomination is used to instantiate the precompile.
	network     *network.UnitTestNetwork
	factory     factory.TxFactory
	grpcHandler grpc.Handler
	keyring     testkeyring.Keyring

	precompile *bank.Precompile
}

func TestPrecompileTestSuite(t *testing.T) {
	s = new(PrecompileTestSuite)
	suite.Run(t, s)
}

func (s *PrecompileTestSuite) SetupTest() {
	// Mint and register a second coin for testing purposes
	// FIXME the RegisterCoin logic will need to be refactored
	// once logic is integrated
	// with the protocol via genesis and/or a transaction
	s.tokenDenom = "xmpl"

	keyring := testkeyring.New(2)
	integrationNetwork := network.NewUnitTestNetwork(
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
		network.WithOtherDenoms([]string{is.tokenDenom}),
	)
	grpcHandler := grpc.NewIntegrationHandler(integrationNetwork)
	txFactory := factory.New(integrationNetwork, grpcHandler)

	ctx := integrationNetwork.GetContext()
	sk := integrationNetwork.App.StakingKeeper
	bondDenom, err := sk.BondDenom(ctx)
	s.Require().NoError(err, "failed to get bond denom")
	s.Require().NotEmpty(bondDenom, "bond denom cannot be empty")

	s.bondDenom = bondDenom
	s.factory = txFactory
	s.grpcHandler = grpcHandler
	s.keyring = keyring
	s.network = integrationNetwork

	// Register EVMOS
	evmosMetadata, found := s.network.App.BankKeeper.GetDenomMetaData(s.network.GetContext(), s.bondDenom)
	s.Require().True(found, "expected evmos denom metadata")

	// FIXME need to refactor this once the RegisterCoin logic is integrated
	// with the protocol via genesis and/or a transaction
	tokenPair, err := s.network.App.Erc20Keeper.RegisterCoin(s.network.GetContext(), evmosMetadata)
	s.Require().NoError(err, "failed to register coin")

	s.evmosAddr = common.HexToAddress(tokenPair.Erc20Address)

	xmplMetadata := banktypes.Metadata{
		Description: "An exemplary token",
		Base:        s.tokenDenom,
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    s.tokenDenom,
				Exponent: 0,
				Aliases:  []string{s.tokenDenom},
			},
			{
				Denom:    s.tokenDenom,
				Exponent: 18,
			},
		},
		Name:    "Exemplary",
		Symbol:  "XMPL",
		Display: s.tokenDenom,
	}

	// FIXME need to refactor this once the RegisterCoin logic is integrated
	// with the protocol via genesis and/or a transaction
	tokenPair, err = s.network.App.Erc20Keeper.RegisterCoin(s.network.GetContext(), xmplMetadata)
	s.Require().NoError(err, "failed to register coin")

	s.xmplAddr = common.HexToAddress(tokenPair.Erc20Address)

	s.precompile = s.setupBankPrecompile()
}
