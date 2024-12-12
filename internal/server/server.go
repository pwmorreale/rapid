//
//  Copyright © 2024 Peter W. Morreale. All Rights Reserved.
//

// Package server implements the REST server.
package server

// Server defines interfaces for manipulating the server instance
//
//go:generate counterfeiter -o ../../test/mocks/fake_server.go . Server
type Server interface {
	// Start begins server execution.
	Start() error
	Stop()
}

// Instance implements a server instance.
type Instance struct {
	addr string
}

// New creates a new instance
func New(addr string) *Instance {
	return &Instance{
		addr: addr,
	}
}

// Start starts the REST server.
func (s *Instance) Start() error {
	return nil
}

// Stop stops the server instance.
func (s *Instance) Stop() {
}
