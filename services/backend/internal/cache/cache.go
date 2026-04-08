package cache

// Cache defines the interface for a cache service
type Cache interface {
	// value must be a pointer,and json.marshal will be used for serialization
	Set(key string, value any) error
	// dest must be a pointer, and json.unmarshal will be used for deserialization
	Get(key string, dest any) error
	SetString(key, value string) error
	GetString(key string) (string, error)
	Delete(key string) error
}
