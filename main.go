package main

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
)

func main() {
	var err error
	var cli *client.Client

	fmt.Printf("Go version: %s\n", runtime.Version())

	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	FailOnError(err, "Error with NewClientWithOpts")

	var out io.ReadCloser
	imageName := "postgres:13.10"
	out, err = cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	FailOnError(err, "Error with ImagePull")

	io.Copy(os.Stdout, out)

	freePort := getFreePort()
	log.Printf("free port: %v", freePort)
	var portBindings map[nat.Port][]nat.PortBinding
	_, portBindings, err = nat.ParsePortSpecs([]string{freePort + ":" + strconv.Itoa(5432)})
	FailOnError(err, "Cannot create portBinding for free port: "+freePort)
	for port, bindings := range portBindings {
		log.Printf("Port: %v; Bindings: %v", port, bindings)
	}

	env := []string{"POSTGRES_USER=myuser", "POSTGRES_PASSWORD=mypassword", "POSTGRES_DB=mydb"}
	var resp container.CreateResponse
	resp, err = cli.ContainerCreate(context.Background(), &container.Config{
		Image: imageName,
		Env:   env,
	}, &container.HostConfig{
		PortBindings: portBindings,
	}, nil, nil, "database")
	FailOnError(err, "Cannot create docker container")
	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	FailOnError(err, "Cannot start docker client")

}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getFreePort() (freePort string) {
	l, err := net.Listen("tcp", ":0")
	FailOnError(err, "Cannot listen on random port")
	defer l.Close()
	freePort = l.Addr().String()
	return
}
