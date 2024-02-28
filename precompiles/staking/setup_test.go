package staking_test

import (
	"testing"

	"github.com/evmos/evmos/v16/precompiles/staking"
	"github.com/evmos/evmos/v16/testutil/integration/evmos/factory"
	"github.com/evmos/evmos/v16/testutil/integration/evmos/grpc"
	testkeyring "github.com/evmos/evmos/v16/testutil/integration/evmos/keyring"
	"github.com/evmos/evmos/v16/testutil/integration/evmos/network"

	// //nolint:revive // dot imports are fine for Ginkgo
	// . "github.com/onsi/ginkgo/v2"
	// //nolint:revive // dot imports are fine for Ginkgo
	// . "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"
)

type PrecompileTestSuite struct {
	suite.Suite

	network     *network.UnitTestNetwork
	factory     factory.TxFactory
	grpcHandler grpc.Handler
	keyring     testkeyring.Keyring

	bondDenom      string
	precompile     *staking.Precompile
}

func TestPrecompileTestSuite(t *testing.T) {
	suite.Run(t, new(PrecompileTestSuite))

	// Run Ginkgo integration tests
	// TODO uncomment
	// RegisterFailHandler(Fail)
	// RunSpecs(t, "Precompile Test Suite")
}

func (s *PrecompileTestSuite) SetupTest() {
	keyring := testkeyring.New(2)
	unitNetwork := network.NewUnitTestNetwork(
		network.WithPreFundedAccounts(keyring.GetAllAccAddrs()...),
	)
	grpcHandler := grpc.NewIntegrationHandler(unitNetwork)
	txFactory := factory.New(unitNetwork, grpcHandler)

	ctx := unitNetwork.GetContext()
	sk := unitNetwork.App.StakingKeeper
	bondDenom, err := sk.BondDenom(ctx)
	s.Require().NoError(err, "failed to get bond denom")
	s.Require().NotEmpty(bondDenom, "bond denom cannot be empty")

	s.bondDenom = bondDenom
	s.factory = txFactory
	s.grpcHandler = grpcHandler
	s.keyring = keyring
	s.network = unitNetwork

	s.precompile = s.setupStakingPrecompile()
}
