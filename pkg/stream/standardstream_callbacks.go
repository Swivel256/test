package stream

func (s *StandardStream) OnStart(cb func()) {
	s.startCallbacks = append(s.startCallbacks, cb)
}

func (s *StandardStream) EmitStart() {
	for _, cb := range s.startCallbacks {
		cb()
	}
}

func (s *StandardStream) OnConnect(cb func()) {
	s.connectCallbacks = append(s.connectCallbacks, cb)
}

func (s *StandardStream) EmitConnect() {
	for _, cb := range s.connectCallbacks {
		cb()
	}
}

func (s *StandardStream) OnDisconnect(cb func()) {
	s.disconnectCallbacks = append(s.disconnectCallbacks, cb)
}

func (s *StandardStream) EmitDisconnect() {
	for _, cb := range s.disconnectCallbacks {
		cb()
	}
}

type StandardStreamEventHub interface {
	OnStart(cb func())

	OnConnect(cb func())

	OnDisconnect(cb func())

	OnTradeUpdate(cb func())
}
