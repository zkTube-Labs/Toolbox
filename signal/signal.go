package signal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type listener struct {
	Addr     string `json:"addr"`
	FD       int    `json:"fd"`
	Filename string `json:"filename"`
}

type signalType struct{}

func (s *signalType) importListener(addr string) (net.Listener, error) {
	// Extract the encoded listener metadata from the environment.
	listenerEnv := os.Getenv("LISTENER")
	if listenerEnv == "" {
		return nil, fmt.Errorf("unable to find LISTENER environment variable")
	}

	// Unmarshal the listener metadata.
	var l listener
	err := json.Unmarshal([]byte(listenerEnv), &l)
	if err != nil {
		return nil, err
	}
	if l.Addr != addr {
		return nil, fmt.Errorf("unable to find listener for %v", addr)
	}

	// The file has already been passed to this process, extract the file
	// descriptor and name from the metadata to rebuild/find the *os.File for
	// the listener.
	listenerFile := os.NewFile(uintptr(l.FD), l.Filename)
	if listenerFile == nil {
		return nil, fmt.Errorf("unable to create listener file: %v", err)
	}
	defer listenerFile.Close()

	// Create a net.Listener from the *os.File.
	ln, err := net.FileListener(listenerFile)
	if err != nil {
		return nil, err
	}

	return ln, nil
}

func (s *signalType) createListener(addr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return ln, nil
}

func (s *signalType) createOrImportListener(addr string) (net.Listener, error) {
	// Try and import a listener for addr. If it's found, use it.
	ln, err := s.importListener(addr)
	if err == nil {
		fmt.Printf("Imported listener file descriptor for %v.\n", addr)
		return ln, nil
	}

	// No listener was imported, that means this process has to create one.
	ln, err = s.createListener(addr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created listener file descriptor for %v.\n", addr)
	return ln, nil
}

func (s *signalType) getListenerFile(ln net.Listener) (*os.File, error) {
	switch t := ln.(type) {
	case *net.TCPListener:
		return t.File()
	case *net.UnixListener:
		return t.File()
	}
	return nil, fmt.Errorf("unsupported listener: %T", ln)
}

func (s *signalType) forkChild(addr string, ln net.Listener) (*os.Process, error) {
	// Get the file descriptor for the listener and marshal the metadata to pass
	// to the child in the environment.
	lnFile, err := s.getListenerFile(ln)
	if err != nil {
		return nil, err
	}
	defer lnFile.Close()
	l := listener{
		Addr:     addr,
		FD:       3,
		Filename: lnFile.Name(),
	}
	listenerEnv, err := json.Marshal(l)
	if err != nil {
		return nil, err
	}

	// Pass stdin, stdout, and stderr along with the listener to the child.
	files := []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
		lnFile,
	}

	// Get current environment and add in the listener to it.
	environment := append(os.Environ(), "LISTENER="+string(listenerEnv))

	// Get current process name and directory.
	execName, err := os.Executable()
	if err != nil {
		return nil, err
	}
	execDir := filepath.Dir(execName)

	// Spawn child process.
	p, err := os.StartProcess(execName, []string{execName}, &os.ProcAttr{
		Dir:   execDir,
		Env:   environment,
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *signalType) waitForSignals(addr string, ln net.Listener, server *http.Server) error {
	signalCh := make(chan os.Signal, 1024)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case signalValue := <-signalCh:
			fmt.Printf("%v signal received.\n", s)
			switch signalValue {
			case syscall.SIGHUP:
				// Fork a child process.
				p, err := s.forkChild(addr, ln)
				if err != nil {
					fmt.Printf("Unable to fork child: %v.\n", err)
					continue
				}
				fmt.Printf("Forked child %v.\n", p.Pid)

				// Create a context that will expire in 5 seconds and use this as a
				// timeout to Shutdown.
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				// Return any errors during shutdown.
				return server.Shutdown(ctx)
			case syscall.SIGINT, syscall.SIGQUIT:
				// Create a context that will expire in 5 seconds and use this as a
				// timeout to Shutdown.
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				// Return any errors during shutdown.
				return server.Shutdown(ctx)
			}
		}
	}
}

func (s *signalType) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from %v!\n", os.Getpid())
}

func (s *signalType) startServer(addr string, ln net.Listener, handle http.Handler) *http.Server {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handle,
	}
	go httpServer.Serve(ln)

	return httpServer
}

// listen and signal
func Listen(addr string, handle http.Handler) error {
	// Parse command line flags for the address to listen on.
	var s signalType
	// Create (or import) a net.Listener and start a goroutine that runs
	// a HTTP server on that net.Listener.
	ln, err := s.createOrImportListener(addr)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to create or import a listener: %v.\n", err))
	}
	server := s.startServer(addr, ln, handle)

	// Wait for signals to either fork or quit.
	err = s.waitForSignals(addr, ln, server)
	if err != nil {
		return errors.New(fmt.Sprintf("Exiting: %v\n", err))
	}
	return nil
}
