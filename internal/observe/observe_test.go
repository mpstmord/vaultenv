package observe_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/vaultenv/vaultenv/internal/observe"
)

// stubFetcher is a minimal Fetcher for testing.
type stubFetcher struct {
	data map[string]interface{}
	err  error
}

func (s *stubFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	return s.data, s.err
}

func TestNew_PanicsOnNilUpstream(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil upstream")
		}
	}()
	observe.New(nil)
}

func TestGetSecretData_CallsObserver(t *testing.T) {
	var received observe.Event
	obs := func(_ context.Context, e observe.Event) { received = e }

	upstream := &stubFetcher{data: map[string]interface{}{"key": "val"}}
	f := observe.New(upstream, obs)

	_, err := f.GetSecretData(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Path != "secret/foo" {
		t.Errorf("expected path %q, got %q", "secret/foo", received.Path)
	}
	if received.Err != nil {
		t.Errorf("expected nil error in event, got %v", received.Err)
	}
}

func TestGetSecretData_ObserverReceivesError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	var received observe.Event
	obs := func(_ context.Context, e observe.Event) { received = e }

	upstream := &stubFetcher{err: sentinel}
	f := observe.New(upstream, obs)

	_, _ = f.GetSecretData(context.Background(), "secret/bar")
	if !errors.Is(received.Err, sentinel) {
		t.Errorf("expected sentinel error in event, got %v", received.Err)
	}
}

func TestMetricsObserver_Counts(t *testing.T) {
	counters := &observe.Counters{}
	upstream := &stubFetcher{data: map[string]interface{}{"x": "1"}}
	f := observe.New(upstream, observe.MetricsObserver(counters))

	for i := 0; i < 3; i++ {
		_, _ = f.GetSecretData(context.Background(), "secret/x")
	}
	if got := counters.Total.Load(); got != 3 {
		t.Errorf("expected total=3, got %d", got)
	}
	if got := counters.Errors.Load(); got != 0 {
		t.Errorf("expected errors=0, got %d", got)
	}
}

func TestMetricsObserver_PanicsOnNilCounters(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil Counters")
		}
	}()
	observe.MetricsObserver(nil)
}

func TestLogObserver_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	upstream := &stubFetcher{data: map[string]interface{}{"k": "v"}}
	f := observe.New(upstream, observe.LogObserver(&buf))

	_, _ = f.GetSecretData(context.Background(), "secret/log")

	if buf.Len() == 0 {
		t.Fatal("expected log output, got nothing")
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"path":"secret/log"`)) {
		t.Errorf("log line missing path field: %s", buf.String())
	}
}

func TestGetSecretData_MultipleObservers_AllCalled(t *testing.T) {
	called := make([]bool, 3)
	var obs []observe.Observer
	for i := range called {
		i := i
		obs = append(obs, func(_ context.Context, _ observe.Event) { called[i] = true })
	}
	upstream := &stubFetcher{data: map[string]interface{}{}}
	f := observe.New(upstream, obs...)
	_, _ = f.GetSecretData(context.Background(), "secret/multi")
	for i, c := range called {
		if !c {
			t.Errorf("observer %d was not called", i)
		}
	}
}
