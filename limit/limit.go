package limit

import (
	"context"
	"fmt"
	"net/http"
	"psychic-rat/log"
	"strings"
	"sync"
	"time"
)

type (
	Limiter struct {
		limits map[string]*bucket
		mu     sync.Mutex
	}

	bucket struct {
		tokens       int
		max          int
		fillWith     int
		fillInterval int
		mu           sync.Mutex
	}
)

func NewBucket(path string, max, fillWith, fillInterval int) *bucket {
	return &bucket{max, max, fillWith, fillInterval, sync.Mutex{}}
}

func New(bs ...bucket) *Limiter {
	return &Limiter{limits: map[string]*bucket{}}
}

func (l *Limiter) CheckLimit(r *http.Request) error {
	id := r.Method + strings.Split(r.RemoteAddr, ":")[0]
	bh := l.getBucketHandler(id)
	if !bh.getToken() {
		return fmt.Errorf("tokens exhausted for %s", id)
	}
	log.Logf(r.Context(), "%s: %v", id, bh)
	return nil
}

func (l *Limiter) getBucketHandler(s string) *bucket {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.limits[s]
	if !ok {
		b = &bucket{15, 15, 5, 5, sync.Mutex{}}
		l.limits[s] = b
		go b.tokenFiller()
		log.Logf(context.Background(), "bucket filler started for %s", s)
	}
	return b
}

func (b *bucket) getToken() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.tokens == 0 {
		return false
	}
	b.tokens--
	return true
}

func (b *bucket) tokenFiller() {
	for {
		time.Sleep(time.Duration(b.fillInterval) * time.Second)
		b.tokens += b.fillWith
		log.Logf(context.Background(), "tokens added")
		if b.tokens >= b.max {
			log.Logf(context.Background(), "bucket is full")
			return
		}
	}
}
