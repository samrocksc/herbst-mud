package combat

// Tick configuration constants
const (
	// DefaultTickInterval is the default time between ticks (1.5 seconds)
	DefaultTickInterval = 1500 // milliseconds
	// FastTick is the fast combat speed (1 second)
	FastTick = 1000
	// SlowTick is the slow/tactical combat speed (2 seconds)
	SlowTick = 2000
)

// Config holds combat configuration
type Config struct {
	TickIntervalMs int64 `json:"tick_interval_ms"`
}

// DefaultConfig returns the default combat configuration
func DefaultConfig() *Config {
	return &Config{
		TickIntervalMs: DefaultTickInterval,
	}
}