package tracing

type ContextSnapshot interface {
	IsValid() bool
}
