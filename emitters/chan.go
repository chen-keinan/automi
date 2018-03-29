package emitters

import (
	"context"
	"errors"
	"log"
	"reflect"

	autoctx "github.com/vladimirvivien/automi/api/context"
)

// ChanEmitter is an emitter that takes in a channel and
// and sets it up as the source of the emitter .
type ChanEmitter struct {
	channel interface{}
	output  chan interface{}
	log     *log.Logger
}

// Chan creates new slice source
func Chan(channel interface{}) *ChanEmitter {
	return &ChanEmitter{
		channel: channel,
		output:  make(chan interface{}, 1024),
	}
}

//GetOutput returns the output channel of this source node
func (c *ChanEmitter) GetOutput() <-chan interface{} {
	return c.output
}

// Open opens the source node to start streaming data on its channel
func (c *ChanEmitter) Open(ctx context.Context) error {
	// ensure channel param is a chan type
	chanType := reflect.TypeOf(c.channel)
	if chanType.Kind() != reflect.Chan {
		return errors.New("ChanEmitter requires channel")
	}
	c.log = autoctx.GetLogger(ctx)
	c.log.Print("opening channel emitter")
	chanVal := reflect.ValueOf(c.channel)

	if !chanVal.IsValid() {
		return errors.New("invalid channel for ChanEmitter")
	}

	go func() {
		defer func() {
			close(c.output)
			c.log.Println("closing slice emitter")
		}()

		for {
			val, open := chanVal.Recv()
			if !open {
				return
			}
			c.output <- val.Interface()
		}
	}()
	return nil
}