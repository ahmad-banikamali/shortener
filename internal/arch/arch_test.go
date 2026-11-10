package arch_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

var rules = []struct {
	dir       string
	forbidden []string
}{
	// core must not know about adapters or frameworks
	{"../core", []string{"/internal/adapter", "net/http", "database/sql", "modernc.org/sqlite"}},
	// domain is the innermost layer
	{"../core/domain", []string{"/internal/core/port", "/internal/core/service"}},
	// adapters must not depend on each other
	{"../adapter/inbound", []string{"/internal/adapter/outbound"}},
}

func TestDependencyRules(t *testing.T) {
	for _, rule := range rules {
		fset := token.NewFileSet()
		err := filepath.WalkDir(rule.dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
				return err
			}
			f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if err != nil {
				return err
			}
			for _, imp := range f.Imports {
				p, _ := strconv.Unquote(imp.Path.Value)
				for _, bad := range rule.forbidden {
					if strings.Contains(p, bad) {
						t.Errorf("%s imports %q (forbidden under %s)", path, p, rule.dir)
					}
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}
