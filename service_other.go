//go:build !windows

package main

func runService(name string) {
	// No-op on non-windows
}
