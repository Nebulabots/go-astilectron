package astilectron

import (
	"context"
	"fmt"
	stdUrl "net/url"
	"path/filepath"

	"github.com/asticode/go-astikit"
)

const (
	EventNameBrowserViewCmdCreate                            = "browser.view.cmd.create"
	EventNameBrowserViewCmdSetBounds                         = "browser.view.cmd.set.bounds"
	EventNameBrowserViewEventSetBounds                       = "browser.view.event.set.bounds"
	EventNameBrowserViewEventDidFinishLoad                   = "browser.view.event.did.finish.load"
	EventNameBrowserViewCmdLoadURL                           = "browser.view.cmd.load.url"
	EventNameBrowserViewEventLoadedURL                       = "browser.view.event.loaded.url"
	EventNameBrowserViewCmdWebContentsExecuteJavaScript      = "browser.view.cmd.web.contents.execute.javascript"
	EventNameBrowserViewEventWebContentsExecutedJavaScript   = "browser.view.event.web.contents.executed.javascript"
	EventNameBrowserViewCmdGetBounds                         = "browser.view.cmd.get.bounds"
	EventNameBrowserViewEventGetBounds                       = "browser.view.event.get.bounds"
	EventNameBrowserViewCmdInterceptStringProtocol           = "browser.view.cmd.intercept.string.protocol"
	EventNameBrowserViewEventInterceptStringProtocol         = "browser.view.event.intercept.string.protocol"
	EventNameBrowserViewEventInterceptStringProtocolCallback = "browser.view.event.intercept.string.protocol.callback"
	EventNameBrowserViewCmdSetBackgroundColor                = "browser.view.cmd.set.background.color"
	EventNameBrowserViewEventSetBackgroundColor              = "browser.view.event.set.background.color"
	EventNameBrowserViewCmdSetAutoResize                     = "browser.view.cmd.set.auto.resize"
	EventNameBrowserViewEventSetAutoResize                   = "browser.view.event.set.auto.resize"
	EventNameBrowserViewCmdSetProxy                          = "browser.view.cmd.web.contents.set.proxy"
	EventNameBrowserViewEventSetProxy                        = "browser.view.event.web.contents.set.proxy"
	EventNameBrowserViewCmdOpenDevTools                      = "browser.view.cmd.open.dev.tools"
	EventNameBrowserViewCmdCloseDevTools                     = "browser.view.cmd.close.dev.tools"
	EventNameBrowserViewCmdSetUserAgent                      = "browser.view.cmd.set.user.agent"
	EventNameBrowserViewEventSetUserAgent                    = "browser.view.event.set.user.agent"
	EventNameBrowserViewCmdUninterceptProtocol               = "browser.view.cmd.unintercept.string.protocol"
	EventNameBrowserViewEventUninterceptProtocol             = "browser.view.event.unintercept.string.protocol"
)

type BrowserView struct {
	*object
	callbackIdentifier *identifier
	l                  astikit.SeverityLogger
	o                  *WindowOptions
	url                *stdUrl.URL
	ID                 string
	Session            *Session
}

func newBrowserView(ctx context.Context, l astikit.SeverityLogger, o Options, p Paths, url string, wo *WindowOptions, s *Session, d *dispatcher, i *identifier, wrt *writer) (b *BrowserView, err error) {
	id := i.new()

	b = &BrowserView{
		callbackIdentifier: newIdentifier(),
		l:                  l,
		o:                  wo,
		object:             newObject(ctx, d, i, wrt, id),
		ID:                 id,
	}

	b.Session = s

	if url != "" {
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
	}

	return
}

func (b *BrowserView) SetBounds(bounds *RectangleOptions) (err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}

	_, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdSetBounds, TargetID: b.id, Bounds: bounds}, EventNameBrowserViewEventSetBounds)
	return
}

func (b *BrowserView) GetBounds(bounds *RectangleOptions) (e Event, err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}

	e, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdGetBounds, TargetID: b.id}, EventNameBrowserViewEventGetBounds)
	return
}

