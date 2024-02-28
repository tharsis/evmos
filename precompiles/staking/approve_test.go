package staking_test

import (
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmosutiltx "github.com/evmos/evmos/v16/testutil/tx"
	"github.com/evmos/evmos/v16/x/evm/statedb"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkauthz "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/evmos/v16/precompiles/authorization"
	cmn "github.com/evmos/evmos/v16/precompiles/common"
	"github.com/evmos/evmos/v16/precompiles/staking"
	"github.com/evmos/evmos/v16/precompiles/testutil"
)

func (s *PrecompileTestSuite) TestApprove() {
	var (
		ctx  sdk.Context
		stDB *statedb.StateDB
	)
	method := s.precompile.Methods[authorization.ApproveMethod]

	testCases := []struct {
		name        string
		malleate    func(*vm.Contract) []interface{}
		postCheck   func(data []byte, inputArgs []interface{})
		gas         uint64
		expError    bool
		errContains string
	}{
		{
			"fail - empty input args",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{}
			},
			func(data []byte, inputArgs []interface{}) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 3, 0),
		},
		{
			"fail - invalid message type",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					abi.MaxUint256,
					[]string{"invalid"},
				}
			},
			func(data []byte, inputArgs []interface{}) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidMsgType, "staking", "invalid"),
		},
		// TODO: enable this test once we check if spender and origin are the same
		// {
		//	"fail - origin address is the same the spender address",
		//	func(_ *vm.Contract) []interface{} {
		//		return []interface{}{
		//			s.keyring.GetAddr(0),
		//			abi.MaxUint256,
		//			[]string{"invalid"},
		//		}
		//	},
		//	func(data []byte, inputArgs []interface{}) {},
		//	200000,
		//	true,
		//	"is the same as spender",
		// },
		{
			"success - MsgDelegate with unlimited coins",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					abi.MaxUint256,
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				s.Require().Equal(data, cmn.TrueValue)
				authz, expirationTime := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.DelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))

				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)
				s.Require().Equal(authz.AuthorizationType, staking.DelegateAuthz)
				var coin *sdk.Coin
				s.Require().Equal(authz.MaxTokens, coin)
			},
			20000,
			false,
			"",
		},
		{
			"success - MsgUndelegate with unlimited coins",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					abi.MaxUint256,
					[]string{staking.UndelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				s.Require().Equal(data, cmn.TrueValue)

				authz, expirationTime := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.UndelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)
				s.Require().Equal(authz.AuthorizationType, staking.UndelegateAuthz)
				var coin *sdk.Coin
				s.Require().Equal(authz.MaxTokens, coin)
			},
			20000,
			false,
			"",
		},
		{
			"success - MsgRedelegate with unlimited coins",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					abi.MaxUint256,
					[]string{staking.RedelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				s.Require().Equal(data, cmn.TrueValue)

				authz, expirationTime := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.RedelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)
				s.Require().Equal(authz.AuthorizationType, staking.RedelegateAuthz)
				var coin *sdk.Coin
				s.Require().Equal(authz.MaxTokens, coin)
			},
			20000,
			false,
			"",
		},
		{
			"success - All staking methods with certain amount of coins",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(1e18),
					[]string{
						staking.DelegateMsg,
						staking.UndelegateMsg,
						staking.RedelegateMsg,
					},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				s.Require().Equal(data, cmn.TrueValue)

				allAuthz, err := s.network.App.AuthzKeeper.GetAuthorizations(ctx, s.keyring.GetAddr(0).Bytes(), s.keyring.GetAddr(0).Bytes())
				s.Require().NoError(err)
				s.Require().Len(allAuthz, 3)
			},
			20000,
			false,
			"",
		},
		{
			"success - remove MsgDelegate authorization",
			func(contract *vm.Contract) []interface{} {
				res, err := s.precompile.Approve(ctx, s.keyring.GetAddr(0), stDB, &method, []interface{}{s.keyring.GetAddr(0), big.NewInt(1), []string{staking.DelegateMsg}})
				s.Require().NoError(err)
				s.Require().Equal(res, cmn.TrueValue)

				authz, expirationTime := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.DelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)

				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(0),
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				s.Require().Equal(data, cmn.TrueValue)

				authz, expirationTime := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.DelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().Nil(authz)
				s.Require().Nil(expirationTime)
			},
			200000,
			false,
			"",
		},
		{
			"success - MsgDelegate with 1 Evmos as limit amount",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(1e18),
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				s.Require().Equal(data, cmn.TrueValue)

				authz, expirationTime := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.DelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)
				s.Require().Equal(authz.AuthorizationType, staking.DelegateAuthz)
				s.Require().Equal(authz.MaxTokens, &sdk.Coin{Denom: s.bondDenom, Amount: math.NewInt(1e18)})
			},
			20000,
			false,
			"",
		},
		{
			"success - Authorization should only be created for validators that are not jailed",
			func(_ *vm.Contract) []interface{} {
				var err error
				// Jail a validator
				valAddr, err = sdk.ValAddressFromBech32(s.network.GetValidators()[0].GetOperator())
				s.Require().NoError(err)
				err = s.network.App.StakingKeeper.Jail(ctx, sdk.ConsAddress(valAddr))
				s.Require().NoError(err, "failed to jail a validator")

				// When a delegator redelegates/undelegates from a validator, the validator
				// switches to Unbonding status.
				// Thus, validators with this status should be considered for the authorization

				// Unbond another validator
				valAddr1, err := sdk.ValAddressFromBech32(s.network.GetValidators()[1].GetOperator())
				s.Require().NoError(err)
				amount, err := s.network.App.StakingKeeper.Unbond(ctx, s.keyring.GetAddr(0).Bytes(), valAddr1, math.LegacyOneDec())
				s.Require().NoError(err, "expected no error unbonding validator")
				s.Require().Equal(math.NewInt(1e18), amount, "expected different amount of tokens to be unbonded")

				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(1e18),
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				s.Require().Equal(data, cmn.TrueValue)

				authz, expirationTime := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.DelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)
				s.Require().Equal(authz.AuthorizationType, staking.DelegateAuthz)
				s.Require().Equal(authz.MaxTokens, &sdk.Coin{Denom: s.bondDenom, Amount: math.NewInt(1e18)})
				// Check that the bonded and unbonding validators are included on the authorization
				s.Require().Len(authz.GetAllowList().Address, 2, "should only be two validators in the allow list")
			},
			1e6,
			false,
			"",
		},
		{
			"success - MsgUndelegate with 1 Evmos as limit amount",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(1e18),
					[]string{staking.UndelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				s.Require().Equal(data, cmn.TrueValue)

				authz, expirationTime := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.UndelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)
				s.Require().Equal(authz.AuthorizationType, staking.UndelegateAuthz)
				s.Require().Equal(authz.MaxTokens, &sdk.Coin{Denom: s.bondDenom, Amount: math.NewInt(1e18)})
			},
			20000,
			false,
			"",
		},
		{
			"success - MsgRedelegate with 1 Evmos as limit amount",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(1e18),
					[]string{staking.RedelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				s.Require().Equal(data, cmn.TrueValue)

				authz, expirationTime := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.RedelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)
			},
			20000,
			false,
			"",
		},
		{
			"success - MsgRedelegate, MsgUndelegate and MsgDelegate with 1 Evmos as limit amount",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(1e18),
					[]string{
						staking.RedelegateMsg,
						staking.UndelegateMsg,
						staking.DelegateMsg,
					},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				s.Require().Equal(data, cmn.TrueValue)

				authz, expirationTime := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.DelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)
				s.Require().Equal(authz.AuthorizationType, staking.DelegateAuthz)
				s.Require().Equal(authz.MaxTokens, &sdk.Coin{Denom: s.bondDenom, Amount: math.NewInt(1e18)})

				authz, expirationTime = CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.UndelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)

				s.Require().Equal(authz.AuthorizationType, staking.UndelegateAuthz)
				s.Require().Equal(authz.MaxTokens, &sdk.Coin{Denom: s.bondDenom, Amount: math.NewInt(1e18)})

				authz, expirationTime = CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.RedelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().NotNil(expirationTime)

				s.Require().Equal(authz.AuthorizationType, staking.RedelegateAuthz)
				s.Require().Equal(authz.MaxTokens, &sdk.Coin{Denom: s.bondDenom, Amount: math.NewInt(1e18)})

				// TODO: Bug here it returns 3 REDELEGATE authorizations
				allAuthz, err := s.network.App.AuthzKeeper.GetAuthorizations(s.network.GetContext(), s.keyring.GetAddr(0).Bytes(), s.keyring.GetAddr(0).Bytes())
				s.Require().NoError(err)
				s.Require().Len(allAuthz, 3)
			},
			20000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			ctx = s.network.GetContext()
			stDB = s.network.GetStateDB()

			var contract *vm.Contract
			contract, ctx = testutil.NewPrecompileContract(s.T(), ctx, s.keyring.GetAddr(0), s.precompile, tc.gas)

			args := tc.malleate(contract)
			bz, err := s.precompile.Approve(ctx, s.keyring.GetAddr(0), stDB, &method, args)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
				s.Require().Empty(bz)
			} else {
				s.Require().NoError(err)
				s.Require().NotEmpty(bz)
				tc.postCheck(bz, args)
			}
		})
	}
}

