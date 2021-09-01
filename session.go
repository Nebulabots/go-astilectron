package astilectron

import (
	"context"
)

// Session event names
const (
	EventNameSessionCmdClearCache        = "session.cmd.clear.cache"
	EventNameSessionEventClearedCache    = "session.event.cleared.cache"
	EventNameSessionCmdFlushStorage      = "session.cmd.flush.storage"
	EventNameSessionEventFlushedStorage  = "session.event.flushed.storage"
	EventNameSessionCmdLoadExtension     = "session.cmd.load.extension"
	EventNameSessionEventLoadedExtension = "session.event.loaded.extension"
	EventNameSessionEventWillDownload    = "session.event.will.download"
	EventNameSessionCmdSetCookies        = "session.cmd.cookies.set"
	EventNameSessionEventSetCookies      = "session.event.cookies.set"
	EventNameSessionCmdGetCookies        = "session.cmd.cookies.get"
	EventNameSessionEventGetCookies      = "session.event.cookies.get"
	EventNameSessionCmdFromPartition     = "session.cmd.from.partition"
	EventNameSessionEventFromPartition   = "session.event.from.partition"
)

// Session represents a session
// TODO Add missing session methods
// TODO Add missing session events
// https://github.com/electron/electron/blob/v1.8.1/docs/api/session.md
type Session struct {
	*object
	Protocol *Protocol
}

// newSession creates a new session
func NewSession(ctx context.Context, d *dispatcher, i *identifier, w *writer) *Session {
	s := &Session{object: newObject(ctx, d, i, w, i.new())}

	s.Protocol = newProtocol(ctx, d, i, w)

	return s
}

// ClearCache clears the Session's HTTP cache
func (s *Session) ClearCache() (err error) {
	if err = s.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(s.ctx, s, s.w, Event{Name: EventNameSessionCmdClearCache, TargetID: s.id}, EventNameSessionEventClearedCache)
	return
}

// FlushStorage writes any unwritten DOMStorage data to disk
func (s *Session) FlushStorage() (err error) {
	if err = s.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(s.ctx, s, s.w, Event{Name: EventNameSessionCmdFlushStorage, TargetID: s.id}, EventNameSessionEventFlushedStorage)
	return
}

// Loads a chrome extension
func (s *Session) LoadExtension(path string) (err error) {
	if err = s.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(s.ctx, s, s.w, Event{Name: EventNameSessionCmdLoadExtension, Path: path, TargetID: s.id}, EventNameSessionEventLoadedExtension)
	return
}

func (s *Session) OnBeforeRequest(fn func(i Event) (bool, string, bool)) (err error) {
	// Setup the event to handle the callback
	s.On(EventNameWebContentsEventSessionWebRequestOnBeforeRequest, func(i Event) (deleteListener bool) {
		cancel, redirectUrl, deleteListener := fn(i)

		// Send message back
		if err = s.w.write(Event{CallbackID: i.CallbackID, Name: EventNameWebContentsEventSessionWebRequestOnBeforeRequestCallback, TargetID: s.id, Cancel: &cancel, RedirectURL: redirectUrl}); err != nil {
			return
		}

		return
	})

	if err = s.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(s.ctx, s, s.w, Event{Name: EventNameWebContentsEventSessionWebRequestOnBeforeRequest, TargetID: s.id}, EventNameWindowEventWebContentsOnBeforeRequest)
	return
}

type SessionCookie struct {
	Url            string   `json:"url"`
	Name           string   `json:"name,omitempty"`
	Value          string   `json:"value,omitempty"`
	Domain         string   `json:"domain,omitempty"`
	Path           string   `json:"path,omitempty"`
	Secure         *bool    `json:"secure,omitempty"`
	HttpOnly       *bool    `json:"httpOnly,omitempty"`
	Session        *bool    `json:"session,omitempty"`
	ExpirationDate *float64 `json:"expirationDate,omitempty"`
	SameSite       string   `json:"sameSite,omitempty"`
}

func (s *Session) SetCookies(cookies []SessionCookie) (err error) {
	if err = s.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(s.ctx, s, s.w, Event{Name: EventNameSessionCmdSetCookies, TargetID: s.id, Cookies: cookies}, EventNameSessionEventSetCookies)
	return
}

func (s *Session) GetCookies() (e Event, err error) {
	if err = s.ctx.Err(); err != nil {
		return
	}
	e, err = synchronousEvent(s.ctx, s, s.w, Event{Name: EventNameSessionCmdGetCookies, TargetID: s.id}, EventNameSessionEventGetCookies)
	return
}

func (s *Session) FromPartition(partition string) (e Event, err error) {
	if err = s.ctx.Err(); err != nil {
		return
	}

	e, err = synchronousEvent(s.ctx, s, s.w, Event{Name: EventNameSessionCmdFromPartition, SessionID: s.id, Partition: partition}, EventNameSessionEventFromPartition)
	return
}
