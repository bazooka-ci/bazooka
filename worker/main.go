package main

import (
	log "github.com/Sirupsen/logrus"
	bzklog "github.com/bazooka-ci/bazooka/commons/logs"

	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	log.SetFormatter(&bzklog.BzkFormatter{})
	log.SetLevel(log.DebugLevel)
}

func main() {
	context := initContext()
	go heartbeat(context)
	go listenForJobs(context)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	log.Infof("Got SIGTERM")
	context.busy.Wait()
	log.Infof("Exiting")
}

func listenForJobs(context *context) {
	for {
		log.Info("Reserving a job ...")
		job, err := context.reserveJob()
		if err != nil {
			log.Errorf("Error while reserving a job: %v", err)
			continue
		}
		context.startJob(job)
	}
}

func heartbeat(context *context) {
	for _ = range time.Tick(10 * time.Second) {
		log.Debugf("Heartbeat ...")
		context.client.Internal.Heartbeat()
	}
}
