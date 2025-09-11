package models

type Status int

const (
	Inactive Status = iota
	Active
)

func (s Status) String() string {
	switch s {
	case Inactive:
		return "Inactive"
	case Active:
		return "Active"
	default:
		return "Unknown"
	}
}
