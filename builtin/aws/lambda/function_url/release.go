package function_url

import "github.com/hashicorp/waypoint-plugin-sdk/component"

func (r *Release) URL() string { return r.Url }

var _ component.Release = (*Release)(nil)
