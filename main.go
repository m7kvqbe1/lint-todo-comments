package main

import (
	"bufio"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var (
	includeFlag string
)

func init() {
	flag.StringVar(&includeFlag, "include", "", "Comma-separated list of file extensions to include (e.g., .ts,.js)")
}

var todoAnalyzer = &analysis.Analyzer{
	Name: "todochecker",
	Doc:  "reports TODO comments within files",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	includeExtensions := map[string]bool{}

	for _, ext := range strings.Split(includeFlag, ",") {
		includeExtensions[strings.TrimSpace(ext)] = true
	}

	for _, f := range pass.Files {
		filename := pass.Fset.Position(f.Pos()).Filename

		if len(includeExtensions) > 0 {
			if _, ok := includeExtensions[filepath.Ext(filename)]; !ok {
				continue // Skip extensions not in include list
			}
		}

		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNumber := 1

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "TODO") {
				pass.Reportf(f.Pos(), "TODO comment detected in %s at line %d", filepath.Base(filename), lineNumber)
			}
			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func main() {
	flag.Parse()
	singlechecker.Main(todoAnalyzer)
}
