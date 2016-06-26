package dto

type ValidRequestDTO interface {
	AlreadySet(dtoName string) bool
	MarkSet(dtoName string)
	MarkAllUnset()
}
