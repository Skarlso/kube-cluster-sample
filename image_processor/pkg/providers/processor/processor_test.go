package processor_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/Skarlso/kube-cluster-sample/image_processor/facerecog"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/models"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers/circuitbreaker"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers/fakes"
	"github.com/Skarlso/kube-cluster-sample/image_processor/pkg/providers/processor"
)

type mockIdentify struct {
	resp *facerecog.IdentifyResponse
	err  error
}

func (m *mockIdentify) Identify(ctx context.Context, in *facerecog.IdentifyRequest, opts ...grpc.CallOption) (*facerecog.IdentifyResponse, error) {
	return m.resp, m.err
}

type mockHealth struct {
	resp *facerecog.HealthCheckResponse
	err  error
}

func (m *mockHealth) HealthCheck(ctx context.Context, in *facerecog.Empty, opts ...grpc.CallOption) (*facerecog.HealthCheckResponse, error) {
	return m.resp, m.err
}

// TODO: Extract nice setup and check log error outputs.
func TestProcessImageSuccessfully(t *testing.T) {
	fakeStorer := &fakes.FakeImageStorer{}
	fakeStorer.GetImageReturns(&models.Image{
		ID:     1,
		Path:   "test/path",
		Person: 1,
		Status: models.PENDING,
	}, nil)
	fakeStorer.GetPersonFromImageReturns(&models.Person{
		ID:   1,
		Name: "Hannibal",
	}, nil)
	fakeStorer.UpdateImageReturns(nil)
	fakeIdentity := &mockIdentify{
		resp: &facerecog.IdentifyResponse{
			ImageName: "hannibal.jpg",
		},
	}
	fakeHealth := &mockHealth{
		resp: &facerecog.HealthCheckResponse{
			Ready: true,
		},
	}
	logger := zerolog.New(os.Stderr)
	cb := circuitbreaker.NewCircuitBreaker(logger)
	p := &processor.Processor{
		Dependencies: processor.Dependencies{
			Logger:         logger,
			Storer:         fakeStorer,
			CircuitBreaker: cb,
		},
		IdentifyClient:    fakeIdentity,
		HealthCheckClient: fakeHealth,
	}
	in := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.ProcessImages(ctx, in)
	in <- 1

	// and then check if the image has been processed and the fakes have been called.
	assert.Eventually(t, func() bool {
		if fakeStorer.GetImageCallCount() != 1 {
			return false
		}
		imageArg := fakeStorer.GetImageArgsForCall(0)
		assert.Equal(t, 1, imageArg)
		personArg := fakeStorer.GetPersonFromImageArgsForCall(0)
		assert.Equal(t, "hannibal.jpg", personArg)
		imageId, personId, status := fakeStorer.UpdateImageArgsForCall(1)
		assert.Equal(t, 1, imageId)
		assert.Equal(t, 1, personId)
		assert.Equal(t, models.PROCESSED, status)
		return true
	}, 5*time.Second, 10*time.Millisecond)
}

func TestProcessImageNotPending(t *testing.T) {
	fakeStorer := &fakes.FakeImageStorer{}
	fakeStorer.GetImageReturns(&models.Image{
		ID:     1,
		Path:   "test/path",
		Person: 1,
		Status: models.PROCESSED,
	}, nil)
	fakeIdentity := &mockIdentify{}
	fakeHealth := &mockHealth{}
	logger := zerolog.New(os.Stderr)
	cb := circuitbreaker.NewCircuitBreaker(logger)
	p := &processor.Processor{
		Dependencies: processor.Dependencies{
			Logger:         logger,
			Storer:         fakeStorer,
			CircuitBreaker: cb,
		},
		IdentifyClient:    fakeIdentity,
		HealthCheckClient: fakeHealth,
	}
	in := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.ProcessImages(ctx, in)
	in <- 1

	// and then check if the image has been processed and the fakes have been called.
	assert.Eventually(t, func() bool {
		if len(in) != 0 {
			return false
		}
		imageArg := fakeStorer.GetImageArgsForCall(0)
		assert.Equal(t, 1, imageArg)
		assert.Equal(t, 0, fakeStorer.GetPersonFromImageCallCount())
		assert.Equal(t, 0, fakeStorer.UpdateImageCallCount())
		return true
	}, 5*time.Second, 10*time.Millisecond)
}

