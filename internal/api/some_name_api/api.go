package somenameapi

import "github.com/solumD/auth-test-task/internal/service"

// for server or handler

// API ...
type API struct {
	someService service.SomeService
}

// New returns new API object
func New(someService service.SomeService) *API {
	return &API{
		someService: someService,
	}
}
