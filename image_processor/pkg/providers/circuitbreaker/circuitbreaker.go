package circuitbreaker

import (
	"errors"
	"log"
	"time"

	"github.com/Skarlso/kube-cluster-sample/image_processor/facerecog"
)

// circuitBreaker defines the functionality of the circuit breaker.
type CircuitBreaker interface {
	Call() (*facerecog.IdentifyResponse, error)
	SetCallF(func() (*facerecog.IdentifyResponse, error))
	SetPingF(func() bool)
}

// circuitBreaker is a circuit of remote calls. The breaker is activated after
// a configured amount of failed tries which disallows all subsequent calls to
// said circuit. This is specific to Identify call thus the function will be
// specific to this application.
//
// Potentially if there are thousands of images that need processing
// and the endpoint for the processing becomes unstable, we will stop sending
// it images for a few seconds to give it a chance to recover.
//
// This could be further improved if we store the requests which didn't go through
// and re-process them after the circuit is alive again. But we leave that up to
// the caller for now.
type circuitBreaker struct {
	TimeOut          time.Duration
	CurrentBreakTime time.Time
	On               bool
	CurrentTries     int
	MaxTries         int
	F                func() (*facerecog.IdentifyResponse, error)
	Ping             func() bool
}

func (c *circuitBreaker) engage() {
	c.CurrentBreakTime = time.Now()
	c.On = true
}

func (c *circuitBreaker) disengage() {
	c.CurrentTries = 0
	c.On = false
}

func (c *circuitBreaker) checkIfOver() {
	if c.CurrentBreakTime.Add(c.TimeOut).Before(time.Now()) {
		log.Printf("timeout over. running ping.")
		if !c.Ping() {
			log.Println("backend still not functioning. extending break.")
			c.engage()
			return
		}
		c.disengage()
	}
}

// SetCallF adds the ability to define a calling function for the circuit breaker.
func (c *circuitBreaker) SetCallF(f func() (*facerecog.IdentifyResponse, error)) {
	c.F = f
}

// SetPingF adds the ability to define a ping function for the circuit breaker.
func (c *circuitBreaker) SetPingF(f func() bool) {
	c.Ping = f
}

// Call the function specified under F.
func (c *circuitBreaker) Call() (*facerecog.IdentifyResponse, error) {
	if c.On {
		c.checkIfOver()
	}
	if c.On {
		log.Printf("max sending try count of %d reached. sending not allowed for %v time period.", c.MaxTries, time.Until(c.CurrentBreakTime.Add(c.TimeOut)))
		return nil, errors.New("circuitbreaker is engaged")
	}
	r, err := c.F()
	if err != nil {
		c.CurrentTries++
		if c.CurrentTries >= c.MaxTries {
			log.Printf("maximum try of %d sends reached. disabling for %v time period.", c.MaxTries, c.TimeOut)
			c.engage()
		}
		return nil, err
	}
	c.CurrentTries = 0
	return r, err
}

// NewcircuitBreaker defines some default parameters for the breaker.
func NewcircuitBreaker() *circuitBreaker {
	c := circuitBreaker{
		CurrentTries: 0,
		MaxTries:     3,
		On:           false,
		TimeOut:      time.Second * 10,
	}
	return &c
}
