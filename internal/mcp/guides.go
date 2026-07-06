package mcp

import (
	"fmt"
	"io/fs"
	"strings"
)

func ResourcesFromFS(fsys fs.FS, uriScheme string) []Resource {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil
	}

	var resources []Resource
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(e.Name(), ".md")

		resources = append(resources, Resource{
			URI:         fmt.Sprintf("%s://%s", uriScheme, name),
			Name:        name,
			Description: fmt.Sprintf("s1ctl guide: %s", name),
			MimeType:    "text/markdown",
			Read:        fsReader(fsys, e.Name()),
		})
	}
	return resources
}

func fsReader(fsys fs.FS, name string) func() (string, error) {
	return func() (string, error) {
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
}
