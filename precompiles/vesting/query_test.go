// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package vesting_test

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v13/precompiles/common"
	"github.com/evmos/evmos/v13/precompiles/vesting"
)

func (s *PrecompileTestSuite) TestBalances() {
	method := s.precompile.Methods[vesting.BalancesMethod]

	testCases := []struct {
		name        string
		malleate    func() []interface{}
		gas         uint64
		postCheck   func(data []byte)
		expError    bool
		errContains string
	}{
		{
			"fail - empty input args",
			func() []interface{} {
				return []interface{}{}
			},
			200000,
			func(data []byte) {},
			true,
			fmt.Sprintf(cmn.ErrInvalidNumberOfArgs, 1, 0),
		},
		{
			"fail - invalid address",
			func() []interface{} {
				return []interface{}{
					"12asji1",
				}
			},
			200000,
			func(data []byte) {},
			true,
			"invalid type for vestingAddress",
		},
		{
			"fail - account is not a vesting account",
			func() []interface{} {
				return []interface{}{
					s.address,
				}
			},
			200000,
			func(data []byte) {},
			true,
			"is not a vesting account",
		},
		{
			"success - should return vesting account balances",
			func() []interface{} {
				s.CreateTestClawbackVestingAccount()
				s.FundTestClawbackVestingAccount()
				return []interface{}{
					toAddr,
				}
			},
			200000,
			func(data []byte) {
				var out vesting.BalancesOutput
				err := s.precompile.UnpackIntoInterface(&out, vesting.BalancesMethod, data)
				s.Require().NoError(err)
				s.Require().Equal(out.Locked, lockupPeriods[0].Amount)
				s.Require().Equal(out.Unvested, lockupPeriods[0].Amount)
			},
			false,
			"",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest() // reset
			contract := vm.NewContract(vm.AccountRef(s.address), s.precompile, big.NewInt(0), tc.gas)

			bz, err := s.precompile.Balances(s.ctx, contract, &method, tc.malleate())

			if tc.expError {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			} else {
				s.Require().NoError(err)
				s.Require().NotEmpty(bz)
				tc.postCheck(bz)
			}
		})
	}
}
