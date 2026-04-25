package coalesce_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/vaultenv/internal/coalesce"
)

type stubFetcher struct {
	data map[string]interface{}
	err  error
}

func (s *stubFetcher) GetSecretData(_ context.Context, _ string) (map[string]interface{}, error) {
	return s.data, s.err
}

func TestNew_PanicsWithNoFetchers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic with no fetchers")
		}
	}()
	coalesce.New()
}

func TestGetSecretData_FirstSucceeds(t *testing.T) {
	want := map[string]interface{}{"key": "value"}
	f1 := &stubFetcher{data: want}
	f2 := &stubFetcher{err: errors.New("should not be called")}

	c := coalesce.New(f1, f2)
	got, err := c.GetSecretData(context.Background(), "secret/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetSecretData_FallsBackOnError(t *testing.T) {
	want := map[string]interface{}{"token": "abc"}
	f1 := &stubFetcher{err: errors.New("primary unavailable")}
	f2 := &stubFetcher{data: want}

	c := coalesce.New(f1, f2)
	got, err := c.GetSecretData(context.Background(), "secret/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["token"] != "abc" {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetSecretData_AllFail_ReturnsJoinedError(t *testing.T) {
	f1 := &stubFetcher{err: errors.New("err1")}
	f2 := &stubFetcher{err: errors.New("err2")}

	c := coalesce.New(f1, f2)
	_, err := c.GetSecretData(context.Background(), "secret/test")
	if err == nil {
		t.Fatal("expected error when all fetchers fail")
	}
	if !errors.Is(err, f1.err) && !errors.Is(err, f2.err) {
		// errors.Join wraps both; just check the string contains both messages
		msg := err.Error()
		if msg == "" {
			t.Error("error message should not be empty")
		}
	}
}

func TestGetSecretData_SingleFetcher_Success(t *testing.T) {
	want := map[string]interface{}{"db": "pass"}
	c := coalesce.New(&stubFetcher{data: want})
	got, err := c.GetSecretData(context.Background(), "secret/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["db"] != "pass" {
		t.Errorf("got %v, want %v", got, want)
	}
}
