package storage

import (
	"errors"
	"time"

	"golang.org/x/sync/syncmap"
)

// SyncMapStorage is a Redis-like storage model based on syncmap implementation
type SyncMapStorage struct {
	ttl time.Duration
	// I've chosen syncmap to avoid manual concurrency management (locking/unlocking mutexes)
	// and to get benefits of its inernal model (read non-only non-blocking access, synchronized write access)
	internalStorage *syncmap.Map
}

// InitSyncMapStorage ...
func InitSyncMapStorage(ttl time.Duration) *SyncMapStorage {
	newStorage := new(SyncMapStorage)
	newStorage.internalStorage = new(syncmap.Map)
	newStorage.ttl = ttl

	return newStorage
}

func (ls *SyncMapStorage) getTtl() time.Duration {
	return ls.ttl
}

// GetAllKeys ....
func (ls *SyncMapStorage) GetAllKeys() []string {
	// having no opportunity to get length of ls.internalStorage i have chosen 0 & 16 magic numbers
	keys := make([]string, 0, 16)
	ls.internalStorage.Range(func(key interface{}, value interface{}) bool {
		if notExpired(value.(*StorableWithMeta), ls) {
			keys = append(keys, key.(string)) // FIXME unsafe
		}
		return true
	})
	return keys
}

// GetValueByKey ....
func (ls *SyncMapStorage) GetValueByKey(key string) (*StorableWithMeta, bool) {
	if value, exists := ls.internalStorage.Load(key); exists {
		return filterExpired(value.(*StorableWithMeta), ls)
	}
	return nil, false
}

// DeleteValueByKey ...
func (ls *SyncMapStorage) DeleteValueByKey(key string) bool {
	ls.internalStorage.Delete(key)
	return true
}

// GetNestedValueByKeyAndIndex ....
func (ls *SyncMapStorage) GetNestedValueByKeyAndIndex(key string, index int) (*StorableWithMeta, bool, error) {
	if maybeSlice, exists := ls.internalStorage.Load(key); exists {
		if swm, yes := maybeSlice.(*StorableWithMeta); yes && notExpired(swm, ls) {
			slice, yes := swm.Entity.([]string) // TODO  looks ugly. consider another generic / polymorphic construction
			if yes {
				if index >= 0 && index < len(slice) {
					valWithMeta := enpackStorable(slice[index], swm)
					return valWithMeta, true, nil
				}
				return nil, false, errors.New("Index out of range")
			}
			return nil, false, errors.New("Stored value is not an array")
		}
	}
	return nil, false, nil
}

// GetNestedValueByKeyAndSubkey ...
func (ls *SyncMapStorage) GetNestedValueByKeyAndSubkey(key string, subKey string) (*StorableWithMeta, bool, error) {
	if maybeDict, exists := ls.internalStorage.Load(key); exists {
		if swm, yes := maybeDict.(*StorableWithMeta); yes && notExpired(swm, ls) {
			dict, yes := swm.Entity.(map[string]string) // TODO  looks ugly. consider another generic / polymorphic construction
			if yes {
				val, ok := dict[subKey]
				valWithMeta := enpackStorable(val, swm)
				return valWithMeta, ok, nil
			}
			return nil, false, errors.New("Stored value is not a dictionary")
		}
	}
	return nil, false, errors.New("Map not found")
}

// UpdateValueByKey ...
func (ls *SyncMapStorage) UpdateValueByKey(key string, newValue Storable) error {
	if _, exists := ls.internalStorage.Load(key); exists {
		// wrapping value into newStorableWithMeta ensures that LastWriteTime will be updated
		// and lifetime of the entity will be prolonged
		ls.internalStorage.Store(key, newStorableWithMeta(newValue))
		return nil
	}
	return errors.New("There is no entry with such key")
}

// AppendNewValue ...
func (ls *SyncMapStorage) AppendNewValue(key string, newValue Storable) error {
	if _, exists := ls.internalStorage.Load(key); !exists {
		ls.internalStorage.Store(key, newStorableWithMeta(newValue))
		return nil
	}
	return errors.New("Entry with such key already exists")
}
