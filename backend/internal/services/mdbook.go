package services

import (
	"fmt"
	"os/exec"
)

func BuildBook(sourceDir, buildDir string) error {
	cmd := exec.Command("mdbook", "build", sourceDir, "-d", buildDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mdbook build failed: %w: %s", err, string(out))
	}
	return nil
}
