package singleprocess

import (
	"context"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) UpsertTask(
	ctx context.Context,
	req *pb.UpsertTaskRequest,
) (*pb.UpsertTaskResponse, error) {
	if err := serverptypes.ValidateUpsertTaskRequest(req); err != nil {
		return nil, err
	}

	result := req.Task
	if err := s.state.TaskPut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertTaskResponse{Task: result}, nil
}

// GetTask returns a Task based on ID
func (s *service) GetTask(
	ctx context.Context,
	req *pb.GetTaskRequest,
) (*pb.GetTaskResponse, error) {
	if err := serverptypes.ValidateGetTaskRequest(req); err != nil {
		return nil, err
	}

	t, err := s.state.TaskGet(req.Ref)
	if err != nil {
		return nil, err
	}

	// Get the Start, Run, and Stop jobs
	resp, err := s.getJobsByTaskRef(ctx, t)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *service) ListTask(
	ctx context.Context,
	req *pb.ListTaskRequest,
) (*pb.ListTaskResponse, error) {
	// NOTE: no ptype validation at the moment, there are no request params

	result, err := s.state.TaskList()
	if err != nil {
		return nil, err
	}

	var tasks []*pb.GetTaskResponse
	for _, t := range result {
		tsk, err := s.getJobsByTaskRef(ctx, t)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, tsk)
	}

	return &pb.ListTaskResponse{Tasks: tasks}, nil
}

func (s *service) setTaskState(
	ctx context.Context,
	t *pb.GetTaskResponse,
) (*pb.GetTaskResponse, error) {
	if t.StartJob == nil {
		return nil, status.Error(codes.FailedPrecondition, "Task StartJob is nil, cannot determine Task state")
	} else if t.TaskJob == nil {
		return nil, status.Error(codes.FailedPrecondition, "TaskJob is nil, cannot determine Task state")
	} else if t.StopJob == nil {
		return nil, status.Error(codes.FailedPrecondition, "Task StartJob is nil, cannot determine Task state")
	}

	if t.StartJob.State == pb.Job_UNKNOWN && t.TaskJob.State == pb.Job_UNKNOWN &&
		t.StopJob.State == pb.Job_UNKNOWN {
		t.State = pb.GetTaskResponse_UNKNOWN
	} else if t.StartJob.State <= 2 && t.TaskJob.State <= 2 && t.StopJob.State <= 2 {
		// Job_State 2 and lower means no jobs have ran yet, at least waiting for a runner
		t.State = pb.GetTaskResponse_PENDING
	}

	if t.StartJob.State == pb.Job_RUNNING {
		t.State = pb.GetTaskResponse_STARTING
	} else if t.StartJob.State == pb.Job_SUCCESS {
		t.State = pb.GetTaskResponse_STARTED
	}

	if t.TaskJob.State == pb.Job_RUNNING {
		t.State = pb.GetTaskResponse_RUNNING
	} else if t.TaskJob.State == pb.Job_SUCCESS {
		t.State = pb.GetTaskResponse_COMPLETED
	}

	if t.StopJob.State == pb.Job_RUNNING {
		t.State = pb.GetTaskResponse_STOPPING
	} else if t.StopJob.State == pb.Job_SUCCESS {
		t.State = pb.GetTaskResponse_STOPPED
	}

	return t, nil
}

func (s *service) getJobsByTaskRef(
	ctx context.Context,
	t *pb.Task,
) (*pb.GetTaskResponse, error) {
	var taskJob, startJob, stopJob *pb.Job

	if t.TaskJob != nil {
		var err error
		taskJob, err = s.GetJob(ctx, &pb.GetJobRequest{JobId: t.TaskJob.Id})
		if err != nil {
			return nil, err
		}
	}

	if t.StartJob == nil {
		var err error
		startJob, err = s.GetJob(ctx, &pb.GetJobRequest{JobId: t.StartJob.Id})
		if err != nil {
			return nil, err
		}
	}

	if t.StopJob == nil {
		var err error
		stopJob, err = s.GetJob(ctx, &pb.GetJobRequest{JobId: t.StopJob.Id})
		if err != nil {
			return nil, err
		}
	}

	return &pb.GetTaskResponse{Task: t, TaskJob: taskJob, StartJob: startJob, StopJob: stopJob}, nil
}
