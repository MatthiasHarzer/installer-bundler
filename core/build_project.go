package core

import (
	"os"
	"os/exec"
)

func BuildProject(projectDir string) (string, error) {
	cmd := exec.Command("make", "build")
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return projectDir + "/bin/installer-runtime.exe", nil
}
