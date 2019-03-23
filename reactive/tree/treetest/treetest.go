package treetest

import (
	"fmt"
	"reflect"

	"zemn.me/debug"
	"zemn.me/reactive/tree"
)

type RecordedError struct {
	C   tree.Component
	Err error
}

type Recorder struct {
	Components       []tree.Component
	ClosedComponents []tree.Component
	Errors           []RecordedError
}

func (r *Recorder) Clear() { *r = Recorder{} }

func (r *Recorder) UnMap(c tree.Component) { r.ClosedComponents = append(r.ClosedComponents, c) }
func (r *Recorder) Map(c tree.Component) {
	debug.Log("mapping %+v: now mapped %d components", reflect.ValueOf(c), len(r.Components)+1)

	r.Components = append(r.Components, c)
}
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
	RenderCalls       []bool
}

func (s *StaticComponent) ForceUpdate() (err error) {
	if len(s.MountCalls) < 1 {
		return fmt.Errorf(
			"trying to force update on %[0]s, but"+
				" %[0]s has no record of being mounted!",
			s.Id,
		)
	}

	s.MountCalls[len(s.MountCalls)-1].StateController.Update()

	return nil
}
func (s StaticComponent) Name() string { return fmt.Sprintf("%s<%s>", s.Id, reflect.TypeOf(s)) }
func (s *StaticComponent) Clear()      { *s = StaticComponent{} }

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
func (c *StaticComponent) Render() ([]tree.Component, error) {
	c.RenderCalls = append(c.RenderCalls, true)
	return c.Children, nil
}