func TestProcessImageFailedToGetImage(t *testing.T) {
	fakeStorer := &fakes.FakeImageStorer{}
	fakeStorer.GetImageReturns(nil, errors.New("nope"))
	fakeIdentity := &mockIdentify{}
	fakeHealth := &mockHealth{}
	logger := zerolog.New(os.Stderr)
	cb := circuitbreaker.NewCircuitBreaker(logger)
	p := &processor.Processor{
		Dependencies: processor.Dependencies{
			Logger:         logger,
			Storer:         fakeStorer,
			CircuitBreaker: cb,
		},
		IdentifyClient:    fakeIdentity,
		HealthCheckClient: fakeHealth,
	}
	in := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.ProcessImages(ctx, in)
	in <- 1

	// and then check if the image has been processed and the fakes have been called.
	assert.Eventually(t, func() bool {
		if len(in) != 0 {
			return false
		}
		imageArg := fakeStorer.GetImageArgsForCall(0)
		assert.Equal(t, 1, imageArg)
		assert.Equal(t, 0, fakeStorer.GetPersonFromImageCallCount())
		assert.Equal(t, 0, fakeStorer.UpdateImageCallCount())
		return true
	}, 5*time.Second, 10*time.Millisecond)
}

func TestProcessImageFailedToCallIdentityServiceCallsPing(t *testing.T) {
	fakeStorer := &fakes.FakeImageStorer{}
	fakeStorer.GetImageReturns(&models.Image{
		ID:     1,
		Path:   "test/path",
		Person: 1,
		Status: models.PENDING,
	}, nil)
	fakeStorer.GetPersonFromImageReturns(&models.Person{
		ID:   1,
		Name: "Hannibal",
	}, nil)
	fakeStorer.UpdateImageReturns(nil)
	fakeIdentity := &mockIdentify{
		err: errors.New("nope"),
	}
	fakeHealth := &mockHealth{
		resp: &facerecog.HealthCheckResponse{
			Ready: false,
		},
	}
	logger := zerolog.New(os.Stderr)
	cb := circuitbreaker.NewCircuitBreaker(logger)
	p := &processor.Processor{
		Dependencies: processor.Dependencies{
			Logger:         logger,
			Storer:         fakeStorer,
			CircuitBreaker: cb,
		},
		IdentifyClient:    fakeIdentity,
		HealthCheckClient: fakeHealth,
	}
	in := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.ProcessImages(ctx, in)
	in <- 1

	// and then check if the image has been processed and the fakes have been called.
	assert.Eventually(t, func() bool {
		if fakeStorer.GetImageCallCount() != 1 {
			return false
		}
		imageArg := fakeStorer.GetImageArgsForCall(0)
		assert.Equal(t, 1, imageArg)
		imageId, personId, status := fakeStorer.UpdateImageArgsForCall(1)
		assert.Equal(t, 1, imageId)
		assert.Equal(t, -1, personId)
		assert.Equal(t, models.FAILEDPROCESSING, status)
		return true
	}, 5*time.Second, 10*time.Millisecond)
}

func TestProcessImageImageNameNotFound(t *testing.T) {
	fakeStorer := &fakes.FakeImageStorer{}
	fakeStorer.GetImageReturns(&models.Image{
		ID:     1,
		Path:   "test/path",
		Person: 1,
		Status: models.PENDING,
	}, nil)
	fakeStorer.GetPersonFromImageReturns(&models.Person{
		ID:   1,
		Name: "Hannibal",
	}, nil)
	fakeStorer.UpdateImageReturns(nil)
	fakeIdentity := &mockIdentify{
		resp: &facerecog.IdentifyResponse{
			ImageName: "not_found",
		},
	}
	fakeHealth := &mockHealth{
		resp: &facerecog.HealthCheckResponse{
			Ready: true,
		},
	}
	logger := zerolog.New(os.Stderr)
	cb := circuitbreaker.NewCircuitBreaker(logger)
	p := &processor.Processor{
		Dependencies: processor.Dependencies{
			Logger:         logger,
			Storer:         fakeStorer,
			CircuitBreaker: cb,
		},
		IdentifyClient:    fakeIdentity,
		HealthCheckClient: fakeHealth,
	}
	in := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.ProcessImages(ctx, in)
	in <- 1

	// and then check if the image has been processed and the fakes have been called.
	assert.Eventually(t, func() bool {
		if fakeStorer.GetImageCallCount() != 1 {
			return false
		}
		imageArg := fakeStorer.GetImageArgsForCall(0)
		assert.Equal(t, 1, imageArg)
		imageId, personId, status := fakeStorer.UpdateImageArgsForCall(1)
		assert.Equal(t, 1, imageId)
		assert.Equal(t, -1, personId)
		assert.Equal(t, models.FAILEDPROCESSING, status)
		return true
	}, 5*time.Second, 10*time.Millisecond)
}

