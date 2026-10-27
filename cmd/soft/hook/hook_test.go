package hook

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

func TestRunCustomHookPropagatesExitError(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(executable)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&0o111 == 0 {
		t.Skip("test binary is not executable on this platform")
	}

	t.Setenv("SOFT_SERVE_TEST_HELPER", "1")
	t.Setenv("SOFT_SERVE_TEST_HELPER_EXIT", "1")
	err = runCustomHook(
		context.Background(),
		log.Default(),
		0,
		executable,
		io.Reader(nil),
		io.Discard,
		io.Discard,
		"-test.run=TestCustomHookHelperProcess",
	)
	if err == nil {
		t.Fatal("runCustomHook returned nil for a failing custom hook")
	}

	if _, ok := errors.AsType[*exec.ExitError](err); !ok {
		t.Fatalf("runCustomHook error = %T %v, want exec.ExitError", err, err)
	}
}

func TestRunCustomHookIgnoresMissingOrNonExecutableHook(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing-hook")
	if err := runCustomHook(context.Background(), log.Default(), 0, missing, nil, io.Discard, io.Discard); err != nil {
		t.Fatalf("runCustomHook(missing) error = %v", err)
	}

	nonExecutable := filepath.Join(t.TempDir(), "non-executable-hook")
	if err := os.WriteFile(nonExecutable, []byte("exit 1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := runCustomHook(context.Background(), log.Default(), 0, nonExecutable, nil, io.Discard, io.Discard); err != nil {
		t.Fatalf("runCustomHook(non-executable) error = %v", err)
	}
}

func TestCustomHookHelperProcess(t *testing.T) {
	if os.Getenv("SOFT_SERVE_TEST_HELPER") != "1" {
		return
	}
	if os.Getenv("SOFT_SERVE_TEST_HELPER_EXIT") == "1" {
		os.Exit(1)
	}
	os.Exit(0)
}
