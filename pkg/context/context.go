package context

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const KeyRequestID = "id"

type Context interface {
	context.Context

	WithDeadline(d time.Time)
	CopyWithDeadline(d time.Time) Context

	WithTimeout(timeout time.Duration)
	CopyWithTimeout(timeout time.Duration) Context
	Cancel()

	Copy() Context

	Value
}

type local struct {
	base       context.Context
	cancelFunc context.CancelFunc
}

func (l local) Deadline() (deadline time.Time, ok bool) {
	return l.base.Deadline()
}

func (l local) Done() <-chan struct{} {
	return l.base.Done()
}

func (l local) Err() error {
	return l.base.Err()
}

func (l *local) Cancel() {
	if l.cancelFunc != nil {
		l.cancelFunc()
	}
}

func (l local) Copy() Context {
	return &l
}

func (l *local) isEmptyID() bool {
	_, ok := l.id()
	return !ok
}

var cancelFunc = func() {}

func Empty() Context {
	ctx := &local{
		base:       context.Background(),
		cancelFunc: cancelFunc,
	}

	ctx.WithValue(KeyRequestID, uuid.New())

	return ctx
}

func New(option interface{}) Context {
	ctx := &local{
		base:       context.Background(),
		cancelFunc: cancelFunc,
	}

	switch baseCtx := option.(type) {
	case gin.Context:
		ctx.base = baseCtx.Request.Context()
		for key, value := range baseCtx.Keys {
			ctx.WithValue(key, value)
		}
	case *gin.Context:
		ctx.base = baseCtx.Request.Context()
		for key, value := range baseCtx.Keys {
			ctx.WithValue(key, value)
		}
	case Context:
		ctx.withValue(KeyRequestID, baseCtx.ID())
	case context.Context:
		ctx.base = baseCtx
	}

	if ctx.isEmptyID() {
		ctx.withValue(KeyRequestID, uuid.New().String())
	}

	return ctx
}

func NewWithTimeout(option interface{}, timeout time.Duration) Context {
	ctx := New(option)
	ctx.WithTimeout(timeout)

	return ctx
}
