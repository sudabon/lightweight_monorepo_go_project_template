package arch

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// sourceRootName はソースルートのディレクトリ/パッケージ名。依存検査は
// import パス中のこのセグメント直後の層名をキーに行う。フォーク時にリネーム。
const sourceRootName = "todo_app"

// forbiddenForDomain は domain 層が import してはならない層。
var forbiddenForDomain = map[string]bool{"app": true, "interfaces": true, "infra": true}

// TestDomainLayerDependencies は domain 層が外側の層へ内向き import しない
// ことを検査する。domain ディレクトリが無い場合はスキップする。
func TestDomainLayerDependencies(t *testing.T) {
	domainPath := filepath.Join(sourceRoot, "domain")
	if _, err := os.Stat(domainPath); os.IsNotExist(err) {
		t.Skipf("domain layer %q not found; add the application package to enable this check", domainPath)
	}

	fset := token.NewFileSet()
	var violations []string

	err := filepath.WalkDir(domainPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, perr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if perr != nil {
			return perr
		}
		for _, imp := range file.Imports {
			layer, ok := layerAfterSourceRoot(strings.Trim(imp.Path.Value, `"`))
			if ok && forbiddenForDomain[layer] {
				rel, _ := filepath.Rel(domainPath, path)
				violations = append(violations,
					rel+": domain layer cannot import from "+layer+" layer")
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk domain layer: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf("dependency rule violations:\n%s", strings.Join(violations, "\n"))
	}
}

// layerAfterSourceRoot は import パス中で sourceRootName の直後に来る層名を返す。
// モジュールパス自体に sourceRootName と同名のセグメントが含まれる場合
// （例: module github.com/sudabon/todo_app 配下の todo_app パッケージ）に
// 誤検出しないよう、最後に出現する sourceRootName を基準にする。
// これは internal/<sourceRootName>/<layer> のネスト構成を前提とする。
// 構成をフラット化する場合は sourceRootName を一意な名前にリネームすること。
func layerAfterSourceRoot(importPath string) (string, bool) {
	segs := strings.Split(importPath, "/")
	for i := len(segs) - 1; i >= 0; i-- {
		if segs[i] == sourceRootName && i+1 < len(segs) {
			return segs[i+1], true
		}
	}
	return "", false
}