func (b *BrowserView) SetBackgroundColor(color string) (err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}

	_, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdSetBackgroundColor, TargetID: b.id, Color: color}, EventNameBrowserViewEventSetBackgroundColor)
	return
}

func (b *BrowserView) SetAutoResize(resizeOptions *ResizeOptions) (err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}

	_, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdSetAutoResize, TargetID: b.id, ResizeOptions: resizeOptions}, EventNameBrowserViewEventSetAutoResize)
	return
}

func (b *BrowserView) Create() (err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}

	if b.url != nil {
		_, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdCreate, TargetID: b.id, URL: b.url.String(), WindowOptions: b.o}, EventNameBrowserViewEventDidFinishLoad)
	} else {
		_, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdCreate, TargetID: b.id, WindowOptions: b.o}, EventNameBrowserViewEventDidFinishLoad)
	}
	return
}

func (b *BrowserView) LoadURL(url string, load *Load) (err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}

	_, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdLoadURL, TargetID: b.id, URL: url, Load: load}, EventNameBrowserViewEventLoadedURL)
	return
}

//todo: code can only return as string rn
func (b *BrowserView) ExecuteJavaScript(code string) (e Event, err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}
	e, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdWebContentsExecuteJavaScript, TargetID: b.id, Code: code}, EventNameBrowserViewEventWebContentsExecutedJavaScript)
	return
}

//func passed in from user
func (b *BrowserView) InterceptStringProtocol(scheme string, fn func(i Event) (mimeType string, data string, deleteListener bool)) (err error) {
	b.On(EventNameBrowserViewEventInterceptStringProtocol, func(i Event) (deleteListener bool) {
		mimeType, data, deleteListener := fn(i)

		if err = b.w.write(Event{CallbackID: i.CallbackID, Name: EventNameBrowserViewEventInterceptStringProtocolCallback, TargetID: b.id, Scheme: scheme, MimeType: mimeType, Data: data}); err != nil {
			return
		}

		return
	})

	if err = b.ctx.Err(); err != nil {
		return
	}

	b.w.write(Event{Name: EventNameBrowserViewCmdInterceptStringProtocol, TargetID: b.id, Scheme: scheme})

	return
}

func (b *BrowserView) SetProxy(proxy *WindowProxyOptions) (err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}

	_, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdSetProxy, TargetID: b.id, Proxy: proxy}, EventNameBrowserViewEventSetProxy)
	return
}

func (b *BrowserView) OpenDevTools() (err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}

	err = b.w.write(Event{Name: EventNameBrowserViewCmdOpenDevTools, TargetID: b.id})
	return
}

func (b *BrowserView) CloseDevTools() (err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}

	err = b.w.write(Event{Name: EventNameBrowserViewCmdCloseDevTools, TargetID: b.id})
	return
}

func (b *BrowserView) UninterceptProtocol(scheme string) (err error) {
	if err = b.ctx.Err(); err != nil {
		return
	}

	_, err = synchronousEvent(b.ctx, b, b.w, Event{Name: EventNameBrowserViewCmdUninterceptProtocol, TargetID: b.id, Scheme: scheme}, EventNameBrowserViewEventUninterceptProtocol)
	return
}

func (b *BrowserView) OnLogin(fn func(i Event) (username, password string, err error)) {
	b.On(EventNameWebContentsEventLogin, func(i Event) (deleteListener bool) {
		// Get username and password
		username, password, err := fn(i)
		if err != nil {
			b.l.Error(fmt.Errorf("getting username and password failed: %w", err))
			return
		}

		// No auth
		if len(username) == 0 && len(password) == 0 {
			return
		}

		// Send message back
		if err = b.w.write(Event{CallbackID: i.CallbackID, Name: EventNameWebContentsEventLoginCallback, Password: password, TargetID: b.id, Username: username}); err != nil {
			b.l.Error(fmt.Errorf("writing login callback message failed: %w", err))
			return
		}
		return
	})

	if err := b.ctx.Err(); err != nil {
		return
	}

	b.w.write(Event{Name: EventNameWebContentsEventLogin, TargetID: b.id})

	return
}
