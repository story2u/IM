package golden

import (
	"fmt"
	"strings"
)

// MarkdownReport renders a compact golden suite report for CI artifacts.
func MarkdownReport(report SuiteReport) string {
	var builder strings.Builder
	builder.WriteString("# Golden HTTP Report\n\n")
	builder.WriteString("This report validates or compares deterministic HTTP cases against baseline and Go endpoints.\n\n")
	builder.WriteString("## Summary\n\n")
	builder.WriteString("| Field | Value |\n")
	builder.WriteString("| --- | --- |\n")
	fmt.Fprintf(&builder, "| Suite | `%s` |\n", escapeTable(report.Suite))
	fmt.Fprintf(&builder, "| Mode | `%s` |\n", escapeTable(report.Mode))
	fmt.Fprintf(&builder, "| Match | `%t` |\n", report.Match)
	fmt.Fprintf(&builder, "| Cases | %d |\n\n", report.CaseCount)
	writeCaseTable(&builder, report.Cases)
	writeResultTable(&builder, report.Results)
	return builder.String()
}

func writeCaseTable(builder *strings.Builder, cases []CaseSummary) {
	builder.WriteString("## Cases\n\n")
	builder.WriteString("| Name | Method | Path |\n")
	builder.WriteString("| --- | --- | --- |\n")
	for _, testCase := range cases {
		fmt.Fprintf(builder,
			"| `%s` | `%s` | `%s` |\n",
			escapeTable(testCase.Name),
			escapeTable(testCase.Method),
			escapeTable(testCase.Path),
		)
	}
	builder.WriteString("\n")
}

func writeResultTable(builder *strings.Builder, results []Result) {
	if len(results) == 0 {
		return
	}
	builder.WriteString("## Results\n\n")
	builder.WriteString("| Case | Match | Baseline Status | Go Status | Diffs |\n")
	builder.WriteString("| --- | --- | ---: | ---: | --- |\n")
	for _, result := range results {
		fmt.Fprintf(builder,
			"| `%s` | `%t` | %d | %d | `%s` |\n",
			escapeTable(result.Case),
			result.Match,
			result.Baseline.StatusCode,
			result.Go.StatusCode,
			escapeTable(strings.Join(result.Diffs, "; ")),
		)
	}
	builder.WriteString("\n")
}

func escapeTable(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "`", "\\`")
	value = strings.ReplaceAll(value, "|", "\\|")
	return value
}
