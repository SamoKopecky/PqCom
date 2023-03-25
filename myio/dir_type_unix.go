//go:build !windows
// +build !windows

package myio

const (
	Log    = ".local/state/pqcom"
	Config = ".config/pqcom"
	Cookie = ".cache/pqcom"
)
