package core

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vamosdalian/launchdate-backend/internal/cache"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/singleflight"
)

const (
	publicCachePrefix          = "public:v1"
	publicCacheMaxEntries      = 2048
	publicDomainLaunch         = "launch"
	publicDomainRocket         = "rocket"
	publicDomainCompany        = "company"
	publicDomainLaunchBase     = "launch_base"
	publicDomainPageBackground = "page_background"
)

var (
	publicDefaultPolicy = cache.Policy{
		SoftTTL: 5 * time.Minute,
		HardTTL: 30 * time.Minute,
	}
	publicNegativePolicy = cache.Policy{
		SoftTTL: 30 * time.Second,
		HardTTL: 2 * time.Minute,
	}
	publicPageBackgroundPolicy = cache.Policy{
		SoftTTL: 30 * time.Minute,
		HardTTL: 2 * time.Hour,
	}
)

type publicCacheManager struct {
	store *cache.SWRStore
	loads singleflight.Group
}

type publicCacheOptions struct {
	key            string
	policy         cache.Policy
	negativePolicy cache.Policy
	shouldCache    bool
}

func newPublicCacheManager() *publicCacheManager {
	return &publicCacheManager{
		store: cache.NewSWRStore(cache.SWRStoreOptions{
			MaxEntries: publicCacheMaxEntries,
		}),
	}
}

func newPublicCacheManagerWithStore(store *cache.SWRStore) *publicCacheManager {
	return &publicCacheManager{store: store}
}

func loadPublicCached[T any](manager *publicCacheManager, opts publicCacheOptions, loader func() (T, error)) (T, error) {
	var zero T
	if manager == nil || manager.store == nil || !opts.shouldCache {
		return loader()
	}

	entry, state, ok := manager.store.Get(opts.key)
	if ok {
		switch state {
		case cache.EntryStateFresh:
			return decodePublicCacheEntry[T](manager, opts.key, entry)
		case cache.EntryStateStale:
			refreshStalePublicCache(manager, opts, loader)
			return decodePublicCacheEntry[T](manager, opts.key, entry)
		case cache.EntryStateExpired:
			value, err := loadPublicCachedSync(manager, opts, loader)
			if err == nil {
				return value, nil
			}

			fallback, fallbackErr := decodePublicCacheEntry[T](manager, opts.key, entry)
			if fallbackErr == nil || errors.Is(fallbackErr, mongo.ErrNoDocuments) {
				logrus.WithError(err).Warnf("serving stale public cache for key %s after refresh failure", opts.key)
				return fallback, fallbackErr
			}

			return zero, err
		}
	}

	return loadPublicCachedSync(manager, opts, loader)
}

func loadPublicCachedSync[T any](manager *publicCacheManager, opts publicCacheOptions, loader func() (T, error)) (T, error) {
	var zero T
	result, err, _ := manager.loads.Do(opts.key, func() (interface{}, error) {
		value, loadErr := loader()
		if loadErr != nil {
			if errors.Is(loadErr, mongo.ErrNoDocuments) {
				manager.store.Set(opts.key, nil, true, opts.negativePolicy)
			}
			return nil, loadErr
		}

		payload, marshalErr := json.Marshal(value)
		if marshalErr != nil {
			return nil, marshalErr
		}

		manager.store.Set(opts.key, payload, false, opts.policy)
		return payload, nil
	})
	if err != nil {
		return zero, err
	}

	payload, ok := result.([]byte)
	if !ok {
		return zero, fmt.Errorf("unexpected public cache load result type %T", result)
	}

	var value T
	if err := json.Unmarshal(payload, &value); err != nil {
		manager.store.Delete(opts.key)
		return zero, err
	}

	return value, nil
}

func decodePublicCacheEntry[T any](manager *publicCacheManager, key string, entry cache.Item) (T, error) {
	var zero T
	if entry.Negative {
		return zero, mongo.ErrNoDocuments
	}

	var value T
	if err := json.Unmarshal(entry.Payload, &value); err != nil {
		manager.store.Delete(key)
		return zero, err
	}

	return value, nil
}

