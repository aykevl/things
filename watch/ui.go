package main

import (
	"time"

	"github.com/aykevl/tinygl"
	"github.com/aykevl/tinygl/style/basic"
	"tinygo.org/x/drivers/pixel"
)

// ViewManager is a kind of window manager for the watch.
type ViewManager[T pixel.Color] struct {
	screen *tinygl.Screen[T]
	*basic.Basic[T]

	// This is a stack of views that can be added on top and popped from when
	// going back to the previous view.
	views []View[T]
}

// Len returns the number of views.
func (v *ViewManager[T]) Len() int {
	return len(v.views)
}

// Push adds the view to the stack of views, displaying it on top of the screen.
func (v *ViewManager[T]) Push(view View[T]) {
	v.views = append(v.views, view)
	v.screen.SetChild(view.Object)
}

// Pop removes the topmost view, revealing the view underneath.
func (v *ViewManager[T]) Pop() {
	v.views[len(v.views)-1] = View[T]{} // allow this view to be GC'd
	v.views = v.views[:len(v.views)-1]
	v.screen.SetChild(v.views[len(v.views)-1].Object)
}

// Replace all views in the stack, and replace it with the given view.
// This is used to set a new homescreen for example.
func (v *ViewManager[T]) ReplaceAll(view View[T]) {
	v.views = v.views[:0]
	v.views = append(v.views, view)
	v.screen.SetChild(view.Object)
}

// Update runs the Update callback attached to this view.
func (v *ViewManager[T]) Update(now time.Time) {
	callback := v.views[len(v.views)-1].Update
	if callback != nil {
		callback(now)
	}
}

// A view is a single full-screen UI that is active at a time.
// It is comparable to an Android activity.
type View[T pixel.Color] struct {
	tinygl.Object[T]
	Update func(now time.Time)
}

// NewView creates a new view with the given values.
func NewView[T pixel.Color](object tinygl.Object[T], update func(now time.Time)) View[T] {
	return View[T]{
		Object: object,
		Update: update,
	}
}
