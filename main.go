package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/IMEsec-USP/container-manager/adapters"
	"github.com/IMEsec-USP/container-manager/applications"
)

// func main() {
// 	dockerClient, err := client.NewEnvClient()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	ctx := context.Background()
// 	ctx, _ = context.WithDeadline(ctx, time.Now().Add(20 * time.Second))
// 	reader, err := dockerClient.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	_, err = io.Copy(os.Stdout, reader)
// 	if err != nil {
// 		panic(err)
// 	}
// }

type point struct {
	x, y int
}

func main() {
	dockerAdapter, err := adapters.NewDockerAdapter()
	if err != nil {
		panic(err)
	}

	reader, err := dockerAdapter.PullImage(context.TODO(), applications.LHTimer{})
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}

	containerConfig, err := dockerAdapter.GetContainerConfig(context.TODO(), applications.LHTimer{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Got container config")

	err = dockerAdapter.RemoveContainer(context.TODO(), containerConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println("removed container")

	err = dockerAdapter.RunImage(context.TODO(), applications.LHTimer{}, containerConfig)
	if err != nil {
		panic(err)
	}

	fmt.Println("Finished! :) :(")

	// reader, err := dockerAdapter.PullImage(context.Background(), applications.Ada{})
	// if err != nil {
	// 	panic(err)
	// }
	//
	// _, err = io.Copy(os.Stdout, reader)
	// if err != nil {
	// 	panic(err)
	// }
}
