package astilectron

import (
	"context"
	"sync"

	"github.com/asticode/go-astikit"
)

const (
	EventNameBrowserViewCmdSetBounds   = "browser.view.cmd.set.bounds"
	EventNameBrowserViewEventSetBounds = "browser.view.event.set.bounds"
)

type BrowserView struct {
	*object
	callbackIdentifier *identifier
	l                  astikit.SeverityLogger
	m                  sync.Mutex // Locks o
	o                  *WindowOptions
}

func newBrowserView(ctx context.Context, l astikit.SeverityLogger, o Options, p Paths, url string, wo *WindowOptions, d *dispatcher, i *identifier, wrt *writer) (*BrowserView, error) {
	b := &BrowserView{
		callbackIdentifier: newIdentifier(),
		l:                  l,
		o:                  wo,
		object:             newObject(ctx, d, i, wrt, i.new()),
	}

	return b, nil
}

func (b *BrowserView) SetBounds(bounds *RectangleOptions) {
	if err := b.ctx.Err(); err != nil {
		return
	}

	synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdSetBounds, TargetID: b.id, Bounds: bounds}, EventNameBrowserViewEventSetBounds)
}
