package service_test

import (
	"fmt"

	"github.com/harmonikit/harmoni/service"
)

// userService implements service.Service.
type userService struct{}

func (s *userService) Name() string { return "user-service" }

func ExampleService() {
	svc := &userService{}
	fmt.Println(svc.Name())
	// Output: user-service
}

func ExampleMiddleware() {
	loggingMW := service.Middleware(func(s service.Service) service.Service {
		fmt.Printf("initializing service: %s\n", s.Name())
		return s
	})

	svc := &userService{}
	_ = loggingMW(svc)
	// Output: initializing service: user-service
}
