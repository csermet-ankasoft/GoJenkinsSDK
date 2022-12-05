package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bndr/gojenkins"
)

func main() {
	/*
		fmt.Println("hello world")
		jobs := GetAllJobsName()
		Build := GetAllBuildsID(jobs[0])
		fmt.Println(Build)

		jenkins, ctx := getJenkins()
		_, err := jenkins.GetFolder(ctx, "Testf", "Testf3")
		if err != nil {
			panic(err)
		}

		//fmt.Printf("Folder name: %s", folder.GetName())
	*/
	CreatJob("Test21", "Testf")
}

func CreateFolder(name string, parent string) {
	jenkins, ctx := getJenkins()
	var folderName = ""
	if parent == "" {
		folder, err := jenkins.CreateFolder(ctx, name)
		if err != nil {
			panic(err)
		}
		folderName = folder.GetName()
	} else {
		folder, err := jenkins.CreateFolder(ctx, name, parent)
		if err != nil {
			panic(err)
		}
		folderName = folder.GetName()
	}

	fmt.Printf("Folder name: %s", folderName)
}

func GetFolder(name string, parent string) {
	jenkins, ctx := getJenkins()
	var folderName = ""
	if parent == "" {
		folder, err := jenkins.GetFolder(ctx, name)
		if err != nil {
			panic(err)
		}
		folderName = folder.GetName()
	} else {
		folder, err := jenkins.GetFolder(ctx, name, parent)
		if err != nil {
			panic(err)
		}
		folderName = folder.GetName()
	}

	fmt.Printf("Folder name: %s", folderName)
}

func getJobConfig(name string, parent string) string {
	jenkins, ctx := getJenkins()
	var jobconfig = ""
	if parent == "" {
		job, err := jenkins.GetJob(ctx, name)
		jobconfig, err = job.GetConfig(ctx)
		if err != nil {
			panic(err)
		}
	} else {
		job, err := jenkins.GetJob(ctx, name, parent)
		jobconfig, err = job.GetConfig(ctx)
		if err != nil {
			panic(err)
		}
	}

	return jobconfig
}

func CreatJob(name string, parent string) {
	config := getJobConfig("Test", "")
	jenkins, ctx := getJenkins()
	jobName := ""
	if parent == "" {
		job, err := jenkins.CreateJob(ctx, config, name)
		jobName = job.GetName()
		if err != nil {
			panic(err)
		}
	} else {
		job, err := jenkins.CreateJobInFolder(ctx, config, name, parent)
		jobName = job.GetName()
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf(jobName)
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

func GetAllJobsName() []*gojenkins.Job {
	jenkins, ctx := getJenkins()
	allJobs, err := jenkins.GetAllJobs(ctx)
	if err != nil {
		panic(err)
	}
	return allJobs
}

func GetAllBuildsID(job *gojenkins.Job) []gojenkins.JobBuild {
	jenkins, ctx := getJenkins()
	allJobs, err := jenkins.GetAllBuildIds(ctx, job.GetName())

	if err != nil {
		panic(err)
	}
	return allJobs
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
