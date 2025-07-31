package main

import (
	"testing"

	"gotest.tools/v3/assert"
)

type filterTestCase struct {
	Src    string
	Result string
}

func TestUtilFunc(t *testing.T) {
	testCases := []filterTestCase{
		{
			Src:    "This is a kerfuffle opinion I need to share with the world",
			Result: "This is a **** opinion I need to share with the world",
		},
	}

	for _, tc := range testCases {
		tc := tc
		assert.Equal(t, CleanInput(tc.Src), tc.Result)
	}
}
