// +build !golinewrap_debug

package golinewrap

// debug is a no-op for release builds
func debug(_ string, _ ...interface{}) {}
