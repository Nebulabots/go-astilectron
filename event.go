package astilectron

import (
	"encoding/json"
	"errors"
)

// Target IDs
const (
	targetIDApp  = "app"
	targetIDDock = "dock"
)

// Event represents an event
type Event struct {
	// This is the base of the event
	Name     string `json:"name"`
	TargetID string `json:"targetID,omitempty"`

	// This is a list of all possible payloads.
	// A choice was made not to use interfaces since it's a pain in the ass asserting each an every payload afterwards
	// We use pointers so that omitempty works
	AcceptLanguages string            `json:"acceptLanguages,omitempty"`
	AuthInfo        *EventAuthInfo    `json:"authInfo,omitempty"`
	Badge           string            `json:"badge,omitempty"`
	BounceType      string            `json:"bounceType,omitempty"`
	Bounds          *RectangleOptions `json:"bounds,omitempty"`
	BrowserViewID   string            `json:"browserViewID,omitempty"`
	Cancel          *bool             `json:"cancel,omitempty"`
	CallbackID      string            `json:"callbackId,omitempty"`
	Color           string            `json:"color,omitempty"`
	Code            string            `json:"code,omitempty"`
	//todo: can only be a string now?
	CodeResult string          `json:"codeResult,omitempty"`
	Cookies    []SessionCookie `json:"cookies,omitempty"`
	// https://www.electronjs.org/docs/api/structures/protocol-response
	Data                  string                 `json:"data,omitempty"`
	Displays              *EventDisplays         `json:"displays,omitempty"`
	DialogOptions         *DialogOptions         `json:"dialogOptions,omitempty"`
	Error                 string                 `json:"error,omitempty"`
	FilePath              string                 `json:"filePath,omitempty"`
	ID                    *int                   `json:"id,omitempty"`
	Filter                *FilterOptions         `json:"filter,omitempty"`
	Image                 string                 `json:"image,omitempty"`
	Index                 *int                   `json:"index,omitempty"`
	Menu                  *EventMenu             `json:"menu,omitempty"`
	Load                  *Load                  `json:"load,omitempty"`
	MenuItem              *EventMenuItem         `json:"menuItem,omitempty"`
	MenuItemOptions       *MenuItemOptions       `json:"menuItemOptions,omitempty"`
	MenuItemPosition      *int                   `json:"menuItemPosition,omitempty"`
	MenuPopupOptions      *MenuPopupOptions      `json:"menuPopupOptions,omitempty"`
	Message               *EventMessage          `json:"message,omitempty"`
	MimeType              string                 `json:"mimeType,omitempty"`
	NotificationOptions   *NotificationOptions   `json:"notificationOptions,omitempty"`
	Password              string                 `json:"password,omitempty"`
	Path                  string                 `json:"path,omitempty"`
	Paths                 []string               `json:"paths,omitempty"`
	Partition             string                 `json:"partition,omitempty"`
	Proxy                 *WindowProxyOptions    `json:"proxy,omitempty"`
	RedirectURL           string                 `json:"redirectURL,omitempty"`
	Reply                 string                 `json:"reply,omitempty"`
	ResizeOptions         *ResizeOptions         `json:"resizeOptions,omitempty"`
	Request               *EventRequest          `json:"request,omitempty"`
	Scheme                string                 `json:"scheme,omitempty"`
	SecondInstance        *EventSecondInstance   `json:"secondInstance,omitempty"`
	SessionID             string                 `json:"sessionId,omitempty"`
	ShowOpenDialogOptions *ShowOpenDialogOptions `json:"showOpenDialogOptions,omitempty"`
	Supported             *Supported             `json:"supported,omitempty"`
	TrayOptions           *TrayOptions           `json:"trayOptions,omitempty"`
	URL                   string                 `json:"url,omitempty"`
	URLNew                string                 `json:"newUrl,omitempty"`
	URLOld                string                 `json:"oldUrl,omitempty"`
	UserAgent             string                 `json:"userAgent,omitempty"`
	Username              string                 `json:"username,omitempty"`
	WindowID              string                 `json:"windowId,omitempty"`
	WindowOptions         *WindowOptions         `json:"windowOptions,omitempty"`
}

// EventAuthInfo represents an event auth info
type EventAuthInfo struct {
	Host    string `json:"host,omitempty"`
	IsProxy *bool  `json:"isProxy,omitempty"`
	Port    *int   `json:"port,omitempty"`
	Realm   string `json:"realm,omitempty"`
	Scheme  string `json:"scheme,omitempty"`
}

// EventDisplays represents events displays
type EventDisplays struct {
	All     []*DisplayOptions `json:"all,omitempty"`
	Primary *DisplayOptions   `json:"primary,omitempty"`
}

type FilterOptions struct {
	Urls []string `json:"urls,omitempty"`
}

// EventMessage represents an event message
type EventMessage struct {
	i interface{}
}

// newEventMessage creates a new event message
func newEventMessage(i interface{}) *EventMessage {
	return &EventMessage{i: i}
}

// MarshalJSON implements the JSONMarshaler interface
func (p *EventMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.i)
}

// Unmarshal unmarshals the payload into the given interface
func (p *EventMessage) Unmarshal(i interface{}) error {
	if b, ok := p.i.([]byte); ok {
		return json.Unmarshal(b, i)
	}
	return errors.New("event message should []byte")
}

// UnmarshalJSON implements the JSONUnmarshaler interface
func (p *EventMessage) UnmarshalJSON(i []byte) error {
	p.i = i
	return nil
}

type ResizeOptions struct {
	Width      bool `json:"width,omitempty"`
	Height     bool `json:"height,omitempty"`
	Horizontal bool `json:"horizontal,omitempty"`
	Vertical   bool `json:"vertical,omitempty"`
}

type Load struct {
	HttpReferrer      string `json:"httpReferrer,omitempty"`
	UserAgent         string `json:"userAgent,omitempty"`
	ExtraHeaders      string `json:"extraHeaders,omitempty"`
	BaseURLForDataURL string `json:"baseURLForDataURL,omitempty"`
}

// EventMenu represents an event menu
type EventMenu struct {
	*EventSubMenu
}

// EventMenuItem represents an event menu item
type EventMenuItem struct {
	ID      string           `json:"id"`
	Options *MenuItemOptions `json:"options,omitempty"`
	RootID  string           `json:"rootId"`
	SubMenu *EventSubMenu    `json:"submenu,omitempty"`
}

type UploadData struct {
	Type  string          `json:"type,omitempty"`
	Bytes map[string]byte `json:"bytes,omitempty"`
}

// EventRequest represents an event request
type EventRequest struct {
	Method     string             `json:"method,omitempty"`
	Referrer   string             `json:"referrer,omitempty"`
	URL        string             `json:"url,omitempty"`
	UploadData map[int]UploadData `json:"uploadData,omitempty"`
}

// EventSecondInstance represents data related to a second instance of the app being started
type EventSecondInstance struct {
	CommandLine      []string `json:"commandLine,omitempty"`
	WorkingDirectory string   `json:"workingDirectory,omitempty"`
}

// EventSubMenu represents a sub menu event
type EventSubMenu struct {
	ID     string           `json:"id"`
	Items  []*EventMenuItem `json:"items,omitempty"`
	RootID string           `json:"rootId"`
}