func refreshStalePublicCache[T any](manager *publicCacheManager, opts publicCacheOptions, loader func() (T, error)) {
	if !manager.store.TryStartRefresh(opts.key) {
		return
	}

	go func() {
		defer manager.store.FinishRefresh(opts.key)

		if _, err := loadPublicCachedSync(manager, opts, loader); err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			logrus.WithError(err).Warnf("failed to refresh stale public cache for key %s", opts.key)
		}
	}()
}

func canonicalPublicQuery(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	normalized := make(map[string]string, len(params))
	keys := make([]string, 0, len(params))
	for key, rawValue := range params {
		value := strings.TrimSpace(rawValue)
		if value == "" {
			continue
		}
		if key == "search" {
			value = strings.ToLower(value)
		}
		normalized[key] = value
		keys = append(keys, key)
	}

	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(normalized[key])))
	}

	return strings.Join(parts, "&")
}

func hashPublicQuery(query string) string {
	sum := sha1.Sum([]byte(query))
	return hex.EncodeToString(sum[:])
}

func detailPublicCacheKey(version uint64, domain string, id int64) string {
	return fmt.Sprintf("%s:%s:detail:ver=%d:id=%d", publicCachePrefix, domain, version, id)
}

func listPublicCacheKey(version uint64, domain string, query map[string]string) string {
	return fmt.Sprintf("%s:%s:list:ver=%d:q=%s", publicCachePrefix, domain, version, hashPublicQuery(canonicalPublicQuery(query)))
}

func shouldCachePublicList(page int, search string) bool {
	if strings.TrimSpace(search) == "" {
		return true
	}
	return page <= 3
}

func (m *MainService) publicCacheVersion(domain string) uint64 {
	if m.publicCache == nil || m.publicCache.store == nil {
		return 0
	}
	return m.publicCache.store.NamespaceVersion(domain)
}

func (m *MainService) bumpPublicCacheDomains(domains ...string) {
	if m.publicCache == nil || m.publicCache.store == nil {
		return
	}
	for _, domain := range domains {
		if strings.TrimSpace(domain) == "" {
			continue
		}
		m.publicCache.store.BumpNamespace(domain)
	}
}

func (m *MainService) publicDetailCacheOptions(domain string, id int64) publicCacheOptions {
	return publicCacheOptions{
		key:            detailPublicCacheKey(m.publicCacheVersion(domain), domain, id),
		policy:         publicDefaultPolicy,
		negativePolicy: publicNegativePolicy,
		shouldCache:    true,
	}
}

func (m *MainService) publicListCacheOptions(domain string, page int, search string, extra map[string]string) publicCacheOptions {
	query := map[string]string{
		"page":   fmt.Sprintf("%d", page),
		"search": search,
	}
	for key, value := range extra {
		query[key] = value
	}

	return publicCacheOptions{
		key:            listPublicCacheKey(m.publicCacheVersion(domain), domain, query),
		policy:         publicDefaultPolicy,
		negativePolicy: publicNegativePolicy,
		shouldCache:    shouldCachePublicList(page, search),
	}
}

func (m *MainService) pageBackgroundCacheOptions() publicCacheOptions {
	return publicCacheOptions{
		key:            listPublicCacheKey(m.publicCacheVersion(publicDomainPageBackground), publicDomainPageBackground, nil),
		policy:         publicPageBackgroundPolicy,
		negativePolicy: publicNegativePolicy,
		shouldCache:    true,
	}
}

func (m *MainService) InvalidatePublicCacheForSync(syncType string) {
	switch syncType {
	case "launch", "update":
		m.bumpPublicCacheDomains(publicDomainLaunch, publicDomainRocket, publicDomainCompany, publicDomainLaunchBase)
	case "agency":
		m.bumpPublicCacheDomains(publicDomainCompany, publicDomainRocket, publicDomainLaunch)
	case "launcher":
		m.bumpPublicCacheDomains(publicDomainRocket, publicDomainLaunch, publicDomainCompany)
	case "location", "pad":
		m.bumpPublicCacheDomains(publicDomainLaunchBase, publicDomainLaunch)
	}
}
