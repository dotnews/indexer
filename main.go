package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/inthenews/indexer/config"
	"github.com/inthenews/indexer/worker"
)

const (
	root          = "./"
	processFlag   = "process"
	processUsage  = "Processing method: [index|delete] defaults to index"
	indexProcess  = "index"
	deleteProcess = "delete"
)

func main() {
	var err error
	c := config.New(root)
	worker := worker.New(c)
	process := flag.String(processFlag, indexProcess, processUsage)
	flag.Parse()

	glog.Infof("Running %s on %s", *process, c.Env)

	switch *process {
	case indexProcess:
		err = worker.Index()
	case deleteProcess:
		err = worker.Delete()
	default:
		glog.Fatalf("Invalid process argument: %s", *process)
	}

	if err != nil {
		glog.Fatalf("Failed running worker: %v", err)
	}
}
