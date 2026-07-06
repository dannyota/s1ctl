package mcp

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Tool struct {
	Name        string
	Description string
	InputSchema map[string]any
	Run         func(args map[string]any) (string, error)
}

var skipCommands = map[string]bool{
	"mcp": true, "completion": true, "docs": true, "help": true,
}

func ToolsFromCobra(root *cobra.Command) []Tool {
	var tools []Tool
	walkCommands(root, nil, &tools)
	return tools
}

func walkCommands(cmd *cobra.Command, path []string, tools *[]Tool) {
	for _, c := range cmd.Commands() {
		if c.Hidden || (len(path) == 0 && skipCommands[c.Name()]) {
			continue
		}
		cur := append(append([]string(nil), path...), c.Name())

		if c.HasSubCommands() {
			walkCommands(c, cur, tools)
			continue
		}

		if c.RunE == nil && c.Run == nil {
			continue
		}

		t := buildTool(c, cur)
		*tools = append(*tools, t)
	}
}

var skipFlags = map[string]bool{
	"help": true, "json": true, "output": true,
	"verbose": true, "no-progress": true, "config": true,
}

func buildTool(cmd *cobra.Command, path []string) Tool {
	name := strings.Join(path, "_")

	desc := cmd.Short
	if cmd.Long != "" {
		desc = cmd.Long
	}
	if hasMutationFlag(cmd) {
		desc += " [mutation: requires --yes to apply, dry-run by default]"
	}

	schema := buildInputSchema(cmd)

	return Tool{
		Name:        name,
		Description: desc,
		InputSchema: schema,
		Run:         makeRunner(cmd, path),
	}
}

func buildInputSchema(cmd *cobra.Command) map[string]any {
	properties := map[string]any{}
	var required []string

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if skipFlags[f.Name] {
			return
		}
		if f.Hidden {
			return
		}

		prop := flagToProperty(f)
		properties[f.Name] = prop

		ann := cmd.Flag(f.Name)
		if ann != nil {
			if _, ok := ann.Annotations[cobra.BashCompOneRequiredFlag]; ok {
				required = append(required, f.Name)
			}
		}
	})

	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if skipFlags[f.Name] {
			return
		}
		if f.Hidden {
			return
		}
		if f.Name == "read-only" {
			properties[f.Name] = flagToProperty(f)
		}
	})

	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func flagToProperty(f *pflag.Flag) map[string]any {
	prop := map[string]any{}

	switch f.Value.Type() {
	case "bool":
		prop["type"] = "boolean"
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		prop["type"] = "integer"
	case "float32", "float64":
		prop["type"] = "number"
	case "stringSlice", "stringArray":
		prop["type"] = "array"
		prop["items"] = map[string]any{"type": "string"}
	default:
		prop["type"] = "string"
	}

	if f.Usage != "" {
		prop["description"] = f.Usage
	}
	if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
		prop["default"] = f.DefValue
	}

	return prop
}

func hasMutationFlag(cmd *cobra.Command) bool {
	return cmd.Flags().Lookup("yes") != nil
}

func makeRunner(cmd *cobra.Command, path []string) func(map[string]any) (string, error) {
	return func(args map[string]any) (string, error) {
		root := cmd.Root()

		cliArgs := append([]string(nil), path...)
		cliArgs = append(cliArgs, "--json", "--no-progress")

		for k, v := range args {
			flag := cmd.Flags().Lookup(k)
			if flag == nil {
				flag = cmd.InheritedFlags().Lookup(k)
			}
			if flag == nil {
				continue
			}

			switch flag.Value.Type() {
			case "bool":
				b, _ := toBool(v)
				if b {
					cliArgs = append(cliArgs, "--"+k)
				}
			case "stringSlice", "stringArray":
				if arr, ok := v.([]any); ok {
					for _, item := range arr {
						cliArgs = append(cliArgs, "--"+k, fmt.Sprint(item))
					}
				}
			default:
				cliArgs = append(cliArgs, "--"+k, fmt.Sprint(v))
			}
		}

		// Capture both cobra output and raw os.Stdout writes.
		// Many commands write to os.Stdout directly via fmt.Printf;
		// os.Pipe intercepts those writes too.
		origStdout := os.Stdout
		pr, pw, err := os.Pipe()
		if err != nil {
			return "", fmt.Errorf("create pipe: %w", err)
		}
		os.Stdout = pw

		var stderr bytes.Buffer
		root.SetArgs(cliArgs)
		root.SetOut(pw)
		root.SetErr(&stderr)

		execErr := root.Execute()

		_ = pw.Close()
		os.Stdout = origStdout

		var stdout bytes.Buffer
		_, _ = stdout.ReadFrom(pr)
		_ = pr.Close()

		out := stdout.String()
		if errOut := stderr.String(); errOut != "" && out == "" {
			out = errOut
		}
		if out == "" && execErr != nil {
			out = execErr.Error()
		}

		return out, execErr
	}
}

func toBool(v any) (bool, bool) {
	switch b := v.(type) {
	case bool:
		return b, true
	case string:
		val, err := strconv.ParseBool(b)
		return val, err == nil
	default:
		return false, false
	}
}
