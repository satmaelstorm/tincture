package events

type BaseEvent struct {
	stop bool
}

func (e *BaseEvent) IsPropagationStopped() bool {
	return e.stop
}
