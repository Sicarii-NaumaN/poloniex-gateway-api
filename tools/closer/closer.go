// Closes help to close all external connections such as sql server, node etc.

package closer

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var globalCloser *closer

func init() {
	globalCloser = New(syscall.SIGTERM, syscall.SIGINT)
}

type closer struct {
	mu      sync.Mutex
	done    chan struct{}
	toClose []Closer
	once    sync.Once
}

type Closer interface {
	Close()
}

func New(sig ...os.Signal) *closer {
	c := &closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.closeAll()
		}()
	}
	return c
}

// Add append connections to be closed
func Add(conn ...Closer) {
	globalCloser.add(conn...)
}

func CloseAll() {
	globalCloser.closeAll()
}

func Wait() {
	globalCloser.wait()
}

func (c *closer) wait() {
	<-c.done
}

func (c *closer) add(conn ...Closer) {
	c.mu.Lock()

	c.toClose = append(c.toClose, conn...)

	c.mu.Unlock()
}

func (c *closer) closeAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.toClose
		c.toClose = nil
		c.mu.Unlock()

		var wg sync.WaitGroup
		for _, f := range funcs {
			wg.Add(1)
			go func(wg *sync.WaitGroup, f Closer) {
				f.Close()
				wg.Done()
			}(&wg, f)
		}
		wg.Wait()
	})
}
