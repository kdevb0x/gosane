// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

// remove the ``
// present so linter and autocomplete will work while developing on linux.
// `+build mips mipsle`

package gosane

import (
	"unsafe"
)

// primitive types
// Type names from the spec of the form `SANE_*` have been shortened to `S*`.

type SByte uint8

// ordered most significant to least significant (big endian)
type SWord int32

type SString []byte

type SStringConst string

type SBool = SWord

const (
	SFALSE SBool = iota
	STRUE
)

type SInt = SWord

type SChar byte

// defined as `typedef void *SHandle;`
type SHandle unsafe.Pointer

type SStatus uint

const (
	Good SStatus = iota
	Unsupported
	Cancelled
	DeviceBusy
	Inval
	Eof
	Jammed
	NoDocs
	CoverOpen
	IoError
	NoMem
	AccessDenied
)

func (e SStatus) Error() string {
	var errString string
	switch e {

	case Unsupported:
		errString = "Operation is not supported"
	case Cancelled:
		errString = "Operation was cancelled"
	case DeviceBusy:
		errString = "Device is busy---retry later"
	case Inval:
		errString = "Data or argument is invalid"
	case Eof:
		errString = "No more data available (end-of-file)"
	case Jammed:
		errString = "Document feeder jammed"
	case NoDocs:
		errString = "Document feeder out of documents"
	case CoverOpen:
		errString = "Scanner cover is open"
	case IoError:
		errString = "Error during device I/O"
	case NoMem:
		errString = "Out of memory"
	case AccessDenied:
		errString = "Access to resource has been denied"
	}
	return errString
}

type Device struct {
	Name SStringConst

	// Vendor should be one of:
	// AGFA 		Microtek
	// Abaton 		Minolta
	// Acer 		Mustek
	// Apple 		NEC
	// Artec 		Nikon
	// Avision 		Plustek
	// CANON 		Polaroid
	// Connectix 		Ricoh
	// Epson 		Sharp
	// Fujitsu 		Siemens
	// Hewlett-Packard  	Tamarack
	// IBM 			UMAX
	// Kodak 		Noname
	// Logitech
	Vendor SStringConst

	Model SStringConst

	// DeviceType should be one of:
	// 	"film scanner"
	// 	"flatbed scanner"
	// 	"frame grabber"
	// 	"handheld scanner"
	// 	"multi-function peripheral"
	// 	"sheetfed scanner"
	// 	"still camera"
	// 	"video camera"
	// 	"virtual device"
	DeviceType SStringConst // named `type` in SANE spec.

	// The name and the email address of the backend author,
	// or a contact person in the format:
	//
	//    `Firstname Lastname <name@domain.org>`
	BackendAuthorEmail SStringConst

	// Backend_website should be set up by the backend with the website,
	// or ftp address of the backend in the format:
	//
	//     `http://www.domain.org/sane-hello/index.html`
	BackendWebsite SStringConst

	// Text that describes where a user can find this device.
	// The text should be configurable by the administrator.
	// This could e.g. look like:
	//
	//    `building 93, 2nd plane, room 2124`
	DeviceLocation SStringConst

	// Comment can be used to display any comment about the device to the user.
	// The text should be configurable by the administrator.
	Comment SStringConst

	ReservedString        SStringConst
	BackendVersionCode    SInt
	BackendCapablityFlags SInt
	ReservedInt           SInt
}

type OptionDescriptor struct {
	// Name uniquely identifies the option.
	Name SStringConst
	// Title is used by frontend as title string.
	Title SStringConst
	// A (potentially very) long string used as help text to describe option.
	Desc SStringConst
	// Type is the type of the OptionValue
	Type ValueType
	// Unit specifies the physical unit of the OptionValue.
	Unit ValueUnit
	// Size of the OptionValue (in bytes)
	// Has a slightly different interpretation depending on the type of the option value:

	//    SString: 	The size is the maximum size of the string. For the
	// purpose of string size calcuations, the terminating NUL character is
	// considered to be part of the string.
	//
	//  	SInt, SFixed:
	// The size must be a positive integer multiple of the size of an SWord.
	// The option value is a vector of length.
	//
	// 		`OptionValue.Size / Sizeof(SWord)`
	//
	//
	// 	SBool:
	// The size must be set to Sizeof(SWord).
	//
	// 	TypeButton, TypeGroup:
	// The option size is ignored.

	Size           SInt
	Cap            Capabilities
	ConstraintType ConstraintType
	// This is a C union in the spec.
	Constraint struct {
		StringList []SStringConst
		WordList   []SWord
		Range      *SRange
	}
}

