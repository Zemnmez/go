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

type StaticComponent struct {
	Id       string
	Children []tree.Component
}

func (StaticComponent) Mount(tree.StateController)                    {}
func (StaticComponent) Close()                                        {}
func (StaticComponent) ShouldUpdate(tree.Component) (_ bool, _ error) { return }
func (c StaticComponent) Render() ([]tree.Component, error)           { return c.Children, nil }
