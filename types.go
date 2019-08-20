// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

// +build mips,mipsle

package gosane

// primitive types

type SANE_Byte uint8

// ordered most significant to least significant (big endian)
type SANE_Word struct {
	int32
}

type SANE_String []byte

type SANE_Bool SANE_Word

type SANE_Handle *int32
