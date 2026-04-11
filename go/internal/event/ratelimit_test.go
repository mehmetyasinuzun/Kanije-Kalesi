package event

import (
	"testing"
	"time"
)

func TestTokenBucketCapacityAndRefill(t *testing.T) {
	tb := newTokenBucket(2, 40*time.Millisecond)
	key := "login_failed"

	if !tb.Allow(key) {
		t.Fatalf("ilk event kabul edilmeliydi")
	}
	if !tb.Allow(key) {
		t.Fatalf("ikinci event kabul edilmeliydi")
	}
	if tb.Allow(key) {
		t.Fatalf("kapasite dolunca event reddedilmeliydi")
	}

	time.Sleep(90 * time.Millisecond)

	if !tb.Allow(key) {
		t.Fatalf("refill sonrasi event tekrar kabul edilmeliydi")
	}
}

func TestDedupCacheWindow(t *testing.T) {
	dc := newDedupCache(50 * time.Millisecond)
	key := "usb_inserted:alice:10.0.0.5"

	if dc.IsDuplicate(key) {
		t.Fatalf("ilk gorulmede duplicate olmamali")
	}
	if !dc.IsDuplicate(key) {
		t.Fatalf("pencere icindeki ikinci gorulmede duplicate olmali")
	}

	time.Sleep(70 * time.Millisecond)

	if dc.IsDuplicate(key) {
		t.Fatalf("pencere dolduktan sonra duplicate olmamali")
	}
}
