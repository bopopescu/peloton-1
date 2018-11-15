package cli

import (
	"context"
	"testing"

	"code.uber.internal/infra/peloton/.gen/peloton/api/v1alpha/peloton"
	"code.uber.internal/infra/peloton/.gen/peloton/api/v1alpha/pod"
	podsvc "code.uber.internal/infra/peloton/.gen/peloton/api/v1alpha/pod/svc"
	"code.uber.internal/infra/peloton/.gen/peloton/api/v1alpha/pod/svc/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/yarpc/yarpcerrors"
)

const testPodName = "941ff353-ba82-49fe-8f80-fb5bc649b04d-1"

type podActionsTestSuite struct {
	suite.Suite
	ctx    context.Context
	client Client

	ctrl      *gomock.Controller
	podClient *mocks.MockPodServiceYARPCClient
}

func (suite *podActionsTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.podClient = mocks.NewMockPodServiceYARPCClient(suite.ctrl)
	suite.ctx = context.Background()
	suite.client = Client{
		Debug:      false,
		podClient:  suite.podClient,
		dispatcher: nil,
		ctx:        suite.ctx,
	}
}

func (suite *podActionsTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

// TestClientPodGetCacheSuccess test the success case of getting cache
func (suite *podActionsTestSuite) TestClientPodGetCacheSuccess() {
	suite.podClient.EXPECT().
		GetPodCache(gomock.Any(), gomock.Any()).
		Return(&podsvc.GetPodCacheResponse{
			Status: &pod.PodStatus{
				State: pod.PodState_POD_STATE_RUNNING,
			},
		}, nil)

	suite.NoError(suite.client.PodGetCacheAction(testPodName))
}

// TestClientPodGetCacheSuccess test the failure case of getting cache
func (suite *podActionsTestSuite) TestClientPodGetCacheFail() {
	suite.podClient.EXPECT().
		GetPodCache(gomock.Any(), gomock.Any()).
		Return(nil, yarpcerrors.InternalErrorf("test error"))

	suite.Error(suite.client.PodGetCacheAction(testPodName))
}

// TestPodGetEventsV1AlphaAction tests PodGetEventsV1AlphaAction
func (suite *podActionsTestSuite) TestPodGetEventsV1AlphaAction() {
	podname := &peloton.PodName{
		Value: "podname",
	}
	podID := &peloton.PodID{
		Value: "podID",
	}
	req := &podsvc.GetPodEventsRequest{
		PodName: podname,
		PodId:   podID,
	}

	podEvent := &pod.PodEvent{
		PodId: &peloton.PodID{
			Value: "podID",
		},
		PrevPodId: &peloton.PodID{
			Value: "prevPodID",
		},
		ActualState:  "PENDING",
		DesiredState: "RUNNING",
	}

	var podEvents []*pod.PodEvent
	podEvents = append(podEvents, podEvent)
	response := &podsvc.GetPodEventsResponse{
		Events: podEvents,
	}
	suite.podClient.EXPECT().GetPodEvents(context.Background(), req).
		Return(response, nil)
	err := suite.client.PodGetEventsV1AlphaAction(podname.GetValue(), podID.GetValue())
	suite.NoError(err)

	// Client with debug set to true
	suite.client.Debug = true
	suite.podClient.EXPECT().GetPodEvents(context.Background(), req).
		Return(response, nil)
	err = suite.client.PodGetEventsV1AlphaAction(podname.GetValue(), podID.GetValue())
	suite.NoError(err)
}

// TestPodGetEventsV1AlphaActionClientFailure tests GetPodEvents API error
func (suite *podActionsTestSuite) TestPodGetEventsV1AlphaActionAPIError() {
	podname := &peloton.PodName{
		Value: "podname",
	}
	podID := &peloton.PodID{
		Value: "podID",
	}
	req := &podsvc.GetPodEventsRequest{
		PodName: podname,
		PodId:   podID,
	}

	suite.podClient.EXPECT().GetPodEvents(suite.ctx, req).
		Return(nil, yarpcerrors.InternalErrorf("test GetPodEvents error"))
	err := suite.client.PodGetEventsV1AlphaAction(podname.GetValue(), podID.GetValue())
	suite.Error(err)
}

// TestClientPodRefreshSuccess test the success case of refreshing pod
func (suite *podActionsTestSuite) TestClientPodRefreshSuccess() {
	suite.podClient.EXPECT().
		RefreshPod(gomock.Any(), gomock.Any()).
		Return(&podsvc.RefreshPodResponse{}, nil)

	suite.NoError(suite.client.PodRefreshAction(testPodName))
}

// TestClientPodRefreshFail test the failure case of refreshing pod
func (suite *podActionsTestSuite) TestClientPodRefreshFail() {
	suite.podClient.EXPECT().
		RefreshPod(gomock.Any(), gomock.Any()).
		Return(nil, yarpcerrors.InternalErrorf("test error"))

	suite.Error(suite.client.PodRefreshAction(testPodName))
}

// TestClientPodStartSuccess test the success case of starting pod
func (suite *podActionsTestSuite) TestClientPodStartSuccess() {
	suite.podClient.EXPECT().
		StartPod(gomock.Any(), gomock.Any()).
		Return(&podsvc.StartPodResponse{}, nil)

	suite.NoError(suite.client.PodStartAction(testPodName))
}

// TestClientPodStartFail test the failure case starting pod
func (suite *podActionsTestSuite) TestClientPodStartFail() {
	suite.podClient.EXPECT().
		StartPod(gomock.Any(), gomock.Any()).
		Return(nil, yarpcerrors.InternalErrorf("test error"))

	suite.Error(suite.client.PodStartAction(testPodName))
}

// TestPodLogsGetActionSuccess tests failure of getting pod logs
// due to file not found error
func (suite *podActionsTestSuite) TestPodLogsGetActionFileNotFound() {
	podname := &peloton.PodName{
		Value: "podname",
	}
	podID := &peloton.PodID{
		Value: "podID",
	}
	req := &podsvc.BrowsePodSandboxRequest{
		PodName: podname,
		PodId:   podID,
	}
	resp := &podsvc.BrowsePodSandboxResponse{
		Hostname: "host1",
		Port:     "8000",
	}

	suite.podClient.EXPECT().
		BrowsePodSandbox(suite.ctx, req).
		Return(resp, nil)
	suite.Error(
		suite.client.PodLogsGetAction(
			"",
			podname.GetValue(),
			podID.GetValue(),
		),
	)
}

// TestPodLogsGetActionSuccess tests failure of getting pod logs
// due to BrowsePodSandbox API error
func (suite *podActionsTestSuite) TestPodLogsGetActionBrowsePodSandboxFailure() {
	suite.podClient.EXPECT().
		BrowsePodSandbox(suite.ctx, gomock.Any()).
		Return(nil, yarpcerrors.InternalErrorf("test error"))
	suite.Error(
		suite.client.PodLogsGetAction(
			"",
			"",
			"",
		),
	)
}

// TestPodLogsGetActionSuccess tests failure of getting pod logs
// due to error while downloading file
func (suite *podActionsTestSuite) TestPodLogsGetActionFileGetFailure() {
	podname := &peloton.PodName{
		Value: "podname",
	}
	podID := &peloton.PodID{
		Value: "podID",
	}
	filename := "filename"
	req := &podsvc.BrowsePodSandboxRequest{
		PodName: podname,
		PodId:   podID,
	}
	resp := &podsvc.BrowsePodSandboxResponse{
		Hostname: "host1",
		Port:     "8000",
		Paths:    []string{filename},
	}

	suite.podClient.EXPECT().
		BrowsePodSandbox(suite.ctx, req).
		Return(resp, nil)
	suite.Error(
		suite.client.PodLogsGetAction(
			filename,
			podname.GetValue(),
			podID.GetValue(),
		),
	)
}

func TestPodActions(t *testing.T) {
	suite.Run(t, new(podActionsTestSuite))
}
