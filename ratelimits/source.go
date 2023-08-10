package ratelimits

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ErrBucketNotFound indicates that the bucket was not found.
var ErrBucketNotFound = fmt.Errorf("bucket not found")

// source is an interface for creating and modifying TATs.
type source interface {
	// Set stores the TAT at the specified bucketKey ('name:id').
	Set(ctx context.Context, bucketKey string, tat time.Time) error

	// Get retrieves the TAT at the specified bucketKey ('name:id').
	Get(ctx context.Context, bucketKey string) (time.Time, error)

	// Delete deletes the TAT at the specified bucketKey ('name:id').
	Delete(ctx context.Context, bucketKey string) error
}

// inmem is an in-memory implementation of the source interface used for
// testing.
type inmem struct {
	sync.RWMutex
	m map[string]time.Time
}

func newInmem() *inmem {
	return &inmem{m: make(map[string]time.Time)}
}

func (in *inmem) Set(_ context.Context, bucketKey string, tat time.Time) error {
	in.Lock()
	defer in.Unlock()
	in.m[bucketKey] = tat
	return nil
}

func (in *inmem) Get(_ context.Context, bucketKey string) (time.Time, error) {
	in.RLock()
	defer in.RUnlock()
	tat, ok := in.m[bucketKey]
	if !ok {
		return time.Time{}, ErrBucketNotFound
	}
	return tat, nil
}

func (in *inmem) Delete(_ context.Context, bucketKey string) error {
	in.Lock()
	defer in.Unlock()
	delete(in.m, bucketKey)
	return nil
}
