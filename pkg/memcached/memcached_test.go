package memcached

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"
)

const (
	host         = "localhost"
	port         = 11211
	maxIdleConns = 10
	maxOpenConns = 10
	TTL          = 1
)

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

	client, err := Connect(
		host,
		WithPort(port),
		WithMaxIdleConns(maxIdleConns),
		WithMaxOpenConns(maxOpenConns),
	)
	if err != nil {
		t.Fatalf("unable to connect to memcached server: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	for _, test := range tests {
		if err = client.Set(ctx, test.key, test.value, TTL); err != nil {
			t.Fatalf("unable to set key: %q value: %q : %v", test.key, test.value, err)
		}
		if gotVal, gotErr := client.Get(ctx, test.key); gotVal != test.value || gotErr == ErrNotFound {
			t.Fatalf("client.Get(%q) = %q, %q, want %q, %v", test.key, gotVal, gotErr, test.value, nil)
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

	client, err := Connect(
		host,
		WithPort(port),
		WithMaxIdleConns(maxIdleConns),
		WithMaxOpenConns(maxOpenConns),
	)
	if err != nil {
		t.Fatalf("unable to connect to memcached server: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	for _, test := range tests {
		if err = client.Set(ctx, test.key, test.value, TTL); err != nil {
			t.Fatalf("unable to set key: %q value: %q : %v", test.key, test.value, err)
		}
		if err = client.Delete(ctx, test.key); err != nil {
			t.Fatalf("unable to delete key: %q : %v", test.key, err)
		}
		if gotVal, gotErr := client.Get(ctx, test.key); gotVal != "" || gotErr != ErrNotFound {
			t.Fatalf("client.Get(%q) = %q, %v, want %q, %v", test.key, gotVal, gotErr, "", ErrNotFound)
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

	client, err := Connect(
		host,
		WithPort(port),
		WithMaxIdleConns(maxIdleConns),
		WithMaxOpenConns(maxOpenConns),
	)
	if err != nil {
		log.Fatalf("unable to connect to memcached server: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	for _, test := range tests {
		if err = client.Set(ctx, test.key, test.value, TTL); err != nil {
			t.Fatalf("unable to set key: %q value: %q : %v", test.key, test.value, err)
		}
	}

	time.Sleep(TTL * time.Second)

	for _, test := range tests {
		if gotVal, gotErr := client.Get(ctx, test.key); gotVal != "" || gotErr != ErrNotFound {
			t.Fatalf("client.Get(%q) = %q, %v, want %q, %v", test.key, gotVal, gotErr, "", ErrNotFound)
		}
	}

}

func TestMemcached(t *testing.T) {
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

	client, err := Connect(
		host,
		WithPort(port),
		WithMaxIdleConns(maxIdleConns),
		WithMaxOpenConns(maxOpenConns),
	)
	if err != nil {
		t.Fatalf("unable to connect to memcached server: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	var wg sync.WaitGroup

	testCacheWithSet := func(tests []Test) {
		for _, test := range tests {
			if err = client.Set(ctx, test.key, test.value, TTL); err != nil {
				t.Fatalf("unable to set key: %q value: %q : %v", test.key, test.value, err)
			}
			if gotVal, gotErr := client.Get(ctx, test.key); gotVal != test.value || gotErr == ErrNotFound {
				t.Fatalf("client.Get(%q) = %q, %q, want %q, %v", test.key, gotVal, gotErr, test.value, nil)
			}
		}
	}

	testCacheWithDelete := func(tests []Test) {
		for _, test := range tests {
			if err = client.Set(ctx, test.key, test.value, TTL); err != nil {
				t.Fatalf("unable to set key: %q value: %q : %v", test.key, test.value, err)
			}
			if err = client.Delete(ctx, test.key); err != nil {
				t.Fatalf("unable to delete key: %q : %v", test.key, err)
			}
			if gotVal, gotErr := client.Get(ctx, test.key); gotVal != "" || gotErr != ErrNotFound {
				t.Fatalf("client.Get(%q) = %q, %v, want %q, %v", test.key, gotVal, gotErr, "", ErrNotFound)
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
