// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package osutil

import "syscall"

// DebugSymlinkForTestsOnly creates newname as a symbolic link to oldname.
// If there is an error, it will be of type *LinkError.
func DebugSymlinkForTestsOnly(oldname, newname string) error {
	// CreateSymbolicLink is not supported before Windows Vista
	if syscall.LoadCreateSymbolicLink() != nil {
		return &LinkError{"symlink", oldname, newname, syscall.EWINDOWS}
	}

	// '/' does not work in link's content
	oldname = fromSlash(oldname)

	// need the exact location of the oldname when its relative to determine if its a directory
	destpath := oldname
	if !isAbs(oldname) {
		destpath = dirname(newname) + `\` + oldname
	}

	fi, err := Lstat(destpath)
	isdir := err == nil && fi.IsDir()

	n, err := syscall.UTF16PtrFromString(fixLongPath(newname))
	if err != nil {
		return &LinkError{"symlink", oldname, newname, err}
	}
	o, err := syscall.UTF16PtrFromString(fixLongPath(oldname))
	if err != nil {
		return &LinkError{"symlink", oldname, newname, err}
	}

	var flags uint32
	if isdir {
		flags |= syscall.SYMBOLIC_LINK_FLAG_DIRECTORY
	}
	flags |= 0x02 // SYMBOLIC_LINK_FLAG_ALLOW_UNPRIVILEGED_CREATE
	err = syscall.CreateSymbolicLink(n, o, flags)
	if err != nil {
		return &LinkError{"symlink", oldname, newname, err}
	}
	return nil
}
