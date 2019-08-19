package sane

// primitive types

type SANE_Byte uint8

// ordered most significant to least significant (big endian)
type SANE_Word [4]SANE_Byte

type SANE_String []byte

type SANE_Handle *int32
