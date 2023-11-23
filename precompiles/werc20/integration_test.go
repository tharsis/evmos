package werc20_test

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/evmos/evmos/v15/precompiles/authorization"
	"github.com/evmos/evmos/v15/precompiles/erc20"
	evmosutiltx "github.com/evmos/evmos/v15/testutil/tx"

	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v15/precompiles/testutil"
	"github.com/evmos/evmos/v15/precompiles/werc20"
	"github.com/evmos/evmos/v15/precompiles/werc20/testdata"
	"github.com/evmos/evmos/v15/testutil/integration/evmos/factory"
	"github.com/evmos/evmos/v15/testutil/integration/evmos/keyring"
	erc20types "github.com/evmos/evmos/v15/x/erc20/types"
	evmtypes "github.com/evmos/evmos/v15/x/evm/types"

	//nolint:revive // dot imports are fine for Ginkgo
	. "github.com/onsi/ginkgo/v2"
	//nolint:revive // dot imports are fine for Ginkgo
	. "github.com/onsi/gomega"
)

var _ = Describe("WEVMOS Extension -", func() {
	var (
		WEVMOSContractAddr common.Address
		err                error
		sender             keyring.Key
		amount             *big.Int

		// contractData is a helper struct to hold the addresses and ABIs for the
		// different contract instances that are subject to testing here.
		contractData ContractData

		failCheck testutil.LogCheckArgs
		passCheck testutil.LogCheckArgs
	)

	BeforeEach(func() {
		s.SetupTest()

		sender = s.keyring.GetKey(0)

		WEVMOSContractAddr, err = s.factory.DeployContract(
			sender.Priv,
			evmtypes.EvmTxArgs{}, // NOTE: passing empty struct to use default values
			factory.ContractDeploymentData{
				Contract:        testdata.WEVMOSContract,
				ConstructorArgs: []interface{}{},
			},
		)
		Expect(err).ToNot(HaveOccurred(), "failed to deploy contract")

		// Register the native EVMOS token in the bank module for the network
		evmosMetadata := banktypes.Metadata{
			Description: "The native token of Evmos",
			Base:        s.bondDenom,
			// NOTE: Denom units MUST be increasing
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    s.bondDenom,
					Exponent: 0,
					Aliases:  []string{s.bondDenom},
				},
				{
					Denom:    s.bondDenom,
					Exponent: 18,
				},
			},
			Name:    "Evmos",
			Symbol:  "EVMOS",
			Display: s.bondDenom,
		}
		s.network.App.BankKeeper.SetDenomMetaData(s.network.GetContext(), evmosMetadata)

		tokenPair := erc20types.NewTokenPair(WEVMOSContractAddr, s.bondDenom, erc20types.OWNER_MODULE)

		precompile, err := werc20.NewPrecompile(
			tokenPair,
			s.network.App.BankKeeper,
			s.network.App.AuthzKeeper,
			s.network.App.TransferKeeper,
		)

		Expect(err).ToNot(HaveOccurred(), "failed to create wevmos extension")
		s.precompile = precompile

		err = s.network.App.EvmKeeper.AddEVMExtensions(s.network.GetContext(), precompile)
		Expect(err).ToNot(HaveOccurred(), "failed to add wevmos extension")

		s.tokenDenom = tokenPair.GetDenom()

		contractData = ContractData{
			ownerPriv:      sender.Priv,
			erc20Addr:      WEVMOSContractAddr,
			erc20ABI:       testdata.WEVMOSContract.ABI,
			precompileAddr: s.precompile.Address(),
			precompileABI:  s.precompile.ABI,
		}

		failCheck = testutil.LogCheckArgs{ABIEvents: s.precompile.Events}
		passCheck = failCheck.WithExpPass(true)

		err = s.network.NextBlock()
		Expect(err).ToNot(HaveOccurred(), "failed to advance block")

		// Default sender and amount
		sender = s.keyring.GetKey(0)
		amount = big.NewInt(1e18)
	})

	Context("WEVMOS specific functions", func() {
		When("calling deposit correctly", func() {
			It("should emit the Deposit event but not modify the balance", func() {
				depositCheck := passCheck.WithExpPass(true).WithExpEvents(werc20.EventTypeDeposit)
				txArgs, callArgs := s.getTxAndCallArgs(erc20Call, contractData, werc20.DepositMethod)
				txArgs.Amount = amount

				_, _, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, callArgs, depositCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				s.checkBalances(failCheck, sender, contractData)
			})

			It("should spend the correct minimum gas", func() {
				depositCheck := passCheck.WithExpPass(true).WithExpEvents(werc20.EventTypeDeposit)
				txArgs, callArgs := s.getTxAndCallArgs(erc20Call, contractData, werc20.DepositMethod)
				txArgs.Amount = amount

				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, callArgs, depositCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				Expect(ethRes.GasUsed).To(BeNumerically(">=", werc20.DepositRequiredGas), "expected different gas used")
			})
		})

		When("calling withdraw correctly", func() {
			It("should emit the Withdrawal event but not modify the balance", func() {
				withdrawCheck := passCheck.WithExpPass(true).WithExpEvents(werc20.EventTypeWithdrawal)
				txArgs, callArgs := s.getTxAndCallArgs(erc20Call, contractData, werc20.WithdrawMethod, amount)

				_, _, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, callArgs, withdrawCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				s.checkBalances(failCheck, sender, contractData)
			})

			It("should spend the correct minimum gas", func() {
				withdrawCheck := passCheck.WithExpPass(true).WithExpEvents(werc20.EventTypeWithdrawal)
				txArgs, callArgs := s.getTxAndCallArgs(erc20Call, contractData, werc20.WithdrawMethod, amount)

				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, callArgs, withdrawCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				Expect(ethRes.GasUsed).To(BeNumerically(">=", werc20.WithdrawRequiredGas), "expected different gas used")
			})
		})

		// TODO: How do we actually check the method types here? We can see the correct ones being populated by printing the line in the cmn.Precompile
		When("calling with incomplete data or amount", func() {
			It("calls no call data, with amount - should call `receive` ", func() {
				txArgs, callArgs := s.getTxAndCallArgs(erc20Call, contractData, "")
				txArgs.Amount = amount

				res, err := s.factory.ExecuteContractCall(sender.Priv, txArgs, callArgs)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				depositCheck := passCheck.WithExpPass(true).WithExpEvents(werc20.EventTypeDeposit)
				depositCheck.Res = res
				err = testutil.CheckLogs(depositCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				s.checkBalances(failCheck, sender, contractData)
			})

			It("calls short call data, with amount - should call `fallback` ", func() {
				txArgs, _ := s.getTxAndCallArgs(erc20Call, contractData, "")
				txArgs.Amount = amount
				txArgs.Input = []byte{1, 2, 3} // 3 dummy bytes

				res, err := s.factory.ExecuteEthTx(sender.Priv, txArgs)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				depositCheck := passCheck.WithExpPass(true).WithExpEvents(werc20.EventTypeDeposit)
				depositCheck.Res = res
				err = testutil.CheckLogs(depositCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				s.checkBalances(failCheck, sender, contractData)
			})

			It("calls with non-existing function, with amount - should call `fallback` ", func() {
				txArgs, _ := s.getTxAndCallArgs(erc20Call, contractData, "")
				txArgs.Input = []byte("nonExistingMethod")
				txArgs.Amount = amount

				res, err := s.factory.ExecuteEthTx(sender.Priv, txArgs)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				depositCheck := passCheck.WithExpPass(true).WithExpEvents(werc20.EventTypeDeposit)
				depositCheck.Res = res
				err = testutil.CheckLogs(depositCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				s.checkBalances(failCheck, sender, contractData)
			})

			It("calls non call data, without amount - should call `fallback` ", func() {
				txArgs, _ := s.getTxAndCallArgs(erc20Call, contractData, "")

				res, err := s.factory.ExecuteEthTx(sender.Priv, txArgs)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				depositCheck := passCheck.WithExpPass(true).WithExpEvents(werc20.EventTypeDeposit)
				depositCheck.Res = res
				err = testutil.CheckLogs(depositCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				s.checkBalances(failCheck, sender, contractData)
			})

			It("calls short call data, without amount - should call `fallback` ", func() {
				txArgs, _ := s.getTxAndCallArgs(erc20Call, contractData, "")
				txArgs.Input = []byte{1, 2, 3} // 3 dummy bytes

				res, err := s.factory.ExecuteEthTx(sender.Priv, txArgs)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				depositCheck := passCheck.WithExpPass(true).WithExpEvents(werc20.EventTypeDeposit)
				depositCheck.Res = res
				err = testutil.CheckLogs(depositCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				s.checkBalances(failCheck, sender, contractData)
			})

			It("calls with non-existing function, without amount -  should call `fallback` ", func() {
				txArgs, _ := s.getTxAndCallArgs(erc20Call, contractData, "")
				txArgs.Input = []byte("nonExistingMethod")

				res, err := s.factory.ExecuteEthTx(sender.Priv, txArgs)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				depositCheck := passCheck.WithExpPass(true).WithExpEvents(werc20.EventTypeDeposit)
				depositCheck.Res = res
				err = testutil.CheckLogs(depositCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				s.checkBalances(failCheck, sender, contractData)
			})
		})
	})

	Context("ERC20 specific functions", func() {
		When("querying name", func() {
			It("should return the correct name", func() {
				// Query the name
				txArgs, nameArgs := s.getTxAndCallArgs(directCall, contractData, erc20.NameMethod)

				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, nameArgs, passCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				var name string
				err = s.precompile.UnpackIntoInterface(&name, erc20.NameMethod, ethRes.Ret)
				Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
				Expect(name).To(Equal("Evmos"), "expected different name")
			})
		})

		When("querying symbol", func() {
			It("should return the correct symbol", func() {
				// Query the symbol
				txArgs, symbolArgs := s.getTxAndCallArgs(directCall, contractData, erc20.SymbolMethod)

				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, symbolArgs, passCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				var symbol string
				err = s.precompile.UnpackIntoInterface(&symbol, erc20.SymbolMethod, ethRes.Ret)
				Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
				Expect(symbol).To(Equal("EVMOS"), "expected different symbol")
			})
		})

		When("querying decimals", func() {
			It("should return the correct decimals", func() {
				// Query the decimals
				txArgs, decimalsArgs := s.getTxAndCallArgs(directCall, contractData, erc20.DecimalsMethod)

				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, decimalsArgs, passCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				var decimals uint8
				err = s.precompile.UnpackIntoInterface(&decimals, erc20.DecimalsMethod, ethRes.Ret)
				Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
				Expect(decimals).To(Equal(uint8(18)), "expected different decimals")
			})
		})

		When("querying balance", func() {
			It("should return an existing balance", func() {
				// Query the balance
				txArgs, balancesArgs := s.getTxAndCallArgs(directCall, contractData, erc20.BalanceOfMethod, sender.Addr)

				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, balancesArgs, passCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				expBalance := s.network.App.BankKeeper.GetBalance(s.network.GetContext(), sender.AccAddr, s.bondDenom)

				var balance *big.Int
				err = s.precompile.UnpackIntoInterface(&balance, erc20.BalanceOfMethod, ethRes.Ret)
				Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
				Expect(balance).To(Equal(expBalance.Amount.BigInt()), "expected different balance")
			})

			It("should return a 0 balance new address", func() {
				// Query the balance
				txArgs, balancesArgs := s.getTxAndCallArgs(directCall, contractData, erc20.BalanceOfMethod, evmosutiltx.GenerateAddress())

				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, balancesArgs, passCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				var balance *big.Int
				err = s.precompile.UnpackIntoInterface(&balance, erc20.BalanceOfMethod, ethRes.Ret)
				Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
				Expect(balance.Int64()).To(Equal(int64(0)), "expected different balance")
			})
		})

		When("querying allowance", func() {
			It("should return an existing allowance", func() {
				grantee := evmosutiltx.GenerateAddress()
				granter := sender
				authzCoins := sdk.Coins{sdk.NewInt64Coin(s.tokenDenom, 100)}

				s.setupSendAuthzForContract(directCall, grantee, granter.Priv, authzCoins)

				txArgs, allowanceArgs := s.getTxAndCallArgs(directCall, contractData, auth.AllowanceMethod, granter.Addr, grantee)

				_, ethRes, err := s.factory.CallContractAndCheckLogs(granter.Priv, txArgs, allowanceArgs, passCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				var allowance *big.Int
				err = s.precompile.UnpackIntoInterface(&allowance, auth.AllowanceMethod, ethRes.Ret)
				Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
				Expect(allowance).To(Equal(authzCoins[0].Amount.BigInt()), "expected different allowance")
			})

			It("should return zero if no balance exists", func() {
				address := evmosutiltx.GenerateAddress()

				// Query the balance
				txArgs, balancesArgs := s.getTxAndCallArgs(directCall, contractData, erc20.BalanceOfMethod, address)

				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, balancesArgs, passCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				var balance *big.Int
				err = s.precompile.UnpackIntoInterface(&balance, erc20.BalanceOfMethod, ethRes.Ret)
				Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
				Expect(balance.Int64()).To(BeZero(), "expected zero balance")
			})
		})

		When("querying total supply", func() {
			It("should return the total supply", func() {
				expSupply, ok := new(big.Int).SetString("11000000000000000000", 10)
				Expect(ok).To(BeTrue(), "failed to parse expected supply")

				// Query the balance
				txArgs, supplyArgs := s.getTxAndCallArgs(directCall, contractData, erc20.TotalSupplyMethod)

				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, supplyArgs, passCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				var supply *big.Int
				err = s.precompile.UnpackIntoInterface(&supply, erc20.TotalSupplyMethod, ethRes.Ret)
				Expect(err).ToNot(HaveOccurred(), "failed to unpack result")
				Expect(supply).To(Equal(expSupply), "expected different supply")
			})
		})

		When("transferring tokens", func() {
			It("it should transfer tokens to a receiver using `transfer`", func() {
				// Get receiver address
				receiver := s.keyring.GetKey(1)

				senderBalance := s.network.App.BankKeeper.GetAllBalances(s.network.GetContext(), sender.AccAddr)
				receiverBalance := s.network.App.BankKeeper.GetAllBalances(s.network.GetContext(), receiver.AccAddr)

				// Transfer tokens
				txArgs, transferArgs := s.getTxAndCallArgs(directCall, contractData, erc20.TransferMethod, receiver.Addr, amount)
				txArgs.GasPrice = big.NewInt(765625000)
				transferCoins := sdk.Coins{sdk.NewInt64Coin(s.tokenDenom, amount.Int64())}

				transferCheck := passCheck.WithExpEvents(erc20.EventTypeTransfer)
				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, transferArgs, transferCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				gasAmount := ethRes.GasUsed * txArgs.GasPrice.Uint64()
				coinsWithGasIncluded := transferCoins.Add(sdk.NewInt64Coin(s.bondDenom, int64(gasAmount)))
				s.ExpectBalances(
					[]ExpectedBalance{
						{address: sender.AccAddr, expCoins: senderBalance.Sub(coinsWithGasIncluded...)},
						{address: receiver.AccAddr, expCoins: receiverBalance.Add(transferCoins...)},
					},
				)
			})

			It("it should transfer tokens to a receiver using `transferFrom`", func() {
				// Get receiver address
				receiver := s.keyring.GetKey(1)

				senderBalance := s.network.App.BankKeeper.GetAllBalances(s.network.GetContext(), sender.AccAddr)
				receiverBalance := s.network.App.BankKeeper.GetAllBalances(s.network.GetContext(), receiver.AccAddr)

				// Transfer tokens
				txArgs, transferArgs := s.getTxAndCallArgs(directCall, contractData, erc20.TransferFromMethod, sender.Addr, receiver.Addr, amount)
				txArgs.GasPrice = big.NewInt(765625000)
				transferCoins := sdk.Coins{sdk.NewInt64Coin(s.tokenDenom, amount.Int64())}

				transferCheck := passCheck.WithExpEvents(erc20.EventTypeTransfer)
				_, ethRes, err := s.factory.CallContractAndCheckLogs(sender.Priv, txArgs, transferArgs, transferCheck)
				Expect(err).ToNot(HaveOccurred(), "unexpected result calling contract")

				gasAmount := ethRes.GasUsed * txArgs.GasPrice.Uint64()
				coinsWithGasIncluded := transferCoins.Add(sdk.NewInt64Coin(s.bondDenom, int64(gasAmount)))
				s.ExpectBalances(
					[]ExpectedBalance{
						{address: sender.AccAddr, expCoins: senderBalance.Sub(coinsWithGasIncluded...)},
						{address: receiver.AccAddr, expCoins: receiverBalance.Add(transferCoins...)},
					},
				)
			})
		})
	})
})
