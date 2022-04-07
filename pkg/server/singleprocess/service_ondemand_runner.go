package singleprocess

import (
	"context"

	empty "google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
)

func (s *Service) UpsertOnDemandRunnerConfig(
	ctx context.Context,
	req *pb.UpsertOnDemandRunnerConfigRequest,
) (*pb.UpsertOnDemandRunnerConfigResponse, error) {
	if err := serverptypes.ValidateUpsertOnDemandRunnerConfigRequest(req); err != nil {
		return nil, err
	}

	if req.Config.TargetRunner == nil {
		req.Config.TargetRunner = &pb.Ref_Runner{
			Target: &pb.Ref_Runner_Any{},
		}
	}
	result := req.Config
	if err := s.state(ctx).OnDemandRunnerConfigPut(result); err != nil {
		return nil, err
	}

	return &pb.UpsertOnDemandRunnerConfigResponse{Config: result}, nil
}

func (s *Service) GetOnDemandRunnerConfig(
	ctx context.Context,
	req *pb.GetOnDemandRunnerConfigRequest,
) (*pb.GetOnDemandRunnerConfigResponse, error) {
	if err := serverptypes.ValidateGetOnDemandRunnerConfigRequest(req); err != nil {
		return nil, err
	}

	result, err := s.state(ctx).OnDemandRunnerConfigGet(req.Config)
	if err != nil {
		return nil, err
	}

	return &pb.GetOnDemandRunnerConfigResponse{
		Config: result,
	}, nil
}

func (s *Service) ListOnDemandRunnerConfigs(
	ctx context.Context,
	req *empty.Empty,
) (*pb.ListOnDemandRunnerConfigsResponse, error) {
	result, err := s.state(ctx).OnDemandRunnerConfigList()
	if err != nil {
		return nil, err
	}

	return &pb.ListOnDemandRunnerConfigsResponse{Configs: result}, nil
}
