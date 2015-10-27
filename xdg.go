// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Package xdg provides an implementation of the XDG Base Directory Specification.
//
// On initialization of this package (happens automatically), the variables
//
//  ConfigHome
//  DataHome
//  CacheHome
//  RuntimeDir
//  ConfigDirs
//  DataDirs
//
// are set to their recommended values.
// Do not change them unless you are absolutely sure of what you are doing.
//
// For more information on the specification, see:
//
//  http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html
//
// This package takes inspiration from github.com/adrg/xdg. Many thanks.
package xdg

import (
	"errors"
	"os"
	"path"
	"strings"
)

var (
	// home is a single base directory of the user's home directory.
	// This directory is defined by the environment variable $HOME.
	//
	// If $HOME is not set, and is required, then this implementation errors
	// out.
	home string

	// ConfigHome is a single base directory relative to which user-specific
	// configuration files should be written. This directory is defined by the
	// environment variable $XDG_CONFIG_HOME.
	//
	// If $XDG_CONFIG_HOME is not set, the default "$HOME/.config" is used.
	ConfigHome string

	// DataHome is a single base directory relative to which user-specific data
	// files should be written. This directory is defined by the environment
	// variable $XDG_DATA_HOME.
	//
	// If $XDG_DATA_HOME is not set, the default "$HOME/.local/share" is used.
	DataHome string

	// CacheHome is a single base directory relative to which user-specific
	// non-essential (cached) data should be written. This directory is defined
	// by the environment variable $XDG_CACHE_HOME.
	//
	// If $XDG_CACHE_HOME is not set, the default "$HOME/.cache" is used.
	CacheHome string

	// RuntimeDir is a single base directory relative to which user-specific
	// runtime files and other file objects should be placed. This directory is
	// defined by the environment variable $XDG_RUNTIME_DIR.
	//
	// The specification has the following to say about $XDG_RUNTIME_DIR:
	//
	//  $XDG_RUNTIME_DIR defines the base directory relative to which
	//  user-specific non-essential runtime files and other file objects (such
	//  as sockets, named pipes, ...) should be stored. The directory MUST be
	//  owned by the user, and he MUST be the only one having read and write
	//  access to it. Its Unix access mode MUST be 0700.
	//
	//  The lifetime of the directory MUST be bound to the user being logged in.
	//  It MUST be created when the user first logs in and if the user fully
	//  logs out the directory MUST be removed. If the user logs in more than
	//  once he should get pointed to the same directory, and it is mandatory
	//  that the directory continues to exist from his first login to his last
	//  logout on the system, and not removed in between. Files in the directory
	//  MUST not survive reboot or a full logout/login cycle.
	//
	//  The directory MUST be on a local file system and not shared with any
	//  other system. The directory MUST by fully-featured by the standards of
	//  the operating system. More specifically, on Unix-like operating systems
	//  AF_UNIX sockets, symbolic links, hard links, proper permissions, file
	//  locking, sparse files, memory mapping, file change notifications,
	//  a reliable hard link count must be supported, and no restrictions on the
	//  file name character set should be imposed. Files in this directory MAY
	//  be subjected to periodic clean-up. To ensure that your files are not
	//  removed, they should have their access time timestamp modified at least
	//  once every 6 hours of monotonic time or the 'sticky' bit should be set
	//  on the file.
	//
	//  If $XDG_RUNTIME_DIR is not set applications should fall back to
	//  a replacement directory with similar capabilities and print a warning
	//  message. Applications should use this directory for communication and
	//  synchronization purposes and should not place larger files in it, since
	//  it might reside in runtime memory and cannot necessarily be swapped out
	//  to disk.
	//
	// In this implementation, we assume that the system takes care of removing
	// the XDG runtime directory at shutdown.
	//
	// If $XDG_RUNTIME_DIR is not set, this implementation fails FOR NOW.
	RuntimeDir string

	// ConfigDirs is a set of preference ordered base directories relative to
	// which configuration files should be searched. This set of directories is
	// defined by the environment variable $XDG_CONFIG_DIRS. The directories in
	// $XDG_CONFIG_DIRS should be seperated with a colon ':'.
	//
	// If $XDG_CONFIG_DIRS is not set, the default "/etc/xdg" is used.
	ConfigDirs []string

	// DataDirs is a set of preference ordered base directories relative to
	// which data files should be searched. This set of directories is defined
	// by the environment variable $XDG_DATA_DIRS.
	//
	// If $XDG_CONFIG_DIRS is not set, the default "/usr/local/share:/usr/share"
	// is used.
	DataDirs []string
)

