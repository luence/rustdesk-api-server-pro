package errcode

import (
	"errors"
	"strings"
	"testing"
)

// TestEnsurePreservesCodedError 验证已有错误码不会被重复包装。
func TestEnsurePreservesCodedError(t *testing.T) {
	err := New(ERR1001.Code, ERR1001.Message)
	if Ensure(err).Error() != err.Error() {
		t.Fatalf("coded error was changed: %s", Ensure(err))
	}
}

// TestEnsureWrapsRawError 验证裸错误对外返回前会补齐通用错误码。
func TestEnsureWrapsRawError(t *testing.T) {
	err := Ensure(errors.New("database unavailable"))
	if !strings.HasPrefix(err.Error(), ERRB010.Code+":") {
		t.Fatalf("raw error has no code: %s", err)
	}
}
