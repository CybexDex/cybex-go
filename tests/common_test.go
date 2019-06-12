package tests

import (
	"testing"
	"time"

	"github.com/CybexDex/cybex-go/api"
	"github.com/CybexDex/cybex-go/config"
	"github.com/CybexDex/cybex-go/types"
	"github.com/stretchr/testify/suite"

	//import operations to initialize types.OperationMap
	_ "github.com/CybexDex/cybex-go/operations"
)

type commonTest struct {
	suite.Suite
	TestAPI api.CybexAPI
}

func (suite *commonTest) SetupTest() {
	suite.TestAPI = NewTestAPI(
		suite.T(),
		WsFullApiUrl,
		RpcFullApiUrl,
	)
}

func (suite *commonTest) TearDownTest() {
	if err := suite.TestAPI.Close(); err != nil {
		suite.FailNow(err.Error(), "Close")
	}
}

func (suite *commonTest) Test_GetChainID() {
	res, err := suite.TestAPI.GetChainID()
	if err != nil {
		suite.FailNow(err.Error(), "GetChainID")
	}

	suite.Equal(res, config.ChainIDCYB)
}

func (suite *commonTest) Test_GetAccountBalances() {
	res, err := suite.TestAPI.GetAccountBalances(UserID2, AssetCYB)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountBalances 1")
	}

	suite.NotNil(res)
	//logging.Dump("balance cyb >", res)

	res, err = suite.TestAPI.GetAccountBalances(UserID2)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountBalances 2")
	}

	suite.NotNil(res)
	//logging.Dump("balances all >", res)
}

func (suite *commonTest) Test_GetAccounts() {
	res, err := suite.TestAPI.GetAccounts(UserID2) //, UserID3, UserID4)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccounts")
	}

	suite.NotNil(res)
	suite.Len(res, 1)

	//logging.Dump("get accounts >", res)
}

func (suite *commonTest) Test_GetFullAccounts() {
	res, err := suite.TestAPI.GetFullAccounts(UserID2)
	if err != nil {
		suite.FailNow(err.Error(), "GetFullAccounts")
	}

	suite.NotNil(res)
	suite.Len(res, 1)

	//logging.Dump("get full accounts >", res)
}

func (suite *commonTest) Test_GetObjects() {
	res, err := suite.TestAPI.GetObjects(
		UserID1,
		AssetCNY,
		BitAssetDataCNY,
		// LimitOrder1,
		// CallOrder1,
		// SettleOrder1,
		OperationHistory1,
		CommiteeMember1,
		Balance1,
	)

	if err != nil {
		suite.FailNow(err.Error(), "GetObjects")
	}

	suite.NotNil(res)
	suite.Len(res, 6)
	//logging.Dump("objects >", res)
}

func (suite *commonTest) Test_GetBlock() {
	res, err := suite.TestAPI.GetBlock(33217575) //26867161)
	if err != nil {
		suite.FailNow(err.Error(), "GetBlock")
	}

	suite.NotNil(res)
	//logging.Dump("get_block >", res)
}

func (suite *commonTest) Test_GetDynamicGlobalProperties() {
	res, err := suite.TestAPI.GetDynamicGlobalProperties()
	if err != nil {
		suite.FailNow(err.Error(), "GetDynamicGlobalProperties")
	}

	suite.NotNil(res)
	//logging.Dump("dynamic global properties >", res)
}

func (suite *commonTest) Test_GetAccountByName() {
	res, err := suite.TestAPI.GetAccountByName("openledger")
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountByName")
	}

	suite.NotNil(res)
	//logging.Dump("accounts >", res)
}

func (suite *commonTest) Test_GetTradeHistory() {
	dtTo := time.Now().UTC()

	dtFrom := dtTo.Add(-time.Hour * 24)
	res, err := suite.TestAPI.GetTradeHistory(AssetCYB, AssetHERO, dtTo, dtFrom, 50)

	if err != nil {
		suite.FailNow(err.Error(), "GetTradeHistory")
	}

	suite.NotNil(res)
	//logging.Dump("markettrades >", res)
}

func (suite *commonTest) Test_GetLimitOrders() {
	res, err := suite.TestAPI.GetLimitOrders(AssetCNY, AssetCYB, 50)
	if err != nil {
		suite.FailNow(err.Error(), "GetLimitOrders")
	}

	suite.NotNil(res)
	//logging.Dump("limitorders >", res)
}

func (suite *commonTest) Test_GetCallOrders() {
	res, err := suite.TestAPI.GetCallOrders(AssetUSD, 50)
	if err != nil {
		suite.FailNow(err.Error(), "GetCallOrders")
	}

	suite.NotNil(res)
	//	logging.Dump("callorders >", res)
}

func (suite *commonTest) Test_GetMarginPositions() {
	res, err := suite.TestAPI.GetMarginPositions(UserID2)
	if err != nil {
		suite.FailNow(err.Error(), "GetMarginPositions")
	}

	suite.NotNil(res)
	//logging.Dump("marginpositions >", res)
}

func (suite *commonTest) Test_GetSettleOrders() {
	res, err := suite.TestAPI.GetSettleOrders(AssetCNY, 50)
	if err != nil {
		suite.FailNow(err.Error(), "GetSettleOrders")
	}

	suite.NotNil(res)
	//logging.Dump("settleorders >", res)logging.SetDebug(true)
}

func (suite *commonTest) Test_ListAssets() {
	res, err := suite.TestAPI.ListAssets("OPEN.DASH", 2)
	if err != nil {
		suite.FailNow(err.Error(), "ListAssets")
	}

	suite.NotNil(res)
	suite.Len(res, 2)
	//logging.Dump("assets >", res)
}

func (suite *commonTest) Test_GetAccountHistory() {

	user := types.NewGrapheneID("1.2.96393")
	start := types.NewGrapheneID("1.11.187698971")
	stop := types.NewGrapheneID("1.11.187658388")

	res, err := suite.TestAPI.GetAccountHistory(user, stop, 30, start)
	if err != nil {
		suite.FailNow(err.Error(), "GetAccountHistory")
	}

	suite.NotNil(res)
	//logging.Dump("history >", res)
}

func (suite *commonTest) Test_GetOrderBook() {
	res, err := suite.TestAPI.GetOrderBook(AssetUSD, AssetCYB, 10)
	if err != nil {
		suite.FailNow(err.Error(), "GetOrderBook")
	}

	suite.NotNil(res)
	suite.Len(res.Asks, 10)
	suite.Len(res.Bids, 10)

	//logging.Dump("test GetOrderBook >", res)
}

func (suite *commonTest) Test_Get24Volume() {
	res, err := suite.TestAPI.Get24Volume(AssetUSD, AssetCYB)
	if err != nil {
		suite.FailNow(err.Error(), "Get24Volume")
	}

	suite.NotNil(res)
	suite.Equal(*AssetUSD, res.Base)
	suite.Equal(*AssetCYB, res.Quote)

	//logging.Dump("test Get24Volume >", res)
}

func TestCommon(t *testing.T) {
	testSuite := new(commonTest)
	suite.Run(t, testSuite)
}
