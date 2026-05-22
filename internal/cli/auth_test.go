// Copyright 2026 sidhartha1s. Licensed under Apache-2.0. See LICENSE.

package cli

import (
	"strings"
	"testing"
)

// Test that the auth command structure is properly set up.
func TestAuthCmd_HasSubcommands(t *testing.T) {
	flags := &rootFlags{}
	cmd := newAuthCmd(flags)

	subcommands := cmd.Commands()
	if len(subcommands) < 3 {
		t.Errorf("expected at least 3 subcommands (status, set-token, logout), got %d", len(subcommands))
	}

	names := make(map[string]bool)
	for _, sub := range subcommands {
		names[sub.Name()] = true
	}

	for _, expected := range []string{"status", "set-token", "logout"} {
		if !names[expected] {
			t.Errorf("missing expected subcommand: %s", expected)
		}
	}
}

// Test that auth set-token is configured to require exactly one argument.
func TestAuthSetTokenCmd_RequiresExactlyOneArg(t *testing.T) {
	flags := &rootFlags{}
	cmd := newAuthCmd(flags)

	// Find set-token subcommand
	var found bool
	for _, sub := range cmd.Commands() {
		if sub.Name() == "set-token" {
			found = true
			// Verify Args is set (cobra.ExactArgs returns a PositionalArgs func)
			if sub.Args == nil {
				t.Error("set-token should have Args validator set")
			}
			// Verify the Use string indicates a required argument
			if !strings.Contains(sub.Use, "<token>") {
				t.Errorf("set-token Use should indicate required token arg, got: %s", sub.Use)
			}
			break
		}
	}

	if !found {
		t.Fatal("set-token subcommand not found")
	}
}

// Test that auth status command exists and has correct short description.
func TestAuthStatusCmd_Description(t *testing.T) {
	flags := &rootFlags{}
	cmd := newAuthCmd(flags)

	var found bool
	for _, sub := range cmd.Commands() {
		if sub.Name() == "status" {
			found = true
			if !strings.Contains(sub.Short, "authentication") {
				t.Errorf("status command short description should mention authentication, got: %s", sub.Short)
			}
			break
		}
	}

	if !found {
		t.Fatal("status subcommand not found")
	}
}
