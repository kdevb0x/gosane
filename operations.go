// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package gosane

const (
	MaxUsernameLen SInt = 128
	MaxPasswordLen SInt = 128
)

type AuthorizationCallback func(resource SStringConst, username SChar, password SChar)

func Init(verionCode SInt, authorize AuthorizationCallback) error {
	return SStatus(Good)
}
