package server

import "testing"

func TestInvalidBaseConfig(t *testing.T) {
	// should not run
}

func TestMailMissingConfig(t *testing.T) {
	// /contact should return 404
}

func TestInvalidMailAddress(t *testing.T) {
	// Should return BadRequest
}
