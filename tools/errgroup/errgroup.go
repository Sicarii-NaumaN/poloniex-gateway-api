package errgroup

import (
	"errors"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"runtime/debug"

	"go.uber.org/zap"
	eg "golang.org/x/sync/errgroup"
)

// ErrPanic is an error occurred when goroutine recovers from panic
var ErrPanic = errors.New("recovered from panic")

type Group struct {
	g   eg.Group
	log *zap.SugaredLogger
}

//  Go calls the given function in a new goroutine.
func (g *Group) Go(f func() error) {
	g.g.Go(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("errgroup: %w: %v", ErrPanic, r)
				g.logger().With("stack", string(debug.Stack())).Error(err.Error())
			}
		}()
		err = f()
		return
	})
}

func (g *Group) logger() *zap.SugaredLogger {
	if g.log != nil {
		return g.log
	}
	return logger.Logger()
}

// Wait blocks until all function calls from the Go method have returned, then returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	return g.g.Wait()
}
