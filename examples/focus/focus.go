package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	// DefaultFormFieldWidth is the default field screen width of form elements
	// whose field width is flexible (0). This is used in the Form class for
	// horizontal layouts.
	DefaultFormFieldWidth = 10

	// DefaultFormFieldHeight is the default field height of multi-line form
	// elements whose field height is flexible (0).
	DefaultFormFieldHeight = 5
)

// RootContainer allows you to combine multiple one-line form elements into a vertical
// or horizontal layout. RootContainer elements include types such as InputField or
// Checkbox. These elements can be optionally followed by one or more buttons
// for which you can define form-wide actions (e.g. Save, Clear, Cancel).
//
// See https://github.com/rivo/tview/wiki/RootContainer for an example.
type RootContainer struct {
	*tview.Box

	// The items of the form (one row per item).
	items []tview.FormItem

	// If set to true, instead of position items and buttons from top to bottom,
	// they are positioned from left to right.
	horizontal bool

	// The alignment of the buttons.
	buttonsAlign int

	// The number of empty cells between items.
	itemPadding int

	// The index of the item or button which has focus. (Items are counted first,
	// buttons are counted last.) This is only used when the form itself receives
	// focus so that the last element that had focus keeps it.
	focusedElement int

	// The label color.
	labelColor tcell.Color

	// The background color of the input area.
	fieldBackgroundColor tcell.Color

	// The text color of the input area.
	fieldTextColor tcell.Color

	// The style of the buttons when they are not focused.
	buttonStyle tcell.Style

	// The style of the buttons when they are focused.
	buttonActivatedStyle tcell.Style

	// The style of the buttons when they are disabled.
	buttonDisabledStyle tcell.Style

	// The last (valid) key that wsa sent to a "finished" handler or -1 if no
	// such key is known yet.
	lastFinishedKey tcell.Key

	// An optional function which is called when the user hits Escape.
	cancel func()
}

// NewRootContainer returns a new form.
func NewRootContainer() *RootContainer {
	box := tview.NewBox().SetBorderPadding(1, 1, 1, 1)

	f := &RootContainer{
		Box:                  box,
		itemPadding:          1,
		labelColor:           tview.Styles.SecondaryTextColor,
		fieldBackgroundColor: tview.Styles.ContrastBackgroundColor,
		fieldTextColor:       tview.Styles.PrimaryTextColor,
		buttonStyle:          tcell.StyleDefault.Background(tview.Styles.ContrastBackgroundColor).Foreground(tview.Styles.PrimaryTextColor),
		buttonActivatedStyle: tcell.StyleDefault.Background(tview.Styles.PrimaryTextColor).Foreground(tview.Styles.ContrastBackgroundColor),
		buttonDisabledStyle:  tcell.StyleDefault.Background(tview.Styles.ContrastBackgroundColor).Foreground(tview.Styles.ContrastSecondaryTextColor),
		lastFinishedKey:      tcell.KeyTab, // To skip over inactive elements at the beginning of the form.
	}

	return f
}

// SetFocus shifts the focus to the form element with the given index, counting
// non-button items first and buttons last. Note that this index is only used
// when the form itself receives focus.
func (f *RootContainer) SetFocus(index int) *RootContainer {
	if index < 0 {
		f.focusedElement = 0
	} else if index >= len(f.items) {
		f.focusedElement = len(f.items)
	} else {
		f.focusedElement = index
	}
	return f
}

// AddFormItem adds a new item to the form. This can be used to add your own
// objects to the form. Note, however, that the Form class will override some
// of its attributes to make it work in the form context. Specifically, these
// are:
//
//   - The label width
//   - The label color
//   - The background color
//   - The field text color
//   - The field background color
func (f *RootContainer) AddFormItem(item tview.FormItem) *RootContainer {
	f.items = append(f.items, item)
	return f
}

// AddTextView adds a text view to the form. It has a label and text, a size
// (width and height) referring to the actual text element (a fieldWidth of 0
// extends it as far right as possible, a fieldHeight of 0 will cause it to be
// [DefaultFormFieldHeight]), a flag to turn on/off dynamic colors, and a flag
// to turn on/off scrolling. If scrolling is turned off, the text view will not
// receive focus.
func (f *RootContainer) AddTextView(label, text string, fieldWidth, fieldHeight int, dynamicColors, scrollable bool) *RootContainer {
	if fieldHeight == 0 {
		fieldHeight = DefaultFormFieldHeight
	}
	textArea := tview.NewTextView().
		SetLabel(label).
		SetSize(fieldHeight, fieldWidth).
		SetDynamicColors(dynamicColors).
		SetScrollable(scrollable).
		SetText(text)
	f.items = append(f.items, textArea)
	return f
}

