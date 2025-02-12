package utils

import (
	"sync"
	"testing"
)

func TestSafeMap_SetGet(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("apple", 10)

	if value, exists := sm.Get("apple"); !exists || value != 10 {
		t.Errorf("Expected apple=10, got %v, exists: %v", value, exists)
	}

	if _, exists := sm.Get("banana"); exists {
		t.Errorf("Expected banana to not exist")
	}
}

func TestSafeMap_Delete(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("apple", 10)
	sm.Delete("apple")

	if _, exists := sm.Get("apple"); exists {
		t.Errorf("Expected apple to be deleted")
	}
}

func TestSafeMap_Keys(t *testing.T) {
	sm := NewSafeMap[string, int]()
	sm.Set("apple", 10)
	sm.Set("banana", 20)

	keys := sm.Keys()
	expectedKeys := map[string]bool{"apple": true, "banana": true}

	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}

	for _, key := range keys {
		if !expectedKeys[key] {
			t.Errorf("Unexpected key found: %v", key)
		}
	}
}

func TestSafeMap_ConcurrentAccess(t *testing.T) {
	sm := NewSafeMap[int, int]()
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sm.Set(i, i*10)
		}(i)
	}

	wg.Wait()

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			defer wg.Done()
			if value, exists := sm.Get(i); !exists || value != i*10 {
				t.Errorf("Mismatch: key=%d, expected=%d, got=%d", i, i*10, value)
			}
		}(i)
	}
	wg.Wait()
}
