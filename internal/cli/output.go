package cli

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"danny.vn/s1/graphql"
	"danny.vn/s1/mgmt"
	"danny.vn/s1/sdl"
)

var (
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).PaddingRight(2)
	cellStyle   = lipgloss.NewStyle().PaddingRight(2)
)

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func printCSV(headers []string, rows [][]string) error {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(headers); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func printTable(headers []string, rows [][]string) {
	t := table.New().
		Headers(headers...).
		Rows(rows...).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return cellStyle
		}).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderColumn(false).
		BorderHeader(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("8")))

	fmt.Println(t)
}

func printOutput(w io.Writer, headers []string, rows [][]string, items any, count, total int, resource string, all bool) error {
	switch outputFormat {
	case "json":
		return printJSON(items)
	case "csv":
		return printCSV(headers, rows)
	default:
		printTable(headers, rows)
		printFooter(w, count, total, resource, all)
		return nil
	}
}

func printFooter(w io.Writer, count, total int, resource string, all bool) {
	if all || total <= count {
		fmt.Fprintf(w, "\n%s\n", pluralize(count, resource))
		return
	}
	if total > 0 {
		fmt.Fprintf(w, "\nShowing %d of %d %ss (use --all to fetch all)\n", count, total, resource)
		return
	}
	fmt.Fprintf(w, "\n%s\n", pluralize(count, resource))
}

func printError(w io.Writer, err error) {
	if outputFormat == "json" {
		printJSONError(err)
		return
	}
	if verbose {
		printVerboseError(w, err)
		return
	}
	fmt.Fprintln(w, "Error:", err)
}

func printJSONError(err error) {
	out := map[string]any{"error": map[string]any{"message": err.Error()}}

	var mgmtErr *mgmt.APIError
	var sdlErr *sdl.APIError
	var gqlHTTP *graphql.HTTPError
	var gqlQuery *graphql.QueryError

	switch {
	case errors.As(err, &mgmtErr):
		out["error"] = map[string]any{
			"status": mgmtErr.Status,
			"title":  mgmtErr.Title,
			"detail": mgmtErr.Detail,
			"body":   string(mgmtErr.RawBody),
		}
	case errors.As(err, &sdlErr):
		out["error"] = map[string]any{
			"status": sdlErr.Status,
			"body":   string(sdlErr.Body),
		}
	case errors.As(err, &gqlHTTP):
		out["error"] = map[string]any{
			"status": gqlHTTP.Status,
			"body":   string(gqlHTTP.Body),
		}
	case errors.As(err, &gqlQuery):
		msgs := make([]string, len(gqlQuery.Errors))
		for i, e := range gqlQuery.Errors {
			msgs[i] = e.Message
		}
		out["error"] = map[string]any{
			"messages": msgs,
		}
	}

	_ = printJSON(out)
}

func printVerboseError(w io.Writer, err error) {
	fmt.Fprintln(w, "Error:", err)

	var mgmtErr *mgmt.APIError
	var sdlErr *sdl.APIError
	var gqlHTTP *graphql.HTTPError

	switch {
	case errors.As(err, &mgmtErr):
		if mgmtErr.Title != "" {
			fmt.Fprintf(w, "\n  Title:  %s\n", mgmtErr.Title)
		}
		if mgmtErr.Detail != "" {
			fmt.Fprintf(w, "  Detail: %s\n", mgmtErr.Detail)
		}
		if len(mgmtErr.RawBody) > 0 {
			fmt.Fprintf(w, "  Body:   %s\n", mgmtErr.RawBody)
		}
	case errors.As(err, &sdlErr):
		if len(sdlErr.Body) > 0 {
			fmt.Fprintf(w, "\n  Body: %s\n", sdlErr.Body)
		}
	case errors.As(err, &gqlHTTP):
		if len(gqlHTTP.Body) > 0 {
			fmt.Fprintf(w, "\n  Body: %s\n", gqlHTTP.Body)
		}
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func boolIcon(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func pluralize(n int, singular string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	last := singular[len(singular)-1]
	if last == 'y' && len(singular) > 1 {
		prev := singular[len(singular)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return fmt.Sprintf("%d %sies", n, singular[:len(singular)-1])
		}
	}
	return fmt.Sprintf("%d %ss", n, singular)
}
