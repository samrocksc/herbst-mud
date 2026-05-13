package effects

// dispatchMessage sends a text message to a character via the game's
// message system. This is called when a message effect fires.
func (s *Service) dispatchMessage(charID int, text string, msgType string) {
	if text == "" {
		return
	}
	// Game chat commands (say, yell, shout, whisper) should be passed through as-is
	// since they represent game commands that will be sent to the server.
	// Other valid message types for in-game display: info, success, error, damage, heal
	if msgType == "" {
		msgType = "info"
	}
	// Messages are dispatched through the game model's message system.
	// The game model will pick up pending messages via the MessageBus.
	s.messageBus.Send(charID, text, msgType)
}

// dispatchStartMessage sends the on_start message from an effect's messages field.
func dispatchStartMessage(msgs map[string]string, bus *MessageBus, charID int) {
	if text, ok := msgs["on_start"]; ok && text != "" {
		bus.Send(charID, text, "info")
	}
}

// dispatchEndMessage sends the on_end message from an effect's messages field.
func dispatchEndMessage(msgs map[string]string, bus *MessageBus, charID int) {
	if text, ok := msgs["on_end"]; ok && text != "" {
		bus.Send(charID, text, "info")
	}
}