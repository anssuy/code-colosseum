package main

import (
	"log"
	"os"
	"os/exec"
)

func build(tag, context string) {
	cmd := exec.Command("docker", "build", "-t", tag, context)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("build %s: %v", tag, err)
	}
}

func main() {
	build("sandbox-node:latest", "internal/docker/node")
	build("sandbox-python:latest", "internal/docker/python")
}
