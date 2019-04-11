package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type event struct {
	ShouldFail bool `json:"should_fail"`
}

func handle(ctx context.Context, e event) (string, error) {
	fmt.Printf(
		"Function of %s-%s invoked",
		os.Getenv("PREFIX"),
		os.Getenv("PROJECT"),
	)

	if e.ShouldFail == true {
		return "", errors.New("Error")
	}

	return "Done", nil
}
