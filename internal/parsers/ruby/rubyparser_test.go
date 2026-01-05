package rubyparser_test

import (
	"testing"

	"github.com/mlw157/scout/internal/models"
	rubyparser "github.com/mlw157/scout/internal/parsers/ruby"
)

const testFilePath = "../../../testcases/parsers/ruby/"

func TestParseGemfileLock(t *testing.T) {
	t.Run("test extract correct number of dependencies", func(t *testing.T) {
		testFile := testFilePath + "Gemfile.lock"
		parser := rubyparser.NewRubyParser()
		dependencies, err := parser.ParseFile(testFile)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		got := len(dependencies)
		want := 127

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

	t.Run("test extract correct dependencies", func(t *testing.T) {
		testFile := testFilePath + "Gemfile.lock"
		parser := rubyparser.NewRubyParser()
		dependencies, _ := parser.ParseFile(testFile)

		// First gem in specs
		assertEqualDependency(t, dependencies[0], models.Dependency{Name: "actionmailer", Version: "4.2.5", Ecosystem: "gem"})
		// A gem in the middle
		found := false
		for _, dep := range dependencies {
			if dep.Name == "nokogiri" && dep.Version == "1.6.7.2" && dep.Ecosystem == "gem" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected to find nokogiri 1.6.7.2")
		}
	})
}

func assertEqualDependency(t testing.TB, got, want models.Dependency) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
