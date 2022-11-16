package main

import (
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"time"
)

func main() {
	fmt.Println("hello world")
	CreateFolder("Testf3", "Testf")

}

func getJenkins() (*gojenkins.Jenkins, context.Context) {
	ctx := context.Background()
	jenkins := gojenkins.CreateJenkins(nil, "http://35.184.233.61:8080/", "caner", "test123")
	// Provide CA certificate if server is using self-signed certificate
	// caCert, _ := ioutil.ReadFile("/tmp/ca.crt")
	// jenkins.Requester.CACert = caCert
	_, err := jenkins.Init(context.Background())
	if err != nil {
		panic("Something Went Wrong")
	}
	return jenkins, ctx
}

func CreateFolder(name string, parent string) {
	jenkins, ctx := getJenkins()
	folder, err := jenkins.CreateFolder(ctx, name, parent)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Folder name: %s", folder.GetName())
}

func BuildJob(name string) {
	jenkins, ctx := getJenkins()
	queueid, err := jenkins.BuildJob(ctx, name, nil)
	if err != nil {
		panic(err)
	}
	build, err := jenkins.GetBuildFromQueueID(ctx, queueid)
	if err != nil {
		panic(err)
	}

	// Wait for build to finish
	for build.IsRunning(ctx) {
		time.Sleep(5000 * time.Millisecond)
		build.Poll(ctx)
	}

	fmt.Printf("build number %d with result: %v\n", build.GetBuildNumber(), build.GetResult())
}
