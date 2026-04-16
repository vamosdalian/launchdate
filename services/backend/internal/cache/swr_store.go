package cache

import (
	"sync"
	"time"
)

type Clock func() time.Time

type Policy struct {
	SoftTTL time.Duration
	HardTTL time.Duration
}

type EntryState int

const (
	EntryStateFresh EntryState = iota
	EntryStateStale
	EntryStateExpired
)

type Item struct {
	Payload       []byte
	Negative      bool
	StoredAt      time.Time
	SoftExpiresAt time.Time
	HardExpiresAt time.Time
}

type SWRStoreOptions struct {
	MaxEntries int
	Clock      Clock
}

type SWRStore struct {
	mu         sync.RWMutex
	entries    map[string]Item
	refreshing map[string]struct{}
	versions   map[string]uint64
	maxEntries int
	now        Clock
}

func NewSWRStore(opts SWRStoreOptions) *SWRStore {
	clock := opts.Clock
	if clock == nil {
		clock = time.Now
	}

	return &SWRStore{
		entries:    make(map[string]Item),
		refreshing: make(map[string]struct{}),
		versions:   make(map[string]uint64),
		maxEntries: opts.MaxEntries,
		now:        clock,
	}
}

func (s *SWRStore) Get(key string) (Item, EntryState, bool) {
	now := s.now().UTC()

	s.mu.RLock()
	item, ok := s.entries[key]
	s.mu.RUnlock()
	if !ok {
		return Item{}, EntryStateExpired, false
	}

	return cloneItem(item), classifyEntry(now, item), true
}

func (s *SWRStore) Set(key string, payload []byte, negative bool, policy Policy) {
	now := s.now().UTC()
	item := Item{
		Payload:       cloneBytes(payload),
		Negative:      negative,
		StoredAt:      now,
		SoftExpiresAt: now.Add(policy.SoftTTL),
		HardExpiresAt: now.Add(policy.HardTTL),
	}
	if policy.SoftTTL <= 0 {
		item.SoftExpiresAt = now
	}
	if policy.HardTTL <= 0 {
		item.HardExpiresAt = now
	}
	if item.HardExpiresAt.Before(item.SoftExpiresAt) {
		item.HardExpiresAt = item.SoftExpiresAt
	}

	s.mu.Lock()
	s.entries[key] = item
	s.pruneLocked(now)
	s.mu.Unlock()
}

func (s *SWRStore) Delete(key string) {
	s.mu.Lock()
	delete(s.entries, key)
	delete(s.refreshing, key)
	s.mu.Unlock()
}

func (s *SWRStore) TryStartRefresh(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.refreshing[key]; ok {
		return false
	}

	s.refreshing[key] = struct{}{}
	return true
}

func (s *SWRStore) FinishRefresh(key string) {
	s.mu.Lock()
	delete(s.refreshing, key)
	s.mu.Unlock()
}

func (s *SWRStore) NamespaceVersion(namespace string) uint64 {
	s.mu.RLock()
	version := s.versions[namespace]
	s.mu.RUnlock()
	return version
}

func (s *SWRStore) BumpNamespace(namespace string) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.versions[namespace]++
	return s.versions[namespace]
}

func (s *SWRStore) pruneLocked(now time.Time) {
	if s.maxEntries <= 0 || len(s.entries) <= s.maxEntries {
		return
	}

	for key, item := range s.entries {
		if classifyEntry(now, item) != EntryStateExpired {
			continue
		}
		delete(s.entries, key)
		if len(s.entries) <= s.maxEntries {
			return
		}
	}

	for len(s.entries) > s.maxEntries {
		var oldestKey string
		var oldestTime time.Time
		first := true
		for key, item := range s.entries {
			if first || item.StoredAt.Before(oldestTime) {
				oldestKey = key
				oldestTime = item.StoredAt
				first = false
			}
		}
		delete(s.entries, oldestKey)
	}
}

func classifyEntry(now time.Time, item Item) EntryState {
	switch {
	case !now.After(item.SoftExpiresAt):
		return EntryStateFresh
	case !now.After(item.HardExpiresAt):
		return EntryStateStale
	default:
		return EntryStateExpired
	}
}

func cloneItem(item Item) Item {
	item.Payload = cloneBytes(item.Payload)
	return item
}

func cloneBytes(payload []byte) []byte {
	if len(payload) == 0 {
		return nil
	}
	cloned := make([]byte, len(payload))
	copy(cloned, payload)
	return cloned
}
