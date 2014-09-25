package main

import (
	"bytes"
	"flag"
	"log"
	"os/exec"

	docker "github.com/fsouza/go-dockerclient"
)

var client *docker.Client

func main() {
	var err error
	client, err = docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		log.Fatal("Docker sock not connected", err)
	}
	var containerName string

	// Set up flags
	flag.StringVar(&containerName, "name", "", "Container Name")

	flag.Parse()
	container := containerByName(containerName)
	if container == nil {
		log.Fatalf("Could not find container %s", containerName)
	}

	log.Printf("%#v\n", container)
}

func dockerCmd(args ...string) (error, bytes.Buffer, bytes.Buffer) {
	cmd := exec.Command("docker", args...)
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	err := cmd.Wait()
	return err, stdout, stderr
}

func listContainers() []docker.APIContainers {
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}
	return containers
}

func containerByName(name string) *docker.APIContainers {
	for _, container := range listContainers() {
		//log.Printf("%#v", container)
		for _, cname := range container.Names {
			log.Printf("%s", cname)
			if name == cname[1:] {
				return &container
			}
		}
	}
	return nil
}
