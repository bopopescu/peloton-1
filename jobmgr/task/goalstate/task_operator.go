package goalstate

import (
	"context"

	"code.uber.internal/infra/peloton/.gen/mesos/v1"
	"code.uber.internal/infra/peloton/.gen/peloton/api/task"
	"code.uber.internal/infra/peloton/.gen/peloton/private/hostmgr/hostsvc"
	"code.uber.internal/infra/peloton/common"
	"go.uber.org/yarpc"
)

// TaskOperator can perform operations on a task, potentially resulting in the
// task changes state.
type TaskOperator interface {
	// StopTask by issueing a kill request. Even if the call is succesfull, the
	// task is not guaranteed to be killed.
	StopTask(ctx context.Context, taskInfo *task.TaskInfo) error
}

// NewTaskOperator from the set of arguments..
func NewTaskOperator(d *yarpc.Dispatcher) TaskOperator {
	return &taskOperator{
		hostmgrClient: hostsvc.NewInternalHostServiceYARPCClient(d.ClientConfig(common.PelotonHostManager)),
	}
}

type taskOperator struct {
	hostmgrClient hostsvc.InternalHostServiceYARPCClient
}

func (o *taskOperator) StopTask(ctx context.Context, taskInfo *task.TaskInfo) error {
	// TODO(mu): Notify RM to also remove these tasks from task queue.
	req := &hostsvc.KillTasksRequest{
		TaskIds: []*mesos_v1.TaskID{taskInfo.GetRuntime().GetMesosTaskId()},
	}
	_, err := o.hostmgrClient.KillTasks(ctx, req)
	return err
}