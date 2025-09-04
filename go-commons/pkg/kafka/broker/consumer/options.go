package consumer

import "time"

// Option -.
type Option func(*Consumer)

// Timeout -.
func Timeout(timeout time.Duration) Option {
	return func(c *Consumer) {
		c.timeout = timeout
	}
}

// ConnWaitTime -.
func ConnWaitTime(timeout time.Duration) Option {
	return func(c *Consumer) {
		c.waitTime = timeout
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Consumer) {
		c.attempts = attempts
	}
}
