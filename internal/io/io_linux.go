// Copyright (c) 2021 Andy Pan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// +build linux

package io

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

// Writev calls writev() on Linux.
func Writev(fd int, iovs [][]byte) (int, error) {
	return unix.Writev(fd, iovs)
}

// Readv calls readv() on Linux.
func Readv(fd int, iovs [][]byte) (int, error) {
	return unix.Readv(fd, iovs)
}

// SendMsgs calls sendmsg() with MSG_ZEROCOPY on Linux sending an iov.
func SendMsgs(fd int, iovs [][]byte) (int, error) {
	iovecs := make([]unix.Iovec, len(iovs))
	for i, iov := range iovs {
		iovecs[i].SetLen(len(iov))
		if len(iov) > 0 {
			iovecs[i].Base = &iov[0]
		} else {
			iovecs[i].Base = (*byte)(unsafe.Pointer(&_zero))
		}
	}
	msg := unix.Msghdr{Iov: &iovecs[0], Iovlen: uint64(len(iovecs))}
	if r, _, errno := unix.Syscall(unix.SYS_SENDMSG, uintptr(fd), uintptr(unsafe.Pointer(&msg)), uintptr(unix.MSG_ZEROCOPY)); errno != 0 {
		return 0, errno
	} else {
		return int(r), nil
	}
}
