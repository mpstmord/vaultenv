package cache_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultenv/internal/cache"
)

func TestCache_SetAndGet(t *testing.T) {
	c := cache.New(5 * time.Minute)
	data := map[string]interface{}{"password": "s3cr3t"}
	c.Set("secret/app", data)

	got, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["password"] != "s3cr3t" {
		t.Errorf("unexpected value: %v", got["password"])
	}
}

func TestCache_Miss(t *testing.T) {
	c := cache.New(5 * time.Minute)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestCache_Expiry(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("secret/app", map[string]interface{}{"k": "v"})
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected entry to be expired")
	}
}

func TestCache_Delete(t *testing.T) {
	c := cache.New(5 * time.Minute)
	c.Set("secret/app", map[string]interface{}{"k": "v"})
	c.Delete("secret/app")
	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestCache_Purge(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("a", map[string]interface{}{"x": 1})
	c.Set("b", map[string]interface{}{"y": 2})
	time.Sleep(20 * time.Millisecond)
	c.Purge()
	_, okA := c.Get("a")
	_, okB := c.Get("b")
	if okA || okB {
		t.Fatal("expected all entries purged")
	}
}
