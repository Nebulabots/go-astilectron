package astilectron

import (
	"context"
	"fmt"
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
)

// Session represents a session
// TODO Add missing session methods
// TODO Add missing session events
// https://github.com/electron/electron/blob/v1.8.1/docs/api/session.md
type Session struct {
	*object
}

// newSession creates a new session
func newSession(ctx context.Context, d *dispatcher, i *identifier, w *writer) *Session {
	return &Session{object: newObject(ctx, d, i, w, i.new())}
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

func (w *Window) OnBeforeRequest(fn func(i Event) (string, string, bool)) (err error) {
	// Setup the event to handle the callback
	w.On(EventNameWebContentsEventSessionWebRequestOnBeforeRequestCallback, func(i Event) (deleteListener bool) {
		// Get mime type, data and whether the listener should be deleted.
		mimeType, data, deleteListener := fn(i)

		// Send message back
		if err = w.w.write(Event{CallbackID: i.CallbackID, Name: EventNameWebContentsEventInterceptStringProtocolCallback, TargetID: w.id, MimeType: mimeType, Data: data}); err != nil {
			w.l.Error(fmt.Errorf("writing on before request callback message failed: %w", err))
			return
		}

		return
	})

	if err = w.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(w.ctx, w, w.w, Event{Name: EventNameWebContentsEventSessionWebRequestOnBeforeRequest, TargetID: w.id}, EventNameWebContentsEventSessionWebRequestOnBeforeRequestCallback)
	return
}
