package loggable

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/copier"
	"reflect"
	"sync"
)

type identityManager struct {
	sync.Mutex
	store map[string]interface{}
}

func newIdentityManager() *identityManager {
	return &identityManager{
		store: map[string]interface{}{},
	}
}

func (im *identityManager) save(value, pk interface{}) {
	im.Lock()
	defer im.Unlock()

	t := reflect.TypeOf(value)
	newValue := reflect.New(t).Interface()
	err := copier.Copy(&newValue, value)
	if err != nil {
		panic(err)
	}

	k := genIdentityHash(value, pk)
	im.store[k] = newValue
}

func (im identityManager) get(value, pk interface{}) interface{} {
	im.Lock()
	defer im.Unlock()

	k := genIdentityHash(value, pk)
	value, ok := im.store[k]
	if !ok {
		return nil
	}

	return value
}

func (im *identityManager) diff(value, pk interface{}) UpdateDiff {
	old := im.get(value, pk)
	if old == nil {
		return nil
	}

	return computeDiff(old, value)
}

func genIdentityHash(value, pk interface{}) string {
	t := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(value)).Interface())
	key := fmt.Sprintf("%v_%s", pk, t.Name())
	b := md5.Sum([]byte(key))

	return hex.EncodeToString(b[:])
}
