//+build mage

package main

import (
	"log"
	"os"
	"os/exec"
)

func DBUp() {
	cmd := exec.Command("docker-compose", "-f", "build/package/docker-compose.yml", "up", "-d", "db")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func DBDown() {
	cmd := exec.Command("docker-compose", "-f ./build/package/docker-compose.yml", "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
