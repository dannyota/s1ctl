package cli

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh/spinner"
)

func isTTY(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func runWithSpinner(title string, fn func() error) error {
	if noProgress || !isTTY(os.Stderr) {
		return fn()
	}
	var fnErr error
	if err := spinner.New().Title(title).Action(func() { fnErr = fn() }).Run(); err != nil {
		return err
	}
	return fnErr
}

func printProgress(resource string, fetched, total int) {
	if noProgress || !isTTY(os.Stderr) {
		return
	}
	if total > 0 {
		fmt.Fprintf(os.Stderr, "\rFetching %s... %d/%d", resource, fetched, total)
	} else {
		fmt.Fprintf(os.Stderr, "\rFetching %s... %d", resource, fetched)
	}
}

func clearProgress() {
	if noProgress || !isTTY(os.Stderr) {
		return
	}
	fmt.Fprint(os.Stderr, "\r\033[K")
}
