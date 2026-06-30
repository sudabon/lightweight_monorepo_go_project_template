package arch

import (
	"os"
	"sort"
	"testing"
)

// sourceRoot は Clean Architecture のソースルート。このテストパッケージ
// （backend/tests/arch）からの相対パス。フォーク時にリネームする。
const sourceRoot = "../../internal/todo_app"

// expectedLayers は sourceRoot 直下に許可される唯一のディレクトリ群
// （内側から外側の順）。
var expectedLayers = []string{"domain", "app", "interfaces", "infra"}

// TestSourceStructure は sourceRoot 直下が4層のみで構成されることを検査する。
// 対象パッケージが未配置の場合はスキップする（テンプレート素の状態）。
func TestSourceStructure(t *testing.T) {
	if _, err := os.Stat(sourceRoot); os.IsNotExist(err) {
		t.Skipf("source root %q not found; add the application package to enable this check", sourceRoot)
	}

	entries, err := os.ReadDir(sourceRoot)
	if err != nil {
		t.Fatalf("read source root: %v", err)
	}

	found := map[string]bool{}
	for _, e := range entries {
		if e.IsDir() {
			found[e.Name()] = true
		}
	}

	for _, layer := range expectedLayers {
		if !found[layer] {
			t.Errorf("missing %q layer directory under %s", layer, sourceRoot)
		}
	}

	expected := map[string]bool{}
	for _, layer := range expectedLayers {
		expected[layer] = true
	}
	var unexpected []string
	for name := range found {
		if !expected[name] {
			unexpected = append(unexpected, name)
		}
	}
	sort.Strings(unexpected)
	if len(unexpected) > 0 {
		t.Errorf("source root should only contain Clean Architecture layers; unexpected: %v", unexpected)
	}
}
