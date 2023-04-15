package main

import (
	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/pixel"
	"github.com/aykevl/tinygl/style"
)

// ViewManager is a kind of window manager for the watch.
type ViewManager[T pixel.Color] struct {
	screen *tinygl.Screen[T]
	scale  style.Scale

	// This is a stack of views that can be added on top and popped from when
	// going back to the previous view.
	views []tinygl.Object[T]
}

// Len returns the number of views.
func (v *ViewManager[T]) Len() int {
	return len(v.views)
}

// Push adds the view to the stack of views, displaying it on top of the screen.
func (v *ViewManager[T]) Push(view tinygl.Object[T]) {
	v.views = append(v.views, view)
	v.screen.SetChild(view)
}

// Pop removes the topmost view, revealing the view underneath.
func (v *ViewManager[T]) Pop() {
	v.views[len(v.views)-1] = nil // allow this view to be GC'd
	v.views = v.views[:len(v.views)-1]
	v.screen.SetChild(v.views[len(v.views)-1])
}
