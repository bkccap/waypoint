package singleprocess

import (
	"context"

	configpkg "github.com/hashicorp/waypoint/pkg/config"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (s *Service) WaypointHclFmt(
	ctx context.Context,
	req *pb.WaypointHclFmtRequest,
) (*pb.WaypointHclFmtResponse, error) {
	result, err := configpkg.Format(req.WaypointHcl, "<input>")
	if err != nil {
		return nil, err
	}

	return &pb.WaypointHclFmtResponse{WaypointHcl: result}, nil
}
