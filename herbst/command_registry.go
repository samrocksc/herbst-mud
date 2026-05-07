package main

// CommandHandler is a function that handles a command
type CommandHandler func(*model, []string)

// CommandRegistry manages command handlers and aliases
type CommandRegistry struct {
	handlers map[string]CommandHandler
	aliases  map[string]string
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		handlers: make(map[string]CommandHandler),
		aliases:  make(map[string]string),
	}
}

// Register adds a command handler with optional aliases
func (r *CommandRegistry) Register(name string, handler CommandHandler, aliases ...string) {
	r.handlers[name] = handler
	for _, alias := range aliases {
		r.aliases[alias] = name
	}
}

// Execute runs a command handler if found
// Returns true if command was handled, false otherwise
func (r *CommandRegistry) Execute(m *model, cmd string, args []string) bool {
	canonical := r.resolveCommand(cmd)
	if handler, ok := r.handlers[canonical]; ok {
		handler(m, args)
		return true
	}
	return false
}

// resolveCommand converts alias to canonical name
func (r *CommandRegistry) resolveCommand(cmd string) string {
	if canonical, ok := r.aliases[cmd]; ok {
		return canonical
	}
	return cmd
}

// IsRegistered checks if a command is registered
func (r *CommandRegistry) IsRegistered(cmd string) bool {
	canonical := r.resolveCommand(cmd)
	_, ok := r.handlers[canonical]
	return ok
}