type Capabilities SInt

const (
	// The option value can be set by a call to ControlOption().
	SoftSelect Capabilities = 0x01

	// The option value can be set by user-intervention (e.g., by flipping a switch).
	// The user-interface should prompt the user to execute the appropriate
	// action to set such an option.
	// This capability is mutually exclusive with SoftSelect (either one of
	// them can be set, but not both simultaneously).
	HardSelect Capabilities = SoftSelect >> 1

	// The option value can be detected by software. If SoftSelect is set,
	// this capability must be set. If HardSelect is set, this capability
	// may or may not be set. If this capability is set but neither SoftSelect
	// nor HardSelect are, then there is no way to control the option.
	// That is, the option provides read-out of the current value only.
	SoftDetect Capabilities = HardSelect >> 1

	// If set, this capability indicates that an option is not directly
	// supported by the device and is instead emulated in the backend. A
	// sophisticated frontend may elect to use its own (presumably better)
	// emulation in lieu of an emulated option.
	Emulated Capabilities = SoftDetect >> 1

	// If set, this capability indicates that the backend (or the device) is
	// capable to picking a reasonable option value automatically. For such
	// options, it is possible to select automatic operation by calling
	// ControlOption() with an action value of ActionSetAuto.
	Automatic Capabilities = Emulated >> 1

	// If set, this capability indicates that the option is not currently
	// active (e.g., because it's meaningful only if another option is set
	// to some other value).
	Inactive Capabilities = Automatic >> 1

	// If set, this capability indicates that the option should be considered
	// an ``advanced user option''. If this capability is set for an option
	// of type TypeGroup, all options belonging to the group are also advanced,
	// even if they don't set the capabilty themselves. A frontend typically
	// displays such options in a less conspicuous way than regular options
	// (e.g., a command line interface may list such options last or a
	// graphical interface may make them available in a seperate ``advanced
	// settings'' dialog).
	Advanced Capabilities = Inactive >> 1

	// If set, this capability indicates that the option shouldn't be displayed
	// to and used by the user directly. Instead a hidden option is supposed
	// to be automatically used by the frontend, like e.g. the ``preview''
	// option. If this capability is set for an option of type TypeGroup,
	// all options belonging to the group are also hidden, even if they don't
	// set the capabilty themselves. A frontend typically doesn't display
	// such options by default but there should be a way to override this
	// default behaviour.
	Hidden Capabilities = Advanced >> 1

	// If set, this capability indicates that the option may be at any time
	// between Open() and Close(). I.e. it's allowed to set it even while an
	// image is acquired.
	AlwaysSettable Capabilities = Hidden >> 1
)

type ConstraintType SInt

const (
	None ConstraintType = iota
	Range
	WordList
	StringList
)

type SRange struct {
	Min   SWord
	Max   SWord
	Quant SWord
}

type ValueType uint

const (
	// Option value is of type SBool.
	TypeBool ValueType = iota

	// Option value is of type SANE_Int.
	TypeInt
	// Option value is of type SANE_Fixed.
	TypeFixed

	// Option value is of type SANE_String.
	TypeString

	// An option of this type has no value. Instead, setting an option of
	// this type has an option-specific side-effect. For example, a button-
	// typed option could be used by a backend to provide a means to select
	// default values or to the tell an automatic document feeder to advance
	// to the next sheet of paper.
	TypeButton

	// An option of this type has no value. This type is used to group logically
	// related options. A group option is in effect up to the point where
	// another group option is encountered (or up to the end of the option list,
	// if there are no other group options). For group options, only members
	// title and type are valid in the option descriptor.
	TypeGroup
)

type ValueUnit uint

const (
	// Value is unit-less (e.g., page count).
	UnitNone ValueUnit = iota

	// Value is in number of pixels.
	UnitPixel

	// Value is in number of bits.
	UnitBit

	// Value is in millimeters.
	UnitMm

	// Value is a resolution in dots/inch.
	UnitDpi

	// Value is a percentage.
	UnitPercent

	// Value is time in µ-seconds.
	UnitMicrosecond
)
