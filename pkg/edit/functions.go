package edit

import (
	"fmt"
	"strings"
	"text/template"
)

var TemplateFunctions template.FuncMap = template.FuncMap{
	"join":     join,
	"decorate": decorate,
}

func join(sep string, values []string) string {
	return strings.Join(values, sep)
}

func decorate(format string, values []string) []string {
	result := make([]string, len(values))

	if format == "" {
		format = "%s"
	}

	for i := 0; i < len(values); i++ {
		result[i] = fmt.Sprintf(format, values[i])
	}

	return result
}
