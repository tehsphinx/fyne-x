package widget

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// SimpleWidget defines an interface for fyne widgets that are simple
// to implement when based on SimpleWidgetBase.
// For more info on how to implement, see the documentation of the
// SimpleWidgetBase.
type SimpleWidget interface {
	fyne.Widget

	Render() ([]fyne.CanvasObject, func(size fyne.Size))
}

// SimpleWidgetBase defines the base for a SimpleWidget implementation.
// To create a new widget base it on SimpleWidgetBase using composition.
// Create a `New` function initialising the widget and make sure to call
// ExtendBaseWidget in the New function. Always use the `New` function to
// create the widget or make sure `ExtendBaseWidget` is called elsewhere.
//
// Overwrite the `Render() (objects []fyne.CanvasObject, layout func(size fyne.Size))`
// function. It returns the (base-)objects needed to render the widgets content,
// as well as a function `layout` responsible for positioning and resizing the
// different objects based on the incoming available space for the widget.
// Try not to define new objects in the `layout` function as they would be
// recreated every time the widget is refreshed.
//
// Other functions defined by the fyne.Widget interface can be overwritten
// and will be used by the SimpleWidgetBase if overwritten.
//
// See ./example/simple_wigdet.go for a bootstraped widget implementation.
type SimpleWidgetBase struct {
	widget.BaseWidget

	propertyLock sync.RWMutex
	impl         SimpleWidget
}

// Render must be overwritten in a widget to create other widgets and
// canvas objects the widget is composed of. It returns a slice of the
// created objects and the layout function.
// The layout fucntion should be used to position and size the objects (widgets
// and canvas objects). New objects should be created in the Render function body
// outside the returned layout function, so they are not re-created
// every time the widget gets refreshed.
func (s *SimpleWidgetBase) Render() (objects []fyne.CanvasObject, layout func(size fyne.Size)) {
	return nil, func(fyne.Size) {}
}

// CreateRenderer implements the Widget interface. It creates a simpleRenderer
// and returns it. No renderer needs to be implemented. If the simpleRenderer
// doesn't do it, SimpleWidget is probably not suitable for the use case.
// Usually this should not be overwritten or called manually.
func (s *SimpleWidgetBase) CreateRenderer() fyne.WidgetRenderer {
	wdgt := s.super()
	objs, layout := wdgt.Render()

	return newSimpleRenderer(wdgt, objs, layout)
}

// SetState sets or changes the state of a widget. A Refresh
// is triggered after the state changes have been applied.
func (s *SimpleWidgetBase) SetState(setState func()) {
	setState()
	s.super().Refresh()
}

// SetStateSafe sets or changes the state of a widget in a safe way. A Refresh
// is triggered after the state changes have been applied.
// The provided sync.Locker should be the same you use for read protection of the
// widget properties.
func (s *SimpleWidgetBase) SetStateSafe(m sync.Locker, setState func()) {
	m.Lock()
	setState()
	m.Unlock()

	s.super().Refresh()
}

// ExtendBaseWidget is used by an extending widget to make use of BaseWidget functionality.
func (s *SimpleWidgetBase) ExtendBaseWidget(wid SimpleWidget) {
	impl := s.getImpl()
	if impl != nil {
		return
	}

	s.BaseWidget.ExtendBaseWidget(wid)

	s.propertyLock.Lock()
	defer s.propertyLock.Unlock()
	s.impl = wid
}

func (s *SimpleWidgetBase) super() SimpleWidget {
	impl := s.getImpl()
	if impl == nil {
		var x interface{} = s
		return x.(SimpleWidget)
	}
	return impl
}

func (s *SimpleWidgetBase) getImpl() SimpleWidget {
	s.propertyLock.RLock()
	impl := s.impl
	s.propertyLock.RUnlock()

	if impl == nil {
		return nil
	}
	return impl
}
