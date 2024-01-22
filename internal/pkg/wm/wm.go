package wm

type WorkspaceManager interface {
	Switch(target string)
	GetCurrentWorkspace() Workspace
	OnChangeWorkspace() <-chan Workspace
}

type Workspace interface {
	Name() string
	IsEmpty() bool
}
