package astilectron

import (
	"context"
	"fmt"
	stdUrl "net/url"
	"path/filepath"
	"sync"

	"github.com/asticode/go-astikit"
)

const (
	EventNameBrowserViewCmdCreate          = "browser.view.cmd.create"
	EventNameBrowserViewCmdSetBounds       = "browser.view.cmd.set.bounds"
	EventNameBrowserViewEventSetBounds     = "browser.view.event.set.bounds"
	EventNameBrowserViewEventDidFinishLoad = "browser.view.event.did.finish.load"
)

type BrowserView struct {
	*object
	callbackIdentifier *identifier
	l                  astikit.SeverityLogger
	m                  sync.Mutex // Locks o
	o                  *WindowOptions
	url                *stdUrl.URL
}

func newBrowserView(ctx context.Context, l astikit.SeverityLogger, o Options, p Paths, url string, wo *WindowOptions, d *dispatcher, i *identifier, wrt *writer) (*BrowserView, error) {
	b := &BrowserView{
		callbackIdentifier: newIdentifier(),
		l:                  l,
		o:                  wo,
		object:             newObject(ctx, d, i, wrt, i.new()),
	}

	var err error

	// Basic parse
	if b.url, err = stdUrl.Parse(url); err != nil {
		err = fmt.Errorf("std parsing of url %s failed: %w", url, err)
		return nil, err
	}

	// File
	if b.url.Scheme == "" {
		// Get absolute path
		if url, err = filepath.Abs(url); err != nil {
			err = fmt.Errorf("getting absolute path of %s failed: %w", url, err)
			return nil, err
		}

		// Set url
		b.url = &stdUrl.URL{Path: filepath.ToSlash(url), Scheme: "file"}
	}

	return b, nil
}

func (b *BrowserView) SetBounds(bounds *RectangleOptions) {
	if err := b.ctx.Err(); err != nil {
		return
	}

	synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdSetBounds, TargetID: b.id, Bounds: bounds}, EventNameBrowserViewEventSetBounds)
}

func (b *BrowserView) Create() (err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdCreate, TargetID: b.id, URL: b.url.String(), WindowOptions: b.o}, EventNameBrowserViewEventDidFinishLoad)
	return
}
