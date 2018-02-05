package storage

import (
	"errors"
	"log"
	"time"

	"golang.org/x/sync/syncmap"
)

// SyncStorage is a Redis-like storage model based on syncmap implementation
type SyncStorage struct {
	ttl time.Duration
	// I've chosen syncmap to avoid manual concurrency management (locking/unlocking mutexes)
	// and to get benefits of its inernal model (read non-only non-blocking access, synchronized write access)
	internalStorage *syncmap.Map
}

// InitSyncStorage ...
func InitSyncStorage(ttl time.Duration) *SyncStorage {
	newStorage := new(SyncStorage)
	newStorage.internalStorage = new(syncmap.Map)
	newStorage.ttl = ttl
	expiratorTicker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-expiratorTicker.C:
				newStorage.expireObsolete()
			}
			// TODO think how to stop it if needed
		}
	}()

	return newStorage
}

// GetAllKeys ....
func (ls *SyncStorage) GetAllKeys() []string {
	// having no opportunity to get length of ls.internalStorage i have chosen 0 & 16 magic numbers
	keys := make([]string, 0, 16)
	ls.internalStorage.Range(func(key interface{}, value interface{}) bool {
		if ls.notExpired(value.(*StorableWithMeta)) {
			keys = append(keys, key.(string)) // FIXME unsafe
		}
		return true
	})
	return keys
}

// GetValueByKey ....
func (ls *SyncStorage) GetValueByKey(key string) (*StorableWithMeta, bool) {
	if value, exists := ls.internalStorage.Load(key); exists {
		return ls.filterExpired(value.(*StorableWithMeta))
	}
	return nil, false
}

// DeleteValueByKey ...
func (ls *SyncStorage) DeleteValueByKey(key string) bool {
	ls.internalStorage.Delete(key)
	return true
}

// GetNestedValueByKeyAndIndex ....
func (ls *SyncStorage) GetNestedValueByKeyAndIndex(key string, index int) (*StorableWithMeta, bool, error) {
	if maybeSlice, exists := ls.internalStorage.Load(key); exists {
		if swm, yes := maybeSlice.(*StorableWithMeta); yes && ls.notExpired(swm) {
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
func (ls *SyncStorage) GetNestedValueByKeyAndSubkey(key string, subKey string) (*StorableWithMeta, bool, error) {
	if maybeDict, exists := ls.internalStorage.Load(key); exists {
		if swm, yes := maybeDict.(*StorableWithMeta); yes && ls.notExpired(swm) {
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
func (ls *SyncStorage) UpdateValueByKey(key string, newValue Storable) error {
	if _, exists := ls.internalStorage.Load(key); exists {
		// wrapping value into newStorableWithMeta ensures that LastWriteTime will be updated
		// and lifetime of the entity will be prolonged
		ls.internalStorage.Store(key, newStorableWithMeta(newValue))
		return nil
	}
	return errors.New("There is no entry with such key")
}

// AppendNewValue ...
func (ls *SyncStorage) AppendNewValue(key string, newValue Storable) error {
	if _, exists := ls.internalStorage.Load(key); !exists {
		ls.internalStorage.Store(key, newStorableWithMeta(newValue))
		return nil
	}
	return errors.New("Entry with such key already exists")
}

func (ls *SyncStorage) expireObsolete() {
	now := time.Now()
	toExpire := make([]string, 0)
	ls.internalStorage.Range(func(key interface{}, value interface{}) bool {
		if value.(*StorableWithMeta).LastWriteTime.Add(ls.ttl).Before(now) {
			toExpire = append(toExpire, key.(string))
		} // FIXME unsafe
		// TODO unsafe
		return true
	})

	for _, key := range toExpire {
		log.Println("Expiring enity with key: " + key)
		ls.internalStorage.Delete(key)
	}
}

func (ls *SyncStorage) filterExpired(entity *StorableWithMeta) (*StorableWithMeta, bool) {
	if ls.notExpired(entity) {
		return entity, true
	}
	return nil, false
}

func (ls *SyncStorage) notExpired(entity *StorableWithMeta) bool {
	now := time.Now()
	return !entity.LastWriteTime.Add(ls.ttl).Before(now)
}
