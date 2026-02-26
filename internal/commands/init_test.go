package commands

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"testing"
)

type failReader struct {
	called bool
}

func (r *failReader) Read(p []byte) (int, error) {
	r.called = true
	return 0, errors.New("read called")
}

func TestShouldOverwriteConfigReinitSkipsPrompt(t *testing.T) {
	reader := &failReader{}
	ok, err := shouldOverwriteConfig(bufio.NewReader(reader), io.Discard, true, true)
	if err != nil {
		t.Fatalf("shouldOverwriteConfig() error = %v", err)
	}
	if !ok {
		t.Fatal("expected overwrite to be allowed with reinit")
	}
	if reader.called {
		t.Fatal("expected reinit to skip reading input")
	}
}

func TestShouldOverwriteConfigPromptYes(t *testing.T) {
	ok, err := shouldOverwriteConfig(bufio.NewReader(strings.NewReader("y\n")), io.Discard, true, false)
	if err != nil {
		t.Fatalf("shouldOverwriteConfig() error = %v", err)
	}
	if !ok {
		t.Fatal("expected overwrite to be allowed when user confirms")
	}
}

func TestShouldOverwriteConfigPromptNo(t *testing.T) {
	ok, err := shouldOverwriteConfig(bufio.NewReader(strings.NewReader("n\n")), io.Discard, true, false)
	if err != nil {
		t.Fatalf("shouldOverwriteConfig() error = %v", err)
	}
	if ok {
		t.Fatal("expected overwrite to be rejected when user declines")
	}
}
