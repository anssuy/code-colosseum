package sandbox

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTimeout  = 5 * time.Second
	defaultMemoryMB = 128
	defaultCPU      = "0.5"
	maxOutputLen    = 64 * 1024 // 64KB
)

type runConfig struct {
	timeout  time.Duration
	memoryMB int
	cpu      string
	stdin    string
}

func defaultRunConfig() runConfig {
	return runConfig{
		timeout:  defaultTimeout,
		memoryMB: defaultMemoryMB,
		cpu:      defaultCPU,
	}
}

func runInContainer(ctx context.Context, image, filename, code string, cmd []string, cfg runConfig) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, cfg.timeout)
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

	containerName := "sandbox-" + randomID()

	args := []string{
		"run", "--rm", "-i",
		"--network", "none",
		"--name", containerName,
		"--memory", strconv.Itoa(cfg.memoryMB) + "m",
		"--cpus", cfg.cpu,
		"--pids-limit", "64",
		"-v", fmt.Sprintf("%s:/code:ro", tmpDir),
		"-w", "/code",
		image,
	}
	args = append(args, cmd...)

	dockerCmd := exec.CommandContext(ctx, "docker", args...)
	if cfg.stdin != "" {
		dockerCmd.Stdin = strings.NewReader(cfg.stdin)
	}

	var out bytes.Buffer
	dockerCmd.Stdout = &limitWriter{buf: &out, max: maxOutputLen}
	dockerCmd.Stderr = dockerCmd.Stdout

	err = dockerCmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		exec.Command("docker", "kill", containerName).Run()
		return out.String(), fmt.Errorf("execution timed out")
	}
	if err != nil {
		return out.String(), fmt.Errorf("execution error: %w", err)
	}
	return out.String(), nil
}

func randomID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
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
