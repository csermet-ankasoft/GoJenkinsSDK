package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bndr/gojenkins"
)

func main() {
	//CheckAllNodeStatus()
	//CreateNode("Test2", 1, "Nope Desc", "/home", "NopeLabel") //map[string]string{"method": "JNLPLauncher"}
	DeleteNode("Test2")
}

func DeleteNode(name string) {
	jenkins, ctx := getJenkins()

	result, err := jenkins.DeleteNode(ctx, name)

	if err != nil {
		fmt.Print("Error: ", err)
	} else {
		if result != true {
			fmt.Print("Node Deleted")
		} else {
			fmt.Print("Node Not Deleted")
		}
	}
}

func CreateNode(name string, numExecutors int, description string, remoteFS string, label string) {
	jenkins, ctx := getJenkins()

	node, err := jenkins.CreateNode(ctx, name, numExecutors, description, remoteFS, label)

	if err != nil {
		fmt.Print("Error: ", err)
	} else {
		fmt.Print(node)
	}
}

func GetNode(name string) {
	jenkins, ctx := getJenkins()

	node, _ := jenkins.GetNode(ctx, name) // "(built-in)"
	if node != nil {
		fmt.Print(node.GetName(), " - ")
		fmt.Print(node.Jenkins.Raw.Jobs, " - \n")
		nodeisOnline, _ := node.IsOnline(ctx)
		if nodeisOnline {
			fmt.Println("Node is Online")
		} else {
			fmt.Println("Node is Offline")
		}
	} else {
		fmt.Print("Node Cannot Found")
	}
}

func CheckAllNodeStatus() {
	jenkins, ctx := getJenkins()

	nodes, _ := jenkins.GetAllNodes(ctx)

	for _, node := range nodes {

		// Fetch Node Data
		node.Poll(ctx)
		fmt.Print(node.GetName())
		fmt.Print(node.Jenkins.Raw.Jobs, " - \n")
		fmt.Print(node.Raw.JnlpAgent, " - \n")
		nodeisOnline, _ := node.IsOnline(ctx)
		if nodeisOnline {
			fmt.Println("Node is Online")
		} else {
			fmt.Println("Node is Offline")
		}
	}
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
		if err.Error() == "404" {
			fmt.Printf("Folder Not Found")
		} else if err != nil {
			panic(err)
		}
		folderName = folder.GetName()
	} else {
		folder, err := jenkins.GetFolder(ctx, name, parent)
		if err.Error() == "404" {
			fmt.Printf("Folder Not Found")
		} else if err != nil {
			panic(err)
		} else {
			folderName = folder.GetName()
			fmt.Printf("Folder name: %s", folderName)
		}

	}
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
	config := "<?xml version=\"1.1\" encoding=\"UTF-8\"?><project>\n  <description/>\n  <keepDependencies>false</keepDependencies>\n  <properties/>\n  <scm class=\"hudson.scm.NullSCM\"/>\n  <canRoam>true</canRoam>\n  <disabled>false</disabled>\n  <blockBuildWhenDownstreamBuilding>false</blockBuildWhenDownstreamBuilding>\n  <blockBuildWhenUpstreamBuilding>false</blockBuildWhenUpstreamBuilding>\n  <triggers/>\n  <concurrentBuild>false</concurrentBuild>\n  <builders/>\n  <publishers/>\n  <buildWrappers/>\n</project>"
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

func findJobByName(name string, parent string) *gojenkins.Job {
	jenkins, ctx := getJenkins()
	if parent == "" {
		job, err := jenkins.GetJob(ctx, name)
		if err != nil {
			panic(err)
		}
		return job
	} else {
		job, err := jenkins.GetJob(ctx, name, parent)
		if err != nil {
			panic(err)
		}
		return job
	}
}

func deleteJob(job *gojenkins.Job) {
	_, ctx := getJenkins()
	returnbool, err := job.Delete(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(returnbool)
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
	jenkins := gojenkins.CreateJenkins(nil, "http://3.89.89.181:8080/", "caner", "cnr1")
	// Provide CA certificate if server is using self-signed certificate
	// caCert, _ := ioutil.ReadFile("/tmp/ca.crt")
	// jenkins.Requester.CACert = caCert
	_, err := jenkins.Init(context.Background())
	if err != nil {
		panic("Something Went Wrong")
	}
	return jenkins, ctx
}
