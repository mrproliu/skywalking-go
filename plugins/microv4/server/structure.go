package server

import "github.com/apache/skywalking-go/plugins/core/tracing"

//skywalking:ref_generate go-micro.dev/v4/util/socket InjectData
type InjectData struct {
	Span     tracing.Span
	Snapshot tracing.ContextSnapshot
}