// GetFormItemCount returns the number of items in the form (not including the
// buttons).
func (f *RootContainer) GetFormItemCount() int {
	return len(f.items)
}

// GetFormItem returns the form item at the given position, starting with index
// 0. Elements are referenced in the order they were added. Buttons are not
// included.
func (f *RootContainer) GetFormItem(index int) tview.FormItem {
	return f.items[index]
}

// RemoveFormItem removes the form element at the given position, starting with
// index 0. Elements are referenced in the order they were added. Buttons are
// not included.
func (f *RootContainer) RemoveFormItem(index int) *RootContainer {
	f.items = append(f.items[:index], f.items[index+1:]...)
	return f
}

// GetFormItemByLabel returns the first form element with the given label. If
// no such element is found, nil is returned. Buttons are not searched and will
// therefore not be returned.
func (f *RootContainer) GetFormItemByLabel(label string) tview.FormItem {
	for _, item := range f.items {
		if item.GetLabel() == label {
			return item
		}
	}
	return nil
}

// GetFormItemIndex returns the index of the first form element with the given
// label. If no such element is found, -1 is returned. Buttons are not searched
// and will therefore not be returned.
func (f *RootContainer) GetFormItemIndex(label string) int {
	for index, item := range f.items {
		if item.GetLabel() == label {
			return index
		}
	}
	return -1
}

// GetFocusedItemIndex returns the indices of the form element or button which
// currently has focus. If they don't, -1 is returned resepectively.
func (f *RootContainer) GetFocusedItemIndex() (formItem, button int) {
	index := f.focusIndex()
	if index < 0 {
		return -1, -1
	}
	if index < len(f.items) {
		return index, -1
	}
	return -1, index - len(f.items)
}

// SetCancelFunc sets a handler which is called when the user hits the Escape
// key.
func (f *RootContainer) SetCancelFunc(callback func()) *RootContainer {
	f.cancel = callback
	return f
}

// Draw draws this primitive onto the screen.
func (f *RootContainer) Draw(screen tcell.Screen) {
	f.Box.DrawForSubclass(screen, f)

	// Determine the actual item that has focus.
	if index := f.focusIndex(); index >= 0 {
		f.focusedElement = index
	}

	// Determine the dimensions.
	x, y, width, height := f.GetInnerRect()
	topLimit := y
	bottomLimit := y + height
	rightLimit := x + width
	startX := x

	// Find the longest label.
	var maxLabelWidth int
	for _, item := range f.items {
		labelWidth := tview.TaggedStringWidth(item.GetLabel())
		if labelWidth > maxLabelWidth {
			maxLabelWidth = labelWidth
		}
	}
	maxLabelWidth++ // Add one space.

	// Calculate positions of form items.
	type position struct{ x, y, width, height int }
	positions := make([]position, len(f.items))
	var (
		focusedPosition position
		lineHeight      = 1
	)
	for index, item := range f.items {
		// Calculate the space needed.
		labelWidth := tview.TaggedStringWidth(item.GetLabel())
		var itemWidth int
		if f.horizontal {
			fieldWidth := item.GetFieldWidth()
			if fieldWidth <= 0 {
				fieldWidth = DefaultFormFieldWidth
			}
			labelWidth++
			itemWidth = labelWidth + fieldWidth
		} else {
			// We want all fields to align vertically.
			labelWidth = maxLabelWidth
			itemWidth = width
		}
		itemHeight := item.GetFieldHeight()
		if itemHeight <= 0 {
			itemHeight = DefaultFormFieldHeight
		}

		// Advance to next line if there is no space.
		if f.horizontal && x+labelWidth+1 >= rightLimit {
			x = startX
			y += lineHeight + 1
			lineHeight = itemHeight
		}

		// Update line height.
		if itemHeight > lineHeight {
			lineHeight = itemHeight
		}

		// Adjust the item's attributes.
		if x+itemWidth >= rightLimit {
			itemWidth = rightLimit - x
		}

		// Save position.
		positions[index].x = x
		positions[index].y = y
		positions[index].width = itemWidth
		positions[index].height = itemHeight
		if item.HasFocus() {
			focusedPosition = positions[index]
		}

		// Advance to next item.
		if f.horizontal {
			x += itemWidth + f.itemPadding
		} else {
			y += itemHeight + f.itemPadding
		}
	}

	// Determine vertical offset based on the position of the focused item.
	var offset int
	if focusedPosition.y+focusedPosition.height > bottomLimit {
		offset = focusedPosition.y + focusedPosition.height - bottomLimit
		if focusedPosition.y-offset < topLimit {
			offset = focusedPosition.y - topLimit
		}
	}

	// Draw items.
	for index, item := range f.items {
		// Set position.
		y := positions[index].y - offset
		height := positions[index].height
		item.SetRect(positions[index].x, y, positions[index].width, height)

		// Is this item visible?
		if y+height <= topLimit || y >= bottomLimit {
			continue
		}

		// Draw items with focus last (in case of overlaps).
		if item.HasFocus() {
			defer item.Draw(screen)
		} else {
			item.Draw(screen)
		}
	}
}

