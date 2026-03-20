package main

// ============================================================
// SCROLLING - Message history scroll intercept logic
// ============================================================

// handleScrollUp scrolls the message history up (older messages)
func (m *model) handleScrollUp() {
	if !m.isScrolling {
		m.isScrolling = true
		m.historyOffset = 1
	} else {
		m.historyOffset++
	}
	maxOffset := len(m.messageHistory) - 1
	if m.historyOffset > maxOffset {
		m.historyOffset = maxOffset
	}
}

// handleScrollDown scrolls the message history down (newer messages)
func (m *model) handleScrollDown() {
	if !m.isScrolling {
		return
	}
	m.historyOffset--
	if m.historyOffset < 0 {
		m.historyOffset = 0
		m.isScrolling = false
	}
}
