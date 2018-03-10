package limit

import (
	"context"
	"fmt"
	"net/http"
	"psychic-rat/log"
	"sync"
	"time"
)

type (
	Limiter struct {
		limits      map[string]*bucket
		mu          sync.Mutex
		max         int
		increment   int
		interval    int
		idGenerator IdGeneratorFunc
	}

	IdGeneratorFunc func(*http.Request) string

	bucket struct {
		tokens int
		mu     sync.Mutex
	}
)

func New(max, increment, interval int, idGen IdGeneratorFunc) *Limiter {
	return &Limiter{limits: map[string]*bucket{}, max: max, increment: increment, interval: interval, idGenerator: idGen}
}

func (l *Limiter) CheckLimit(r *http.Request) error {
	id := l.idGenerator(r)
	bh := l.getBucketFor(id)
	if !bh.consumeToken() {
		return fmt.Errorf("tokens exhausted for %s", id)
	}
	return nil
}

func (l *Limiter) getBucketFor(s string) *bucket {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.limits[s]
	if !ok {
		b = &bucket{tokens: l.max}
		l.limits[s] = b
		go b.tokenFiller(l.interval, l.increment, l.max, func() { l.deleteBucket(s) })
		log.Logf(context.Background(), "bucket filler started for %s", s)
	}
	return b
}

func (l *Limiter) deleteBucket(s string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.limits, s)
}

func (b *bucket) consumeToken() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.tokens == 0 {
		return false
	}
	b.tokens--
	return true
}

func (b *bucket) tokenFiller(interval, increment, max int, done func()) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		b.tokens += increment
		if b.tokens >= max {
			done()
			return
		}
	}
}
