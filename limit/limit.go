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
		clients     map[string]*bucket
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
	return &Limiter{clients: map[string]*bucket{}, max: max, increment: increment, interval: interval, idGenerator: idGen}
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

	b, ok := l.clients[s]
	if !ok {
		b = &bucket{tokens: l.max}
		l.clients[s] = b
		go func() {
			log.Logf(context.Background(), "bucket filler started for %s", s)
			b.tokenFiller(l.interval, l.increment, l.max)
			l.deleteBucket(s)
			log.Logf(context.Background(), "bucket filler completed for %s", s)
		}()
	}
	return b
}

func (l *Limiter) deleteBucket(s string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.clients, s)
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

func (b *bucket) tokenFiller(interval, increment, max int) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		b.tokens += increment
		if b.tokens >= max {
			return
		}
	}
}
