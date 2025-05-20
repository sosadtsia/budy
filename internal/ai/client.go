package ai

// Client defines the interface for AI services
type Client interface {
	Ask(query string) error
}
