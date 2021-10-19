package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/Skarlso/kube-cluster-sample/image_processor/facerecog"
)

func TestEngageDisengage(t *testing.T) {
	c := NewCircuitBreaker(zerolog.Logger{})
	c.engage()
	if !c.On {
		t.Fatal("want c.On = true, got: ", c.On)
	}
	c.CurrentTries = 5
	c.disengage()
	if c.On {
		t.Fatal("want c.On = false, got: ", c.On)
	}
	if c.CurrentTries != 0 {
		t.Fatal("want currentTries = 0; got: ", c.CurrentTries)
	}
}

func TestCheckIfOverPing(t *testing.T) {
	c := NewCircuitBreaker(zerolog.Logger{})
	c.TimeOut = 0 * time.Second
	before := time.Now()
	c.CurrentBreakTime = before

	t.Run("ping is false", func(t *testing.T) {
		c.Ping = func() bool {
			return false
		}
		c.checkIfOver()
		if !c.On {
			t.Fatal("circuit breaker should still be engaged")
		}
		if !c.CurrentBreakTime.After(before) {
			t.Fatal("current time should have been extended. was: ", c.CurrentBreakTime)
		}
	})
	t.Run("ping is true", func(t *testing.T) {
		c.Ping = func() bool {
			return true
		}
		c.checkIfOver()
		if c.On {
			t.Fatal("circuit breaker should have disengaged")
		}
	})
}

func TestCallShouldFailAfterXTimesOfFail(t *testing.T) {
	c := NewCircuitBreaker(zerolog.Logger{})
	c.TimeOut = 0 * time.Second
	before := time.Now()
	c.CurrentBreakTime = before
	c.MaxTries = 1
	c.F = func() (*facerecog.IdentifyResponse, error) {
		return nil, errors.New("test error")
	}
	if _, err := c.Call(); err == nil {
		t.Fatal("error should have been returned")
	}
	if !c.On {
		t.Fatal("should have engaged after one failure. was not on.")
	}
	// now try again with a passing call, this should disengage the breaker.
	c.F = func() (*facerecog.IdentifyResponse, error) {
		return &facerecog.IdentifyResponse{ImageName: "stuff.png"}, nil
	}
	c.Ping = func() bool {
		return true
	}
	if _, err := c.Call(); err != nil {
		t.Fatal(err)
	}
	if c.On {
		t.Fatal("should have disengaged after successful call. was still on.")
	}
}

func TestCallShouldNotFailIfFunctionWorks(t *testing.T) {
	c := NewCircuitBreaker(zerolog.Logger{})
	c.TimeOut = 0 * time.Second
	before := time.Now()
	c.CurrentBreakTime = before
	c.MaxTries = 1
	c.F = func() (*facerecog.IdentifyResponse, error) {
		return nil, nil
	}
	_, err := c.Call()
	if err != nil {
		t.Fatal(err)
	}
	if c.On {
		t.Fatal("circuit breaker should not have engaged.")
	}
}

func TestCallShouldReturnErrorInCaseTheBreakerIsOn(t *testing.T) {
	c := NewCircuitBreaker(zerolog.Logger{})
	c.TimeOut = 0 * time.Second
	before := time.Now()
	c.CurrentBreakTime = before
	c.MaxTries = 1
	c.On = true
	c.F = func() (*facerecog.IdentifyResponse, error) {
		return nil, errors.New("test error")
	}
	c.Ping = func() bool {
		return false
	}
	_, err := c.Call()
	if err == nil {
		t.Fatal("error should not have been empty")
	}
	if err.Error() != "circuitbreaker is engaged" {
		t.Fatal("wanted: 'circuitbreaker is engaged', got: ", err.Error())
	}
}
