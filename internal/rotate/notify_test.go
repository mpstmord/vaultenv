package rotate

import (
	"errors"
	"log"
	"os"
	"testing"
)

func TestLogHandler_NoError(t *testing.T) {
	logger := log.New(os.Stdout, "", 0)
	h := LogHandler(logger)
	if err := h("secret/app", map[string]interface{}{"k": "v"}); err != nil {
		t.Fatal(err)
	}
}

func TestChainHandlers_AllOK(t *testing.T) {
	var order []int
	h1 := Handler(func(_ string, _ map[string]interface{}) error { order = append(order, 1); return nil })
	h2 := Handler(func(_ string, _ map[string]interface{}) error { order = append(order, 2); return nil })
	chain := ChainHandlers(h1, h2)
	if err := chain("p", nil); err != nil {
		t.Fatal(err)
	}
	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("unexpected order: %v", order)
	}
}

func TestChainHandlers_StopsOnError(t *testing.T) {
	var second bool
	h1 := Handler(func(_ string, _ map[string]interface{}) error { return errors.New("fail") })
	h2 := Handler(func(_ string, _ map[string]interface{}) error { second = true; return nil })
	chain := ChainHandlers(h1, h2)
	if err := chain("p", nil); err == nil {
		t.Fatal("expected error")
	}
	if second {
		t.Fatal("second handler should not have run")
	}
}

func TestExecHandler_InvalidCommand(t *testing.T) {
	h := ExecHandler("exit 1")
	if err := h("p", nil); err == nil {
		t.Fatal("expected error from exit 1")
	}
}
