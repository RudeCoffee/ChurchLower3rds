//go:build windows

package main

import (
	"log"
	"os"

	"golang.org/x/sys/windows/svc"
)

type churchService struct{}

func (m *churchService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	// Always change working directory to executable location so kjv.json and web files load correctly
	if exePath, err := os.Executable(); err == nil {
		os.Chdir(filepathDir(exePath))
	}

	go runServer()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			changes <- svc.Status{State: svc.StopPending}
			break loop
		default:
			log.Printf("Unexpected service control request #%d", c.Cmd)
		}
	}

	return
}

func filepathDir(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[:i]
		}
	}
	return "."
}

func runService(name string) {
	err := svc.Run(name, &churchService{})
	if err != nil {
		log.Fatalf("Service %s failed: %v", name, err)
	}
}
