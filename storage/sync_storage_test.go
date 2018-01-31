package storage

import (
	"testing"
  "time"
)

var testTTL = time.Second * 2
var testStorage Storage = InitSyncStorage(testTTL)

func TestEmpty(t *testing.T) {
	if len(testStorage.GetAllKeys()) != 0 {
		t.Error("Test storage is not empty after init")
	}
}

func TestBasicCrud(t *testing.T) {
	key, value := "zxcvf", "London is the capital of ..."
  newValue := "Mordor is the land of..."

	testStorage.AppendNewValue(key, value)
	if storedVal, ok := testStorage.GetValueByKey(key); !ok || storedVal.Entity != value {
		t.Error("Test storage does not contain the test value just appended")
	}
  // todo error when appending the save value
  testStorage.UpdateValueByKey(key, newValue)
  if storedVal, ok := testStorage.GetValueByKey(key); !ok || storedVal.Entity != newValue {
		t.Error("Test storage does not contain the new test value just updated")
	}

  testStorage.DeleteValueByKey(key)
  if _, ok := testStorage.GetValueByKey(key); ok {
    t.Error("Test storage contains the new value just deleted")
  }
}

func TestNestedArrayOps(t *testing.T) {
  key, value := "hjkl", []string{"Alpha", "Bravo", "Charlie"}
  ind := 1
  testStorage.AppendNewValue(key, value)

  resp, exists, err := testStorage.GetNestedValueByKeyAndIndex(key, 1)
  if !exists || err != nil || resp.Entity != value[ind] {
    t.Error("Can not get acces to the nested array item by its index")
  }

  resp, exists, err = testStorage.GetNestedValueByKeyAndIndex(key, 10)
  if err == nil || exists {
    t.Error("Outbounding index does not lead to error")
  }

  resp, exists, err = testStorage.GetNestedValueByKeyAndIndex(key, -13)
  if err == nil || exists {
    t.Error("Negative index does not lead to error")
  }
}

func TestNestedDictionaryOps(t *testing.T) {
  key, subkey := "qwegs", "the_second"
  value := map[string]string{
    "the_first": "Nicolas",
    "the_second": "Francois",
  }
  testStorage.AppendNewValue(key, value)

  resp, exists, err := testStorage.GetNestedValueByKeyAndSubkey(key, subkey)

  if !exists || err != nil || resp.Entity != value[subkey] {
    t.Error("Can not get acces to the nested dictionary item by its subkey")
  }
}

func TestExpiration(t *testing.T) {
  key, value := "zxcvf", "London is the capital of ..."

	testStorage.AppendNewValue(key, value)
	if storedVal, ok := testStorage.GetValueByKey(key); !ok || storedVal.Entity != value {
		t.Error("Test storage does not contain the test value just appended")
	}
  // TODO get TTL from config
  time.Sleep(testTTL * 2)
  if _, ok := testStorage.GetValueByKey(key); ok {
    t.Error("Test storage still contains the test value that should be already expired")
  }
}
