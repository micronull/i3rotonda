package wm

type WorkspaceManager interface {
	Switch(target string)
	GetCurrentWorkspace() Workspace
	OnChangeWorkspace() <-chan Workspace
}

type Workspace interface {
	GetName() string
	IsEmpty() bool
}
