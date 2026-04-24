package hook

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"charm.land/log/v2"
)

func TestSummarizeOutput(t *testing.T) {
	t.Parallel()

	short := "hello"
	if got := summarizeOutput(short); got != short {
		t.Fatalf("summarizeOutput(short) = %q, want %q", got, short)
	}

	long := bytes.Repeat([]byte("a"), 600)
	got := summarizeOutput(string(long))
	if len(got) <= 512 {
		t.Fatalf("expected truncated output, got len=%d", len(got))
	}
}

func TestRunCommandTimeout(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	start := time.Now()
	err := runCommand(
		context.Background(),
		log.Default(),
		1,
		bytes.NewBufferString(""),
		io.Writer(&stdout),
		io.Writer(&stderr),
		"bash",
		"-lc",
		"sleep 2",
	)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if elapsed := time.Since(start); elapsed >= 2*time.Second {
		t.Fatalf("expected hook timeout before command completion, elapsed=%s", elapsed)
	}
}
