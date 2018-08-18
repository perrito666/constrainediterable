package constrainediterable

import (
	"sort"
	"time"
)

type internalElement struct {
	value        interface{}
	creationTime time.Time
}

// NewMap returns a Map with the passed size and age limits.
func NewMap(size int, age time.Duration) *Map {
	return &Map{
		innerMap:   map[string]*internalElement{},
		sortedKeys: []string{},
		sizeLimit:  size,
		ageLimit:   age,
	}
}

// NewMapFromMap acts like `NewMap` but receives a `map[string]interface{}` to pre-populate
// the new map, eviction will be called just in case.
func NewMapFromMap(size int, age time.Duration, m map[string]interface{}) *Map {
	now := time.Now()
	mc := make(map[string]*internalElement, len(m))
	keys := make([]string, len(m))
	i := 0
	for k, v := range m {
		mc[k] = &internalElement{
			value:        v,
			creationTime: now,
		}
		keys[i] = k
		i++
	}
	rm := &Map{
		innerMap:   mc,
		sizeLimit:  size,
		ageLimit:   age,
		sortedKeys: keys,
	}
	rm.evict()
	return rm
}

// Map is a hashmap that allos limiting by size and age of its elements
type Map struct {
	// innerMap is where the data is actually stored.
	innerMap map[string]*internalElement
	// sizeLimit holds the amount of items to be held in the map, if it's set
	// stepping over the thresold will trigger eviction.
	sizeLimit int
	// ageLimit represents how old the elements on the map should be, when
	// eviction is triggered and this is set, all the elements older than
	// the specified age will be evicted.
	ageLimit time.Duration
	// sortedKeys contains a list of keys sorted by time, it should be assumed
	// dissordered.
	sortedKeys []string
}

// Len is the number of elements in the collection.
func (m *Map) Len() int {
	return len(m.innerMap)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (m *Map) Less(i, j int) bool {
	return m.innerMap[m.sortedKeys[i]].creationTime.Unix() >
		m.innerMap[m.sortedKeys[j]].creationTime.Unix()
}

// Swap swaps the elements with indexes i and j.
func (m *Map) Swap(i, j int) {
	m.sortedKeys[i], m.sortedKeys[j] = m.sortedKeys[j], m.sortedKeys[i]
}

// Add adds an element to the map.
func (m *Map) Add(k string, v interface{}) {
	now := time.Now()
	m.innerMap[k] = &internalElement{
		value:        v,
		creationTime: now,
	}
	m.sortedKeys = append(m.sortedKeys, k)
	if m.sizeLimit == 0 {
		return
	}

	if len(m.innerMap) > m.sizeLimit {
		m.evict()
	}
}

// Get returns the value of key k if it exists and bool indicating existence
func (m *Map) Get(k string) (interface{}, bool) {
	ie, ok := m.innerMap[k]
	if ok {
		return ie.value, ok
	}
	return nil, ok
}

func (m *Map) evictOld() {
	now := time.Now()
	deletedUpTo := -1
	for i := len(m.sortedKeys) - 1; i >= 0; i-- {
		k := m.sortedKeys[i]
		if now.Sub(m.innerMap[k].creationTime) > m.ageLimit {
			delete(m.innerMap, k)
			deletedUpTo = i
		}
	}
	if deletedUpTo >= 0 {
		m.sortedKeys = m.sortedKeys[:deletedUpTo]
	}
}

func (m *Map) evictN(n int) {
	if n < 1 {
		return
	}

	evictIndex := len(m.sortedKeys) - n
	for i := len(m.sortedKeys) - 1; i >= evictIndex; i-- {
		k := m.sortedKeys[i]
		delete(m.innerMap, k)
	}
	m.sortedKeys = m.sortedKeys[:evictIndex]
}

func (m *Map) evict() {
	sort.Sort(m)
	if m.ageLimit != 0 {
		m.evictOld()
	}

	if len(m.innerMap) > m.sizeLimit {
		m.evictN(len(m.innerMap) - m.sizeLimit)
	}
}
