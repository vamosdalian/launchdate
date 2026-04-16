package core

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vamosdalian/launchdate-backend/internal/cache"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCanonicalPublicQueryStable(t *testing.T) {
	t.Parallel()

	first := canonicalPublicQuery(map[string]string{
		"search":        " Falcon Heavy ",
		"homepage_only": "true",
		"page":          "2",
	})
	second := canonicalPublicQuery(map[string]string{
		"page":          "2",
		"homepage_only": "true",
		"search":        "falcon heavy",
	})

	require.Equal(t, "homepage_only=true&page=2&search=falcon+heavy", first)
	require.Equal(t, first, second)
}

func TestPublicCacheStaleWhileRevalidate(t *testing.T) {
	now := time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC)
	manager := newPublicCacheManagerWithStore(cache.NewSWRStore(cache.SWRStoreOptions{
		Clock: func() time.Time { return now },
	}))
	opts := publicCacheOptions{
		key:            "launch-detail",
		policy:         publicDefaultPolicy,
		negativePolicy: publicNegativePolicy,
		shouldCache:    true,
	}

	var calls atomic.Int32
	refreshStarted := make(chan struct{})
	releaseRefresh := make(chan struct{})
	loader := func() (string, error) {
		call := calls.Add(1)
		if call == 1 {
			return "v1", nil
		}
		close(refreshStarted)
		<-releaseRefresh
		return "v2", nil
	}

	value, err := loadPublicCached(manager, opts, loader)
	require.NoError(t, err)
	require.Equal(t, "v1", value)

	now = now.Add(6 * time.Minute)

	value, err = loadPublicCached(manager, opts, loader)
	require.NoError(t, err)
	require.Equal(t, "v1", value)

	require.Eventually(t, func() bool {
		select {
		case <-refreshStarted:
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)

	close(releaseRefresh)

	require.Eventually(t, func() bool {
		value, err := loadPublicCached(manager, opts, loader)
		return err == nil && value == "v2"
	}, time.Second, 10*time.Millisecond)
	require.EqualValues(t, 2, calls.Load())
}

func TestPublicCacheHardExpiryForcesSynchronousReload(t *testing.T) {
	now := time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC)
	manager := newPublicCacheManagerWithStore(cache.NewSWRStore(cache.SWRStoreOptions{
		Clock: func() time.Time { return now },
	}))
	opts := publicCacheOptions{
		key:            "rocket-detail",
		policy:         publicDefaultPolicy,
		negativePolicy: publicNegativePolicy,
		shouldCache:    true,
	}

	var calls atomic.Int32
	reloadStarted := make(chan struct{})
	releaseReload := make(chan struct{})
	loader := func() (string, error) {
		call := calls.Add(1)
		if call == 1 {
			return "v1", nil
		}
		close(reloadStarted)
		<-releaseReload
		return "v2", nil
	}

	value, err := loadPublicCached(manager, opts, loader)
	require.NoError(t, err)
	require.Equal(t, "v1", value)

	now = now.Add(31 * time.Minute)

	done := make(chan struct{})
	var got string
	var gotErr error
	go func() {
		got, gotErr = loadPublicCached(manager, opts, loader)
		close(done)
	}()

	require.Eventually(t, func() bool {
		select {
		case <-reloadStarted:
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)

	select {
	case <-done:
		t.Fatal("hard-expired cache unexpectedly returned before reload finished")
	default:
	}

	close(releaseReload)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for hard-expired reload")
	}

	require.NoError(t, gotErr)
	require.Equal(t, "v2", got)
	require.EqualValues(t, 2, calls.Load())
}

func TestPublicCacheSingleflightOnMiss(t *testing.T) {
	manager := newPublicCacheManagerWithStore(cache.NewSWRStore(cache.SWRStoreOptions{}))
	opts := publicCacheOptions{
		key:            "company-list",
		policy:         publicDefaultPolicy,
		negativePolicy: publicNegativePolicy,
		shouldCache:    true,
	}

	var calls atomic.Int32
	release := make(chan struct{})
	loader := func() (string, error) {
		calls.Add(1)
		<-release
		return "value", nil
	}

	const workers = 8
	var wg sync.WaitGroup
	results := make([]string, workers)
	errs := make([]error, workers)
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(index int) {
			defer wg.Done()
			results[index], errs[index] = loadPublicCached(manager, opts, loader)
		}(i)
	}

	require.Eventually(t, func() bool {
		return calls.Load() == 1
	}, time.Second, 10*time.Millisecond)

	close(release)
	wg.Wait()

	require.EqualValues(t, 1, calls.Load())
	for i := 0; i < workers; i++ {
		require.NoError(t, errs[i])
		require.Equal(t, "value", results[i])
	}
}

func TestPublicCacheVersionBumpChangesKeys(t *testing.T) {
	service := &MainService{
		publicCache: newPublicCacheManagerWithStore(cache.NewSWRStore(cache.SWRStoreOptions{})),
	}

	initial := service.publicDetailCacheOptions(publicDomainLaunch, 42).key
	service.bumpPublicCacheDomains(publicDomainLaunch)
	bumped := service.publicDetailCacheOptions(publicDomainLaunch, 42).key

	require.Equal(t, "public:v1:launch:detail:ver=0:id=42", initial)
	require.Equal(t, "public:v1:launch:detail:ver=1:id=42", bumped)
}

func TestPublicCacheNegativeCache(t *testing.T) {
	manager := newPublicCacheManagerWithStore(cache.NewSWRStore(cache.SWRStoreOptions{}))
	opts := publicCacheOptions{
		key:            "missing-launch",
		policy:         publicDefaultPolicy,
		negativePolicy: publicNegativePolicy,
		shouldCache:    true,
	}

	var calls atomic.Int32
	loader := func() (string, error) {
		calls.Add(1)
		return "", mongo.ErrNoDocuments
	}

	_, err := loadPublicCached(manager, opts, loader)
	require.ErrorIs(t, err, mongo.ErrNoDocuments)

	_, err = loadPublicCached(manager, opts, loader)
	require.ErrorIs(t, err, mongo.ErrNoDocuments)
	require.EqualValues(t, 1, calls.Load())

	entry, _, ok := manager.store.Get(opts.key)
	require.True(t, ok)
	require.True(t, entry.Negative)
}

func TestPublicCacheHardExpiryFallsBackToStaleOnReloadError(t *testing.T) {
	now := time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC)
	manager := newPublicCacheManagerWithStore(cache.NewSWRStore(cache.SWRStoreOptions{
		Clock: func() time.Time { return now },
	}))
	opts := publicCacheOptions{
		key:            "launch-list",
		policy:         publicDefaultPolicy,
		negativePolicy: publicNegativePolicy,
		shouldCache:    true,
	}

	var calls atomic.Int32
	loader := func() (string, error) {
		call := calls.Add(1)
		if call == 1 {
			return "cached", nil
		}
		return "", errors.New("mongo timeout")
	}

	value, err := loadPublicCached(manager, opts, loader)
	require.NoError(t, err)
	require.Equal(t, "cached", value)

	now = now.Add(31 * time.Minute)

	value, err = loadPublicCached(manager, opts, loader)
	require.NoError(t, err)
	require.Equal(t, "cached", value)
	require.EqualValues(t, 2, calls.Load())
}