// Focus is called by the application when the primitive receives focus.
func (f *RootContainer) Focus(delegate func(p tview.Primitive)) {
	// Hand on the focus to one of our child elements.
	if f.focusedElement < 0 || f.focusedElement >= len(f.items) {
		f.focusedElement = 0
	}
	var handler func(key tcell.Key)
	handler = func(key tcell.Key) {
		if key >= 0 {
			f.lastFinishedKey = key
		}
		switch key {
		case tcell.KeyTab, tcell.KeyEnter:
			f.focusedElement++
			f.Focus(delegate)
		case tcell.KeyBacktab:
			f.focusedElement--
			if f.focusedElement < 0 {
				f.focusedElement = len(f.items) - 1
			}
			f.Focus(delegate)
		case tcell.KeyEscape:
			if f.cancel != nil {
				f.cancel()
			} else {
				f.focusedElement = 0
				f.Focus(delegate)
			}
		default:
			if key < 0 && f.lastFinishedKey >= 0 {
				// Repeat the last action.
				handler(f.lastFinishedKey)
			}
		}
	}

	// Track whether a form item has focus.
	var itemFocused bool

	// Set the handler and focus for all items and buttons.
	for index, item := range f.items {
		item.SetFinishedFunc(handler)
		if f.focusedElement == index {
			itemFocused = true
			func(i tview.FormItem) { // Wrapping might not be necessary anymore in future Go versions.
				defer delegate(i)
			}(item)
		}
	}

	// If no item was focused, focus the form itself.
	if !itemFocused {
		f.Box.Focus(delegate)
	}
}

// HasFocus returns whether or not this primitive has focus.
func (f *RootContainer) HasFocus() bool {
	if f.focusIndex() >= 0 {
		return true
	}
	return f.Box.HasFocus()
}

// focusIndex returns the index of the currently focused item, counting form
// items first, then buttons. A negative value indicates that no containeed item
// has focus.
func (f *RootContainer) focusIndex() int {
	for index, item := range f.items {
		if item.HasFocus() {
			return index
		}
	}
	return -1
}

// MouseHandler returns the mouse handler for this primitive.
func (f *RootContainer) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return f.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		// At the end, update f.focusedElement and prepare current item/button.
		defer func() {
			if consumed {
				index := f.focusIndex()
				if index >= 0 {
					f.focusedElement = index
				}
			}
		}()

		// Determine items to pass mouse events to.
		for _, item := range f.items {
			// Exclude TextView items from mouse-down events as they are
			// read-only items and thus should not be focused.
			if _, ok := item.(*tview.TextView); ok && action == tview.MouseLeftDown {
				continue
			}

			consumed, capture = item.MouseHandler()(action, event, setFocus)
			if consumed {
				return
			}
		}

		// A mouse down anywhere else will return the focus to the last selected
		// element.
		if action == tview.MouseLeftDown && f.InRect(event.Position()) {
			f.Focus(setFocus)
			consumed = true
		}

		return
	})
}

// InputHandler returns the handler for this primitive.
func (f *RootContainer) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return f.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		for _, item := range f.items {
			if item != nil && item.HasFocus() {
				if handler := item.InputHandler(); handler != nil {
					handler(event, setFocus)
					return
				}
			}
		}
	})
}

func main() {
	app := tview.NewApplication()
	form := NewRootContainer().
		AddTextView("Notes", "This is just a demo.\nYou can enter whatever you wish.", 40, 2, true, true).
		SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
