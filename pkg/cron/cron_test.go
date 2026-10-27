package cron

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"charm.land/log/v2"
)

func TestCronLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetLevel(log.DebugLevel)
	clogger := cronLogger{logger}
	clogger.Info("foo")
	clogger.Error(fmt.Errorf("bar"), "test")
	out := buf.String()
	for _, want := range []string{"foo", "test", "err=bar"} {
		if !strings.Contains(out, want) {
			t.Errorf("log output %q does not contain %q", out, want)
		}
	}
}

func TestSchedularAddRemove(t *testing.T) {
	s := NewScheduler(context.TODO())
	id, err := s.AddFunc("* * * * *", func() {})
	if err != nil {
		t.Fatal(err)
	}
	s.Remove(id)
}
