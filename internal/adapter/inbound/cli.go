package inbound

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

type CLI struct {
	shorten port.ShortenUseCase
	resolve port.ResolveUseCase
	stats   port.StatsUseCase
	baseURL string
	out     io.Writer
	errOut  io.Writer
}

func NewCLI(s port.ShortenUseCase, r port.ResolveUseCase, st port.StatsUseCase,
	baseURL string, out, errOut io.Writer) *CLI {
	return &CLI{shorten: s, resolve: r, stats: st, baseURL: baseURL, out: out, errOut: errOut}
}

// Run returns the process exit code.
func (c *CLI) Run(ctx context.Context, args []string) int {
	if len(args) != 2 {
		return c.usage()
	}
	cmd, arg := args[0], args[1]

	switch cmd {
	case "shorten":
		link, err := c.shorten.Shorten(ctx, arg, 0)
		if err != nil {
			return c.fail(err)
		}
		fmt.Fprintf(c.out, "%s/%s\n", c.baseURL, link.Code)
	case "resolve":
		link, err := c.resolve.Resolve(ctx, arg)
		if err != nil {
			return c.fail(err)
		}
		fmt.Fprintln(c.out, link.URL)
	case "stats":
		link, err := c.stats.Stats(ctx, arg)
		if err != nil {
			return c.fail(err)
		}
		fmt.Fprintf(c.out, "code:    %s\nurl:     %s\nvisits:  %d\ncreated: %s\n",
			link.Code, link.URL, link.Visits, link.CreatedAt.Format("2006-01-02 15:04:05"))
	default:
		return c.usage()
	}
	return 0
}

func (c *CLI) usage() int {
	fmt.Fprintln(c.errOut, "usage: shortener <shorten|resolve|stats> <url|code>")
	return 2
}

// Same job as writeError in the HTTP adapter: domain error → this protocol's response.
func (c *CLI) fail(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidURL):
		fmt.Fprintln(c.errOut, "error: invalid url")
	case errors.Is(err, domain.ErrNotFound):
		fmt.Fprintln(c.errOut, "error: link not found")
	default:
		fmt.Fprintln(c.errOut, "error: internal error")
	}
	return 1
}
