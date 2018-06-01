package cli

import (
	"context"
	"strconv"
	"testing"

	mesos "code.uber.internal/infra/peloton/.gen/mesos/v1"
	"code.uber.internal/infra/peloton/.gen/peloton/private/hostmgr/hostsvc"
	"code.uber.internal/infra/peloton/util"
	log "github.com/sirupsen/logrus"

	hostMocks "code.uber.internal/infra/peloton/.gen/peloton/private/hostmgr/hostsvc/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

var (
	_cpuName     = "cpus"
	_pelotonRole = "peloton"
	_memName     = "mem"
	_diskName    = "disk"
	_gpuName     = "gpus"
	_portsName   = "ports"
	_testAgent   = "agent"

	_cpuRes = util.NewMesosResourceBuilder().
		WithName(_cpuName).
		WithValue(1.0).
		Build()
	_memRes = util.NewMesosResourceBuilder().
		WithName(_memName).
		WithValue(1.0).
		Build()
	_diskRes = util.NewMesosResourceBuilder().
			WithName(_diskName).
			WithValue(1.0).
			Build()
	_gpuRes = util.NewMesosResourceBuilder().
		WithName(_gpuName).
		WithValue(1.0).
		Build()
	_portsRes = util.NewMesosResourceBuilder().
			WithName(_portsName).
			WithRanges(util.CreatePortRanges(
			map[uint32]bool{1: true, 2: true})).
		Build()
)

type offersActionsTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	mockHostMgr *hostMocks.MockInternalHostServiceYARPCClient
	ctx         context.Context
}

func (suite *offersActionsTestSuite) SetupSuite() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockHostMgr = hostMocks.NewMockInternalHostServiceYARPCClient(suite.mockCtrl)
	suite.ctx = context.Background()
}

func (suite *offersActionsTestSuite) SetupTest() {
	log.Debug("SetupTest")
}

func (suite *offersActionsTestSuite) TearDownTest() {
	log.Debug("TearDownTest")
}

func (suite *offersActionsTestSuite) TearDownSuite() {
	suite.mockCtrl.Finish()
	suite.ctx.Done()
}

func (suite *offersActionsTestSuite) TestNoOutstandingOffers() {
	c := Client{
		Debug:         false,
		hostMgrClient: suite.mockHostMgr,
		dispatcher:    nil,
		ctx:           suite.ctx,
	}

	resp := &hostsvc.GetOutstandingOffersResponse{
		Offers: nil,
		Error: &hostsvc.GetOutstandingOffersResponse_Error{
			NoOffers: &hostsvc.NoOffersError{
				Message: "no offers present in offer pool",
			},
		},
	}

	suite.mockHostMgr.EXPECT().GetOutstandingOffers(
		gomock.Any(),
		&hostsvc.GetOutstandingOffersRequest{}).Return(resp, nil)

	suite.NoError(c.OffersGetAction())
}

func (suite *offersActionsTestSuite) TestGetOutstandingOffers() {
	c := Client{
		Debug:         false,
		hostMgrClient: suite.mockHostMgr,
		dispatcher:    nil,
		ctx:           suite.ctx,
	}

	resp := &hostsvc.GetOutstandingOffersResponse{
		Offers: suite.createUnreservedMesosOffers(5),
		Error:  nil,
	}

	suite.mockHostMgr.EXPECT().GetOutstandingOffers(
		gomock.Any(),
		&hostsvc.GetOutstandingOffersRequest{}).Return(resp, nil)

	suite.NoError(c.OffersGetAction())
}

func TestOffersAction(t *testing.T) {
	suite.Run(t, new(offersActionsTestSuite))
}

func (suite *offersActionsTestSuite) createUnreservedMesosOffer(
	offerID string) *mesos.Offer {
	rs := []*mesos.Resource{
		_cpuRes,
		_memRes,
		_diskRes,
		_gpuRes,
	}

	return &mesos.Offer{
		Id: &mesos.OfferID{
			Value: &offerID,
		},
		AgentId: &mesos.AgentID{
			Value: &_testAgent,
		},
		Hostname:  &_testAgent,
		Resources: rs,
	}
}

func (suite *offersActionsTestSuite) createUnreservedMesosOffers(count int) []*mesos.Offer {
	var offers []*mesos.Offer
	for i := 0; i < count; i++ {
		offers = append(offers, suite.createUnreservedMesosOffer("offer-id-"+strconv.Itoa(i)))
	}
	return offers
}