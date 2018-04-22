package storage

import "time"

// StorableWithMeta ...
type StorableWithMeta struct {
	LastWriteTime time.Time
	Entity        Storable
}

func newStorableWithMeta(entity Storable) *StorableWithMeta {
	s := new(StorableWithMeta)
	s.Entity = entity
	s.LastWriteTime = time.Now()
	return s
}

// enpackStorable call wraps internal element (cell of slice or value extracted from map)
// to the StorableWithMeta with LastWriteTime nested from top-level storable entity
func enpackStorable(entity Storable, ref *StorableWithMeta) *StorableWithMeta {
	s := new(StorableWithMeta)
	s.Entity = entity
	s.LastWriteTime = ref.LastWriteTime
	return s
}

// Storable is the storable primitive (string, slice of string or map string -> string)
// I had to use generic-like interface{} and runtime checks because considered it as
// a solution demanding less coding
// TODO other alternative should be considered
type Storable interface{}

// Storage is a basic interface for any kind of redis-like registry implementation
type Storage interface {
	getTtl() time.Duration

	GetAllKeys() []string

	GetValueByKey(key string) (*StorableWithMeta, bool)

	DeleteValueByKey(key string) bool

	GetNestedValueByKeyAndIndex(key string, index int) (*StorableWithMeta, bool, error)

	GetNestedValueByKeyAndSubkey(key string, subKey string) (*StorableWithMeta, bool, error)

	UpdateValueByKey(key string, newValue Storable) error

	AppendNewValue(key string, newValue Storable) error
}

// TODO refactor : pass ttl
func filterExpired(entity *StorableWithMeta, ls Storage) (*StorableWithMeta, bool) {
	if notExpired(entity, ls) {
		return entity, true
	}
	return nil, false
}

// TODO refactor : pass ttl
func notExpired(entity *StorableWithMeta, ls Storage) bool {
	now := time.Now()
	return !entity.LastWriteTime.Add(ls.getTtl()).Before(now)
}
