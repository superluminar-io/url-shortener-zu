package main

import (
	"context"
	"testing"
)

func TestHandle(t *testing.T) {
	res, err := handle(context.Background(), event{ShouldFail: false})

	if err != nil {
		t.Error("Unexpected error")
	}

	if res != "Done" {
		t.Error("Unexpected response")
	}
}

func TestHandleError(t *testing.T) {
	res, err := handle(context.Background(), event{ShouldFail: true})

	if err == nil {
		t.Error("Expected error")
	}

	if err.Error() != "Error" {
		t.Error("Unexpected error message")
	}

	if res != "" {
		t.Error("Unexpected response")
	}
}
