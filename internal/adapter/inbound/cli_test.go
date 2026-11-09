package inbound_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/ahmad-banikamali/shortener/internal/adapter/inbound"
	"github.com/ahmad-banikamali/shortener/internal/core/domain"
)

func TestCLI_Shorten_PrintsShortURL(t *testing.T) {
	var out, errOut bytes.Buffer
	cli := inbound.NewCLI(fakeShorten{}, fakeResolve{}, fakeStats{}, "http://x", &out, &errOut)

	code := cli.Run(context.Background(), []string{"shorten", "https://example.com"})

	if code != 0 || out.String() != "http://x/abc\n" {
		t.Fatalf("code=%d out=%q", code, out.String())
	}
}

func TestCLI_Resolve_NotFound(t *testing.T) {
	var out, errOut bytes.Buffer
	cli := inbound.NewCLI(fakeShorten{}, fakeResolve{err: domain.ErrNotFound}, fakeStats{}, "http://x", &out, &errOut)

	code := cli.Run(context.Background(), []string{"resolve", "nope"})

	if code != 1 || !strings.Contains(errOut.String(), "not found") {
		t.Fatalf("code=%d err=%q", code, errOut.String())
	}
}

func TestCLI_BadUsage(t *testing.T) {
	var out, errOut bytes.Buffer
	cli := inbound.NewCLI(fakeShorten{}, fakeResolve{}, fakeStats{}, "http://x", &out, &errOut)

	if code := cli.Run(context.Background(), []string{"oops"}); code != 2 {
		t.Fatalf("expected exit code 2, got %d", code)
	}
}
