package tests

import (
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestGetAllConfig(t *testing.T) {
	// Create a new Consul client
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}

	// Get the KV client
	kv := client.KV()

	// Set up some test key-value pairs
	keys := []string{"test-key-1", "test-key-2", "test-key-3"}
	values := []string{"test-value-1", "test-value-2", "test-value-3"}
	for i, key := range keys {
		_, err := kv.Put(&api.KVPair{Key: key, Value: []byte(values[i])}, nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Get all configuration
	pairs, _, err := kv.List("", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Check that all key-value pairs are returned
	if len(pairs) != len(keys) {
		t.Errorf("expected %d key-value pairs, got %d", len(keys), len(pairs))
	}

	// Check that the values are correct
	for _, pair := range pairs {
		for i, key := range keys {
			if pair.Key == key {
				if string(pair.Value) != values[i] {
					t.Errorf("expected value %q for key %q, got %q", values[i], key, string(pair.Value))
				}
			}
		}
	}
}
