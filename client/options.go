package client

import "time"

type Option struct {
	watchInterval time.Duration
}

type Options func(o *Option)

func WithWatchInterval(interval time.Duration) Options {
	return func(o *Option) {
		o.watchInterval = interval
	}
}