func TestProcessImagePersonNotFound(t *testing.T) {
	fakeStorer := &fakes.FakeImageStorer{}
	fakeStorer.GetImageReturns(&models.Image{
		ID:     1,
		Path:   "test/path",
		Person: 1,
		Status: models.PENDING,
	}, nil)
	fakeStorer.GetPersonFromImageReturns(nil, errors.New("nope"))
	fakeStorer.UpdateImageReturns(nil)
	fakeIdentity := &mockIdentify{
		resp: &facerecog.IdentifyResponse{
			ImageName: "hannibal1.jpg",
		},
	}
	fakeHealth := &mockHealth{
		resp: &facerecog.HealthCheckResponse{
			Ready: true,
		},
	}
	logger := zerolog.New(os.Stderr)
	cb := circuitbreaker.NewCircuitBreaker(logger)
	p := &processor.Processor{
		Dependencies: processor.Dependencies{
			Logger:         logger,
			Storer:         fakeStorer,
			CircuitBreaker: cb,
		},
		IdentifyClient:    fakeIdentity,
		HealthCheckClient: fakeHealth,
	}
	in := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.ProcessImages(ctx, in)
	in <- 1

	// and then check if the image has been processed and the fakes have been called.
	assert.Eventually(t, func() bool {
		if fakeStorer.GetImageCallCount() != 1 {
			return false
		}
		imageArg := fakeStorer.GetImageArgsForCall(0)
		assert.Equal(t, 1, imageArg)
		assert.Equal(t, 1, fakeStorer.UpdateImageCallCount())
		return true
	}, 5*time.Second, 10*time.Millisecond)
}

func TestProcessImageFailedToUpdateImage(t *testing.T) {
	fakeStorer := &fakes.FakeImageStorer{}
	fakeStorer.GetImageReturns(&models.Image{
		ID:     1,
		Path:   "test/path",
		Person: 1,
		Status: models.PENDING,
	}, nil)
	fakeStorer.GetPersonFromImageReturns(&models.Person{
		ID:   1,
		Name: "Hannibal",
	}, nil)
	fakeStorer.UpdateImageReturnsOnCall(1, errors.New("nope"))
	fakeIdentity := &mockIdentify{
		resp: &facerecog.IdentifyResponse{
			ImageName: "hannibal.jpg",
		},
	}
	fakeHealth := &mockHealth{
		resp: &facerecog.HealthCheckResponse{
			Ready: true,
		},
	}
	buf := &bytes.Buffer{}
	logger := zerolog.New(buf)
	cb := circuitbreaker.NewCircuitBreaker(logger)
	p := &processor.Processor{
		Dependencies: processor.Dependencies{
			Logger:         logger,
			Storer:         fakeStorer,
			CircuitBreaker: cb,
		},
		IdentifyClient:    fakeIdentity,
		HealthCheckClient: fakeHealth,
	}
	in := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.ProcessImages(ctx, in)
	in <- 1

	// and then check if the image has been processed and the fakes have been called.
	assert.Eventually(t, func() bool {
		if fakeStorer.GetImageCallCount() != 1 {
			return false
		}
		imageArg := fakeStorer.GetImageArgsForCall(0)
		assert.Equal(t, 1, imageArg)
		personArg := fakeStorer.GetPersonFromImageArgsForCall(0)
		assert.Equal(t, "hannibal.jpg", personArg)
		imageId, personId, status := fakeStorer.UpdateImageArgsForCall(1)
		assert.Equal(t, 1, imageId)
		assert.Equal(t, 1, personId)
		assert.Equal(t, models.PROCESSED, status)
		assert.True(t, bytes.Contains(buf.Bytes(), []byte("warning: could not update image record")))
		return true
	}, 5*time.Second, 10*time.Millisecond)
}
