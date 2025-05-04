//
//  Copyright © 2025 Peter W. Morreale. All Rights Reserved.
//

// Package rest executes REST calls
package rest

import (
	"context"

	"github.com/pwmorreale/rapid/internal/config"
)

// Rest  defines the interface for managing requests and responses
//
//go:generate go tool counterfeiter -o ../../test/mocks/fake_rest.go . Rest
type Rest interface {
	Execute(context.Context, *config.Request)
}

// Context defines a scenario context.
type Context struct {
}

// New creates a new instance.
func New() *Context {
	return &Context{}
}

// Execute creates and executes the request then validates the response.
func (r *Context) Execute(_ context.Context, _ *config.Request) {

}
