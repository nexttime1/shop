package enum

type RoleType int8

const (
	AdminRole RoleType = 1
	UserRole  RoleType = 2
	Viewer    RoleType = 3
)
