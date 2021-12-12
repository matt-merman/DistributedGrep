package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {

	// Setup our test.
	api := new(API)
	input := Input{"test.txt", "che"}
	var reply *string

	// Perform our operation.
	err := api.Mapper(input, reply)

	// Perform our validation
	assert.NoError(t, err)
	assert.Equal(t, 6, reply, "Did not multiply correctly!`")
}
