package middleware

import (
	"context"
	"time"

	"github.com/go-martini/martini"
)

type TimeoutContext interface {
	context.Context
	MartiniContext() martini.Context
}

type timeoutContext struct {
	context.Context
	martiniContext martini.Context
}

func (t *timeoutContext) MartiniContext() martini.Context {
	return t.martiniContext
}

// WithFiveMinuteTimout adds a five minute timeout context to the request.
// it explicitly ignores the context of the http request.
// we do not want to stop the restarting of a container when a request cancels.
func WithFiveMinuteTimeout(c martini.Context) {
	tc, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Minute))
	t := &timeoutContext{Context: tc, martiniContext: c}
	c.Map(t)
}