func (s *PrecompileTestSuite) TestDecreaseAllowance() {
	var (
		ctx  sdk.Context
		stDB *statedb.StateDB
	)
	method := s.precompile.Methods[authorization.DecreaseAllowanceMethod]

	testCases := []struct {
		name        string
		malleate    func(_ *vm.Contract) []interface{}
		postCheck   func(data []byte, inputArgs []interface{})
		gas         uint64
		expError    bool
		errContains string
	}{
		{
			"fail - empty input args",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{}
			},
			func(data []byte, inputArgs []interface{}) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 3, 0),
		},
		//// TODO: enable this once we check origin is not the spender
		// {
		//	"fail - origin address is the spender address",
		//	func(_ *vm.Contract) []interface{} {
		//		return []interface{}{
		//			s.keyring.GetAddr(0),
		//			abi.MaxUint256,
		//			[]string{staking.DelegateMsg},
		//		}
		//	},
		//	func(data []byte, inputArgs []interface{}) {},
		//	200000,
		//	true,
		//	"is the same as spender",
		// },
		{
			"fail - delegate authorization does not exists",
			func(_ *vm.Contract) []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(15000),
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
			},
			200000,
			true,
			"authorization to /cosmos.staking.v1beta1.MsgDelegate",
		},
		{
			"fail - delegate authorization is a generic Authorization",
			func(_ *vm.Contract) []interface{} {
				authz := sdkauthz.NewGenericAuthorization(staking.DelegateMsg)
				exp := time.Now().Add(time.Hour)
				err := s.network.App.AuthzKeeper.SaveGrant(ctx, s.keyring.GetAddr(0).Bytes(), s.keyring.GetAddr(0).Bytes(), authz, &exp)
				s.Require().NoError(err)
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(15000),
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
			},
			200000,
			true,
			sdkauthz.ErrUnknownAuthorizationType.Error(),
		},
		{
			"fail - decrease allowance amount is greater than the authorization limit",
			func(contract *vm.Contract) []interface{} {
				approveArgs := []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(1e18),
					[]string{staking.DelegateMsg},
				}
				resp, err := s.precompile.Approve(ctx, s.keyring.GetAddr(0), stDB, &method, approveArgs)
				s.Require().NoError(err)
				s.Require().Equal(resp, cmn.TrueValue)

				authz, _ := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.DelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().Equal(authz.AuthorizationType, staking.DelegateAuthz)
				s.Require().Equal(authz.MaxTokens, &sdk.Coin{Denom: s.bondDenom, Amount: math.NewInt(1e18)})

				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(2e18),
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {},
			200000,
			true,
			"amount by which the allowance should be decreased is greater than the authorization limit",
		},
		{
			"success - decrease delegate authorization allowance by 1 Evmos",
			func(_ *vm.Contract) []interface{} {
				s.ApproveAndCheckAuthz(method, staking.DelegateMsg, big.NewInt(2e18))
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(1e18),
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				authz, _ := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.DelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().Equal(authz.AuthorizationType, staking.DelegateAuthz)
				s.Require().Equal(authz.MaxTokens, &sdk.Coin{Denom: s.bondDenom, Amount: math.NewInt(1e18)})
			},
			200000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest() // reset
			ctx = s.network.GetContext()
			stDB = s.network.GetStateDB()

			var contract *vm.Contract
			contract, ctx = testutil.NewPrecompileContract(s.T(), ctx, s.keyring.GetAddr(0), s.precompile, tc.gas)

			args := tc.malleate(contract)
			bz, err := s.precompile.DecreaseAllowance(ctx, s.keyring.GetAddr(0), stDB, &method, args)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
				s.Require().Empty(bz)
			} else {
				s.Require().NoError(err)
				s.Require().NotEmpty(bz)
				tc.postCheck(bz, args)
			}
		})
	}
}

