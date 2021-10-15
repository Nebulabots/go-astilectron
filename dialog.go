package astilectron

import "context"

// Message box types
const (
	EventNameDialogEventDestroyed = "dialog.event.destroyed"
	EventNameDialogEventCreated   = "dialog.event.created"
	EventNameDialogCmdCreate      = "dialog.cmd.create"
	EventNameDialogCmdDestroy     = "dialog.cmd.destroy"

	EventNameDialogCmdShowOpenDialog   = "dialog.cmd.show.open.dialog"
	EventNameDialogEventShowOpenDialog = "dialog.event.show.open.dialog"

	MessageBoxTypeError    = "error"
	MessageBoxTypeInfo     = "info"
	MessageBoxTypeNone     = "none"
	MessageBoxTypeQuestion = "question"
	MessageBoxTypeWarning  = "warning"
)

// MessageBoxOptions represents message box options
// We must use pointers since GO doesn't handle optional fields whereas NodeJS does. Use astikit.BoolPtr, astikit.IntPtr or astikit.StrPtr
// to fill the struct
// https://github.com/electron/electron/blob/v1.8.1/docs/api/dialog.md#dialogshowmessageboxbrowserwindow-options-callback
type MessageBoxOptions struct {
	Buttons         []string `json:"buttons,omitempty"`
	CancelID        *int     `json:"cancelId,omitempty"`
	CheckboxChecked *bool    `json:"checkboxChecked,omitempty"`
	CheckboxLabel   string   `json:"checkboxLabel,omitempty"`
	ConfirmID       *int     `json:"confirmId,omitempty"`
	DefaultID       *int     `json:"defaultId,omitempty"`
	Detail          string   `json:"detail,omitempty"`
	Icon            string   `json:"icon,omitempty"`
	Message         string   `json:"message,omitempty"`
	NoLink          *bool    `json:"noLink,omitempty"`
	Title           string   `json:"title,omitempty"`
	Type            string   `json:"type,omitempty"`
}

type DialogOptions struct{}

type FileFilter struct {
	Name       string   `json:"name,omitempty"`
	Extensions []string `json:"extensions,omitempty"`
}

type FileFilterOptions []FileFilter

type ShowOpenDialogOptions struct {
	Title       string            `json:"title,omitempty"`
	DefaultPath string            `json:"defaultPath,omitempty"`
	ButtonLabel string            `json:"buttonLabel,omitempty"`
	Filters     FileFilterOptions `json:"filters,omitempty"`
	Properties  []string          `json:"properties,omitempty"`
}

type Dialog struct {
	*object
	o *DialogOptions
}

// newDialog creates a new dialog
func newDialog(ctx context.Context, o *DialogOptions, d *dispatcher, i *identifier, wrt *writer) (t *Dialog) {
	// Init
	t = &Dialog{
		o:      o,
		object: newObject(ctx, d, i, wrt, i.new()),
	}

	// Make sure the dialogs's context is cancelled once the destroyed event is received
	t.On(EventNameDialogEventDestroyed, func(e Event) (deleteListener bool) {
		t.cancel()
		return true
	})
	return
}

// Create creates the dialog
func (t *Dialog) Create() (err error) {
	if err = t.ctx.Err(); err != nil {
		return
	}
	var e = Event{Name: EventNameDialogCmdCreate, TargetID: t.id, DialogOptions: t.o}
	_, err = synchronousEvent(t.ctx, t, t.w, e, EventNameDialogEventCreated)
	return
}

// Destroy destroys the dialog
func (t *Dialog) Destroy() (err error) {
	if err = t.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(t.ctx, t, t.w, Event{Name: EventNameDialogCmdDestroy, TargetID: t.id}, EventNameDialogEventDestroyed)
	return
}

func (t *Dialog) ShowOpenDialog(o *ShowOpenDialogOptions) (e Event, err error) {
	if err = t.ctx.Err(); err != nil {
		return
	}
	e, err = synchronousEvent(t.ctx, t, t.w, Event{Name: EventNameDialogCmdShowOpenDialog, TargetID: t.id, ShowOpenDialogOptions: o}, EventNameDialogEventShowOpenDialog)
	return
}
