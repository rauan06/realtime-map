package producer

import "time"

// Option -.
type Option func(*Producer)

// Timeout -.
func Timeout(timeout time.Duration) Option {
	return func(p *Producer) {
		p.timeout = timeout
	}
}

// ConnWaitTime -.
func ConnWaitTime(timeout time.Duration) Option {
	return func(p *Producer) {
		p.waitTime = timeout
	}
}

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(p *Producer) {
		p.attempts = attempts
	}
}
