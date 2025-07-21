//go:build windows

package main

import (
	"log"
	"os"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// isWindowsService returns true if running under the Windows service manager.
func isWindowsService() (bool, error) {
	return svc.IsWindowsService()
}

type service struct{}

func (m *service) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (bool, uint32) {
	const accepted = svc.AcceptStop | svc.AcceptShutdown
	s <- svc.Status{State: svc.StartPending}
	go func() {
		if err := runServer(); err != nil {
			log.Printf("server error: %v", err)
		}
	}()
	s <- svc.Status{State: svc.Running, Accepts: accepted}
	for {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			s <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			s <- svc.Status{State: svc.StopPending}
			return false, 0
		default:
		}
	}
}

func runService(name string) error {
	return svc.Run(name, &service{})
}

func installService(name, desc string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return nil
	}
	s, err = m.CreateService(name, exe, mgr.Config{DisplayName: name, Description: desc})
	if err != nil {
		return err
	}
	defer s.Close()
	return nil
}

func removeService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return err
	}
	defer s.Close()
	return s.Delete()
}
