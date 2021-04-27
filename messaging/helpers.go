package messaging

func sendPulse(heartbeat chan interface{}) {
	select {
	case heartbeat <- struct{}{}:
	default:
	}
}
