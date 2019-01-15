package conv

import "fmt"

type UnsupportedTypeError struct {
	Type string
	For  string
}

func (a *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("Type %v is not supported for %v", a.Type, a.For)
}

type NotPointerError struct {
	Type string
}

func (a *NotPointerError) Error() string {
	return fmt.Sprintf("Type %v must be a pointer of this type", a.Type)
}

type BadTypeError struct {
	WantedType string
	GotType    string
}

func (a *BadTypeError) Error() string {
	return fmt.Sprintf("Error interface{} must be %v but is %v", a.WantedType, a.GotType)
}

type KeyStatus string

const (
	NotFound   KeyStatus = "not found"
	AlreadySet KeyStatus = "already set"
)

type KeyError struct {
	Name   string
	Status KeyStatus
}

func (a *KeyError) Error() string {
	return fmt.Sprintf("Error key %v status %v", a.Name, a.Status)
}