func (s *PrecompileTestSuite) TestIncreaseAllowance() {
	var (
		ctx  sdk.Context
		stDB *statedb.StateDB
	)
	method := s.precompile.Methods[authorization.IncreaseAllowanceMethod]

	testCases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func(data []byte, inputArgs []interface{})
		gas         uint64
		expError    bool
		errContains string
	}{
		{
			"fail - empty input args",
			func() []interface{} {
				return []interface{}{}
			},
			func(data []byte, inputArgs []interface{}) {},
			200000,
			true,
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 3, 0),
		},
		// TODO: enable this once we check origin is not the same as spender
		// {
		//	"fail - origin address is the spender address",
		//	func(_ *vm.Contract) []interface{} {
		//		return []interface{}{
		//			s.keyring.GetAddr(0),
		//			abi.MaxUint256,
		//			[]string{staking.DelegateMsg},
		//		}
		//	},
		//	func(data []byte, inputArgs []interface{}) {},
		//	200000,
		//	true,
		//	"is the same as spender",
		// },
		{
			"fail - delegate authorization does not exists",
			func() []interface{} {
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(15000),
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
			},
			200000,
			true,
			"authorization to /cosmos.staking.v1beta1.MsgDelegate",
		},
		{
			"success - no-op, allowance amount is already set to the maximum value",
			func() []interface{} {
				approveArgs := []interface{}{
					s.keyring.GetAddr(0),
					abi.MaxUint256,
					[]string{staking.DelegateMsg},
				}
				resp, err := s.precompile.Approve(ctx, s.keyring.GetAddr(0), stDB, &method, approveArgs)
				s.Require().NoError(err)
				s.Require().Equal(resp, cmn.TrueValue)

				authz, _ := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.DelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().Equal(authz.AuthorizationType, staking.DelegateAuthz)
				var coin *sdk.Coin
				s.Require().Equal(authz.MaxTokens, coin)

				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(2e18),
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {},
			200000,
			false,
			"",
		},
		{
			"success - increase delegate authorization allowance by 1 Evmos",
			func() []interface{} {
				s.ApproveAndCheckAuthz(method, staking.DelegateMsg, big.NewInt(1e18))
				return []interface{}{
					s.keyring.GetAddr(0),
					big.NewInt(1e18),
					[]string{staking.DelegateMsg},
				}
			},
			func(data []byte, inputArgs []interface{}) {
				authz, _ := CheckAuthorization(ctx, s.network.App.AuthzKeeper, staking.DelegateAuthz, s.keyring.GetAddr(0), s.keyring.GetAddr(0))
				s.Require().NotNil(authz)
				s.Require().Equal(authz.AuthorizationType, staking.DelegateAuthz)
				s.Require().Equal(authz.MaxTokens, &sdk.Coin{Denom: s.bondDenom, Amount: math.NewInt(2e18)})
			},
			200000,
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest() // reset
			ctx = s.network.GetContext()
			stDB = s.network.GetStateDB()

			args := tc.malleate()
			bz, err := s.precompile.IncreaseAllowance(ctx, s.keyring.GetAddr(0), stDB, &method, args)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
				s.Require().Empty(bz)
			} else {
				s.Require().NoError(err)
				s.Require().NotEmpty(bz)
				tc.postCheck(bz, args)
			}
		})
	}
}

