package cache

import (
	"sync"
	"testing"
	"time"
)

const TTL = 500 * time.Millisecond

func TestGetSet(t *testing.T) {
	type Test struct {
		key   string
		value string
	}

	tests := []Test{
		{"key11", "val11"},
		{"key12", "val12"},
		{"key13", "val13"},
	}

	cache := New()

	for _, test := range tests {
		cache.Set(test.key, test.value, TTL)
		if gotVal, gotOk := cache.Get(test.key); gotVal != test.value || gotOk != true {
			t.Errorf("cache.Get(%q) = %q, %t, want %q, %t", test.key, gotVal, gotOk, test.value, true)
		}
	}
}

func TestDelete(t *testing.T) {
	type Test struct {
		key   string
		value string
	}

	tests := []Test{
		{"key11", "val11"},
		{"key12", "val12"},
		{"key13", "val13"},
	}

	cache := New()

	for _, test := range tests {
		cache.Set(test.key, test.value, TTL)
		cache.Delete(test.key)
		if gotVal, gotOk := cache.Get(test.key); gotVal != "" || gotOk != false {
			t.Errorf("cache.Get(%q) = %q, %t, want %q, %t", test.key, gotVal, gotOk, "", false)
		}
	}
}

func TestTTL(t *testing.T) {
	type Test struct {
		key   string
		value string
	}

	tests := []Test{
		{"key11", "val11"},
		{"key12", "val12"},
		{"key13", "val13"},
	}

	cache := New()

	for _, test := range tests {
		cache.Set(test.key, test.value, TTL)
	}

	time.Sleep(TTL)

	for _, test := range tests {
		if gotVal, gotOk := cache.Get(test.key); gotVal != "" || gotOk != false {
			t.Errorf("cache.Get(%q) = %q, %t, want %q, %t", test.key, gotVal, gotOk, "", false)
		}
	}

}

func TestCache(t *testing.T) {
	type Test struct {
		key   string
		value string
	}

	tests1 := []Test{
		{"key11", "val11"},
		{"key12", "val12"},
		{"key13", "val13"},
	}

	tests2 := []Test{
		{"key21", "val21"},
		{"key22", "val22"},
		{"key23", "val23"},
	}

	cache := New()

	var wg sync.WaitGroup

	testCacheWithSet := func(tests []Test) {
		for _, test := range tests {
			cache.Set(test.key, test.value, TTL)
			if gotVal, gotOk := cache.Get(test.key); gotVal != test.value || gotOk != true {
				t.Errorf("cache.Get(%q) = %q, %t, want %q, %t", test.key, gotVal, gotOk, test.value, true)
			}
		}
	}

	testCacheWithDelete := func(tests []Test) {
		for _, test := range tests {
			cache.Set(test.key, test.value, TTL)
			cache.Delete(test.key)
			if gotVal, gotOk := cache.Get(test.key); gotVal != "" || gotOk != false {
				t.Errorf("cache.Get(%q) = %q, %t, want %q, %t", test.key, gotVal, gotOk, "", false)
			}
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		testCacheWithSet(tests1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		testCacheWithDelete(tests2)
	}()

	wg.Wait()
}
