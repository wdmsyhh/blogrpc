package component

import "sync"

type safeMap struct {
	lock *sync.RWMutex
	sm   map[string]interface{}
}

func NewSafeMap() *safeMap {
	return &safeMap{
		lock: new(sync.RWMutex),
		sm:   make(map[string]interface{}),
	}
}

func (m *safeMap) Get(k string) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.sm[k]; ok {
		return val
	}
	return nil
}

func (m *safeMap) Set(k string, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.sm[k]; !ok {
		m.sm[k] = v
	} else if val != v {
		m.sm[k] = v
	} else {
		return false
	}
	return true
}

func (m *safeMap) Exists(k string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.sm[k]; ok {
		return true
	}
	return false
}

func (m *safeMap) Delete(k string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.sm, k)
}

func (m *safeMap) List() map[string]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.sm
}