func (s *PrecompileTestSuite) TestRevoke() {
	var ctx sdk.Context

	method := s.precompile.Methods[authorization.RevokeMethod]
	granteeAddr := evmosutiltx.GenerateAddress()
	granterAddr := s.keyring.GetAddr(0)
	createdAuthz := staking.DelegateAuthz
	approvedCoin := &sdk.Coin{Denom: s.bondDenom, Amount: math.NewInt(1e18)}

	testCases := []struct {
		name        string
		malleate    func() []interface{}
		postCheck   func(data []byte, inputArgs []interface{})
		expError    bool
		errContains string
	}{
		{
			name: "fail - empty input args",
			malleate: func() []interface{} {
				return []interface{}{}
			},
			expError:    true,
			errContains: fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 2, 0),
		},
		{
			name: "fail - authorization does not exist",
			malleate: func() []interface{} {
				return []interface{}{
					granteeAddr,
					[]string{staking.UndelegateMsg},
				}
			},
			postCheck: func(data []byte, inputArgs []interface{}) {
				// expect authorization to still be there
				authz, _ := CheckAuthorization(ctx, s.network.App.AuthzKeeper, createdAuthz, granteeAddr, granterAddr)
				s.Require().NotNil(authz)
			},
			expError:    true,
			errContains: "authorization not found",
		},
		{
			name: "pass - authorization revoked",
			malleate: func() []interface{} {
				return []interface{}{
					granteeAddr,
					[]string{staking.DelegateMsg},
				}
			},
			postCheck: func(data []byte, inputArgs []interface{}) {
				// expect authorization to be removed
				authz, _ := CheckAuthorization(ctx, s.network.App.AuthzKeeper, createdAuthz, granteeAddr, granterAddr)
				s.Require().Nil(authz, "expected authorization to be removed")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			s.SetupTest() // reset
			ctx = s.network.GetContext()

			// Create a delegate authorization
			validators, err := s.network.App.StakingKeeper.GetLastValidators(ctx)
			s.Require().NoError(err)
			valAddrs := make([]sdk.ValAddress, len(validators))
			for i, val := range validators {
				valAddrs[i] = sdk.ValAddress(val.GetOperator())
			}
			delegationAuthz, err := stakingtypes.NewStakeAuthorization(
				valAddrs,
				nil,
				createdAuthz,
				approvedCoin,
			)
			s.Require().NoError(err)

			expiration := ctx.BlockTime().Add(time.Hour * 24 * 365).UTC()
			err = s.network.App.AuthzKeeper.SaveGrant(ctx, granteeAddr.Bytes(), granterAddr.Bytes(), delegationAuthz, &expiration)
			s.Require().NoError(err, "failed to save authorization")
			authz, _ := CheckAuthorization(ctx, s.network.App.AuthzKeeper, createdAuthz, granteeAddr, granterAddr)
			s.Require().NotNil(authz, "expected authorization to be set")

			args := tc.malleate()
			bz, err := s.precompile.Revoke(ctx, granterAddr, s.network.GetStateDB(), &method, args)

			if tc.expError {
				s.Require().ErrorContains(err, tc.errContains)
				s.Require().Empty(bz)
			} else {
				s.Require().NoError(err)
				s.Require().NotEmpty(bz)
				tc.postCheck(bz, args)
			}
		})
	}
}
