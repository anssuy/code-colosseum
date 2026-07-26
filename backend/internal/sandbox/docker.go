package sandbox

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	execTimeout  = 5 * time.Second
	maxOutputLen = 64 * 1024 // 64KB
)

func runInContainer(ctx context.Context, image, filename, code string, cmd []string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, execTimeout)
	defer cancel()

	tmpDir, err := os.MkdirTemp("", "sandbox-*")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	codePath := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(codePath, []byte(code), 0o644); err != nil {
		return "", fmt.Errorf("write code file: %w", err)
	}

	args := []string{
		"run", "--rm",
		"--network", "none",
		"--memory", "128m",
		"--cpus", "0.5",
		"--pids-limit", "64",
		"-v", fmt.Sprintf("%s:/code:ro", tmpDir),
		"-w", "/code",
		image,
	}
	args = append(args, cmd...)

	dockerCmd := exec.CommandContext(ctx, "docker", args...)

	var out bytes.Buffer
	dockerCmd.Stdout = &limitWriter{buf: &out, max: maxOutputLen}
	dockerCmd.Stderr = dockerCmd.Stdout

	err = dockerCmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("execution timed out")
	}
	if err != nil {
		return out.String(), fmt.Errorf("execution error: %w", err)
	}
	return out.String(), nil
}

type limitWriter struct {
	buf *bytes.Buffer
	max int
}

func (w *limitWriter) Write(p []byte) (int, error) {
	remaining := w.max - w.buf.Len()
	if remaining <= 0 {
		return len(p), nil
	}
	if len(p) > remaining {
		p = p[:remaining]
	}
	w.buf.Write(p)
	return len(p), nil
}
