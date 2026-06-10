package cmd

import "testing"

func TestRootCommandRegistersIgnoreArrayOrderFlag(t *testing.T) {
	flag := rootCmd.Flags().Lookup("ignore-array-order")
	if flag == nil {
		t.Fatal("expected --ignore-array-order flag")
	}
	if flag.DefValue != "false" {
		t.Fatalf("expected default false, got %q", flag.DefValue)
	}
}
