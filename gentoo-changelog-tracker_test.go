package main

import (
    "testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
    query := "portage/eix"
	expected := "app-portage/eix"
	answer := lookupAtomEix(query)

	if answer != expected {
        t.Fatalf(`Got %q from query %q but should find %q`, answer, query, expected)
    }
}
