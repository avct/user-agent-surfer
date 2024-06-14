// Package uasurfer provides fast and reliable abstraction
// of HTTP User-Agent strings. The philosophy is to identify
// technologies that holds >1% market share, and to avoid
// expending resources and accuracy on guessing at esoteric UA
// strings.
package uasurfer

import (
	"strings"
	"sync"
	"unsafe"
)

//go:generate stringer -type=DeviceType,BrowserName,OSName,Platform -output=const_string.go

// DeviceType (int) returns a constant.
type DeviceType int

// A complete list of supported devices in the
// form of constants.
const (
	DeviceUnknown DeviceType = iota
	DeviceComputer
	DeviceTablet
	DevicePhone
	DeviceConsole
	DeviceWearable
	DeviceTV
)

// StringTrimPrefix is like String() but trims the "Device" prefix
func (d DeviceType) StringTrimPrefix() string {
	return strings.TrimPrefix(d.String(), "Device")
}

// BrowserName (int) returns a constant.
type BrowserName int

// A complete list of supported web browsers in the
// form of constants.
const (
	BrowserUnknown BrowserName = iota
	BrowserChrome
	BrowserIE
	BrowserSafari
	BrowserFirefox
	BrowserAndroid
	BrowserOpera
	BrowserBlackberry
	BrowserUCBrowser
	BrowserSilk
	BrowserNokia
	BrowserNetFront
	BrowserQQ
	BrowserMaxthon
	BrowserSogouExplorer
	BrowserSpotify
	BrowserNintendo
	BrowserSamsung
	BrowserYandex
	BrowserCocCoc
	BrowserBot // Bot list begins here
	BrowserAppleBot
	BrowserBaiduBot
	BrowserBingBot
	BrowserDuckDuckGoBot
	BrowserFacebookBot
	BrowserGoogleBot
	BrowserLinkedInBot
	BrowserMsnBot
	BrowserPingdomBot
	BrowserTwitterBot
	BrowserYandexBot
	BrowserCocCocBot
	BrowserYahooBot // Bot list ends here
)

// StringTrimPrefix is like String() but trims the "Browser" prefix
func (b BrowserName) StringTrimPrefix() string {
	return strings.TrimPrefix(b.String(), "Browser")
}

// OSName (int) returns a constant.
type OSName int

// A complete list of supported OSes in the
// form of constants. For handling particular versions
// of operating systems (e.g. Windows 2000), see
// the README.md file.
const (
	OSUnknown OSName = iota
	OSWindowsPhone
	OSWindows
	OSMacOSX
	OSiOS
	OSAndroid
	OSBlackberry
	OSChromeOS
	OSKindle
	OSWebOS
	OSLinux
	OSPlaystation
	OSXbox
	OSNintendo
	OSBot
)

// StringTrimPrefix is like String() but trims the "OS" prefix
func (o OSName) StringTrimPrefix() string {
	return strings.TrimPrefix(o.String(), "OS")
}

// Platform (int) returns a constant.
type Platform int

// A complete list of supported platforms in the
// form of constants. Many OSes report their
// true platform, such as Android OS being Linux
// platform.
const (
	PlatformUnknown Platform = iota
	PlatformWindows
	PlatformMac
	PlatformLinux
	PlatformiPad
	PlatformiPhone
	PlatformiPod
	PlatformBlackberry
	PlatformWindowsPhone
	PlatformPlaystation
	PlatformXbox
	PlatformNintendo
	PlatformBot
)

// StringTrimPrefix is like String() but trims the "Platform" prefix
func (p Platform) StringTrimPrefix() string {
	return strings.TrimPrefix(p.String(), "Platform")
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) Less(c Version) bool {
	if v.Major < c.Major {
		return true
	}

	if v.Major > c.Major {
		return false
	}

	if v.Minor < c.Minor {
		return true
	}

	if v.Minor > c.Minor {
		return false
	}

	return v.Patch < c.Patch
}

type UserAgent struct {
	Browser    Browser
	OS         OS
	DeviceType DeviceType
}

type Browser struct {
	Name    BrowserName
	Version Version
}

type OS struct {
	Platform Platform
	Name     OSName
	Version  Version
}

// Reset resets the UserAgent to it's zero value
func (ua *UserAgent) Reset() {
	ua.Browser = Browser{}
	ua.OS = OS{}
	ua.DeviceType = DeviceUnknown
}

// IsBot returns true if the UserAgent represent a bot
func (ua *UserAgent) IsBot() bool {
	if ua.Browser.Name >= BrowserBot && ua.Browser.Name <= BrowserYahooBot {
		return true
	}
	if ua.OS.Name == OSBot {
		return true
	}
	if ua.OS.Platform == PlatformBot {
		return true
	}
	return false
}

// Parse accepts a raw user agent (string) and returns the UserAgent.
func Parse(ua string) *UserAgent {
	dest := new(UserAgent)
	parse(ua, dest)
	return dest
}

// ParseUserAgent is the same as Parse, but populates the supplied UserAgent.
// It is the caller's responsibility to call Reset() on the UserAgent before
// passing it to this function.
func ParseUserAgent(ua string, dest *UserAgent) {
	parse(ua, dest)
}

func parse(ua string, dest *UserAgent) {
	bp := bytesPool.Get().(*[]byte)
	b := *bp

	b = append(b[:0], ua...)
	lowercaseBytes(b)
	ua = b2s(b)

	switch {
	case len(ua) == 0:
		dest.OS.Platform = PlatformUnknown
		dest.OS.Name = OSUnknown
		dest.Browser.Name = BrowserUnknown
		dest.DeviceType = DeviceUnknown

	// stop on on first case returning true
	case dest.evalOS(ua):
	case dest.evalBrowserName(ua):
	default:
		dest.evalBrowserVersion(ua)
		dest.evalDevice(ua)
	}

	*bp = b
	bytesPool.Put(bp)
}

// b2s converts a byte slice to a string without allocating.
// WARNING: changing the byte slice will change the string as well!
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

const toLower = 'a' - 'A'

var (
	bytesPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 0, 1024)
			return &b
		},
	}

	toLowerTable = func() [256]byte {
		var a [256]byte
		for i := 0; i < 256; i++ {
			c := byte(i)
			if c >= 'A' && c <= 'Z' {
				c += toLower
			}
			a[i] = c
		}
		return a
	}()
)

// Lowercase all ascii characters in b.
func lowercaseBytes(b []byte) {
	for i := 0; i < len(b); i++ {
		p := &b[i]
		*p = toLowerTable[*p]
	}
}
