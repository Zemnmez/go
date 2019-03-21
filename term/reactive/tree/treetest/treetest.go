package treetest

import (
	"zemn.me/term/reactive/tree"
)

type RecordedError struct {
	C   tree.Component
	Err error
}

type Recorder struct {
	Components []tree.Component
	Errors     []RecordedError
}

func (r *Recorder) Map(c tree.Component) { r.Components = append(r.Components, c) }
func (r *Recorder) Error(c tree.Component, err error) {
	r.Errors = append(r.Errors, RecordedError{c, err})
}

type MountCall struct{ StateController tree.StateController }
type CloseCall struct{}
type ShouldUpdateCall struct{ NewComponent tree.Component }

type StaticComponent struct {
	Id                string
	Children          []tree.Component
	MountCalls        []MountCall
	CloseCalls        []CloseCall
	ShouldUpdateCalls []ShouldUpdateCall
}

func (s *StaticComponent) Mount(sc tree.StateController) {
	s.MountCalls = append(s.MountCalls, MountCall{StateController: sc})
}

func (s *StaticComponent) Close() {
	s.CloseCalls = append(s.CloseCalls, CloseCall{})
}

func (s *StaticComponent) ShouldUpdate(c tree.Component) (_ bool, _ error) {
	s.ShouldUpdateCalls = append(s.ShouldUpdateCalls, ShouldUpdateCall{c})
	return
}
func (c *StaticComponent) Render() ([]tree.Component, error) { return c.Children, nil }
