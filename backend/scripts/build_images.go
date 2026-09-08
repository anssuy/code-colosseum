package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/anssuy/code-colosseum/backend/internal/language"
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
	for _, def := range language.Registry() {
		build(def.Image, def.DockerContext)
	}
}
