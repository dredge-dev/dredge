package workflow

import (
	"fmt"
	"strings"

	"github.com/dredge-dev/dredge/internal/config"
)

func insertGo(insert *config.Insert, currentContent, text string) (string, error) {
	if insert.Section == "import" {
		return insertImportGo(currentContent, text), nil
	} else if strings.HasPrefix(insert.Section, "func") {
		return insertInFuncGo(insert, currentContent, text), nil
	} else {
		return "", fmt.Errorf("Unknown section %s (only import and func are supported in go)", insert.Section)
	}
}

func insertImportGo(currentContent, text string) string {
	var output []string

	lines := strings.Split(currentContent, "\n")
	var packageLine string
	var imports []string
	inImports := false
	headerAdded := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "package") {
			packageLine = line
		} else if len(trimmedLine) == 0 {
			// Skip empty lines
		} else if strings.HasPrefix(trimmedLine, "import") {
			if strings.Contains(trimmedLine, "(") {
				inImports = true
				continue
			} else {
				imports = append(imports, strings.TrimSpace(strings.Trim(trimmedLine, "import ")))
			}
		} else if inImports {
			if strings.Contains(trimmedLine, ")") {
				inImports = false
			} else {
				imports = append(imports, trimmedLine)
			}
		} else if !headerAdded {
			output = append(output, packageLine)
			output = append(output, "\nimport (")
			importMap := make(map[string]bool)
			imports = append(imports, strings.Split(text, "\n")...)
			for _, i := range imports {
				if _, present := importMap[i]; !present {
					importMap[i] = true
					output = append(output, "\t"+i)
				}
			}
			output = append(output, ")\n")
			output = append(output, line)
			headerAdded = true
		} else {
			output = append(output, line)
		}
	}

	return strings.Join(output, "\n")
}

func insertInFuncGo(insert *config.Insert, currentContent, text string) string {
	var output []string

	lines := strings.Split(currentContent, "\n")
	atEnd := insert.Placement == config.INSERT_END
	inSection := false
	bracketCounter := 0

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, insert.Section) {
			if !atEnd {
				output = append(output, line)
				output = append(output, text)
				continue
			}
			inSection = true
			bracketCounter = 0
		}
		if inSection {
			bracketCounter += strings.Count(trimmedLine, "{")
			bracketCounter -= strings.Count(trimmedLine, "}")
			if bracketCounter == 0 {
				inSection = false
				if atEnd {
					output = append(output, text)
				}
			}
		}
		output = append(output, line)
	}

	return strings.Join(output, "\n")
}