// Errors contains all errors that occurred during initialization.
var Errors []error

var ErrHomeInvalid = errors.New("environment variable HOME is invalid or not set")

func init() {
	home = os.Getenv("HOME")
	if path.IsAbs(home) {
		home = ""
		Errors = append(Errors, ErrHomeInvalid)
	}

	ConfigHome = xdgPath("XDG_CONFIG_HOME", "$HOME/.config")
	DataHome = xdgPath("XDG_DATA_HOME", "$HOME/.config")
	CacheHome = xdgPath("XDG_CACHE_HOME", "$HOME/.config")
	RuntimeDir = xdgPath("XDG_RUNTIME_DIR", "")
	ConfigDirs = xdgPaths("XDG_CONFIG_DIRS", "/etc/xdg")
	DataDirs = xdgPaths("XDG_DATA_DIRS", "/usr/local/share:/usr/share")
}

func xdgPath(env, def string) string {
	x := os.Getenv(env)

	if x == "" {
		if strings.Contains(def, "$HOME") {
			if home != "" {
				x = strings.Replace(def, "$HOME", home, -1)
			}
		} else {
			x = def
		}
	}

	// The XDG specification states:
	//
	//  All paths set in these environment variables must be absolute. If an
	//  implementation encounters a relative path in any of these variables it
	//  should consider the path invalid and ignore it.
	if path.IsAbs(x) {
		return x
	}
	Errors = append(Errors, errors.New("no value set for "+env))
	return ""
}

func xdgPath(env, def string) []string {
	xs := os.Getenv(env)

	if xs == "" {
		xs = def
	}

	var fs []string
	for _, x := range strings.Split(xs, ":") {
		// See comment in xdgPath.
		if path.IsAbs(x) {
			fs = append(fs, x)
		} else {
			Errors = append(Errors, errors.New("ignoring "+env+" path element: "+x))
		}
	}
	return fs
}

func OpenConfigFile(suffix string) (*os.File, error) {
	return nil, nil
}

func OpenDataFile(suffix string) (*os.File, error) {
	return nil, nil
}

func OpenCacheFile(suffix string) (*os.File, error) {
	return nil, nil
}

func OpenRuntimeFile(suffix string) (*os.File, error) {
	return nil, nil
}

func FindConfigFiles(suffix string) []string {
	return nil
}

func FindDataFiles(suffix string) []string {
	return nil
}

func FindConfigFile(suffix string) string {
	return ""
}

func FindDataFile(suffix string) string {
	return ""
}

func FindCacheFile(suffix string) string {
	return ""
}

func FindRuntimeFile(suffix string) string {
	return ""
}

// MergeFunc is given to the Merge*Files functions to handle the files that it
// finds. It receives an absolute path to a file, which MergeFunc can then try
// to open. When MergeFunc is done with the file (for example, it couldn't read
// the file, or it was empty) then it can return nil. If an error is returned,
// then the Merge*Files function aborts and returns this error. If an error
// hasn't occurred, but no files need be further inspected, Skip can be returned.
type MergeFunc func(string) error

// Skip can be returned by a MergeFunc which causes the Merge*Files functions
// to skip the rest of the files to be merged.
var Skip = errors.New("skip the rest of the files to be merged")

func MergeDataFiles(suffix string, f MergeFunc) error {
	return mergeFiles(suffix, f, DataHome, DataDirs...)
}

func MergeDataFilesReverse(suffix string, f MergeFunc) error {
	return mergeFilesReverse(suffix, f, DataHome, DataDirs...)
}

func MergeConfigFiles(suffix string, f MergeFunc) error {
	return mergeFiles(suffix, f, ConfigHome, ConfigDirs...)
}

func MergeConfigFilesReverse(suffix string, f MergeFunc) error {
	return mergeFilesReverse(suffix, f, ConfigHome, ConfigDirs...)
}

func mergeFilesReverse(suffix string, f MergeFunc, paths ...string) error {
	var err error
	for _, s := range findFilesReverse(suffix, paths...) {
		if err = f(s); err != nil {
			break
		}
	}
	if err == Skip {
		return nil
	}
	return err
}

func mergeFilesReverse(suffix string, f MergeFunc, paths ...string) error {
	var err error
	for _, s := range findFiles(suffix, paths...) {
		if err = f(s); err != nil {
			break
		}
	}
	if err == Skip {
		return nil
	}
	return err
}
