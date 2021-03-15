package plex

import (
	"encoding/xml"
	"time"

	"github.com/DirtyCajunRice/go-utility/types"
	_ "github.com/DirtyCajunRice/go-utility/types"
)

// Pin Endpoint
type PinData struct {
	Errors           []Error   `json:"errors"`
	ID               int       `json:"id"`
	Code             string    `json:"code"`
	Product          string    `json:"product"`
	Trusted          bool      `json:"trusted"`
	ClientIdentifier string    `json:"clientIdentifier"`
	Location         Location  `json:"location"`
	ExpiresIn        int       `json:"expiresIn"`
	CreatedAt        time.Time `json:"createdAt"`
	ExpiresAt        time.Time `json:"expiresAt"`
	AuthToken        string    `json:"authToken"`
	NewRegistration  bool      `json:"newRegistration"`
}
type Location struct {
	Code         string `json:"code"`
	Country      string `json:"country"`
	City         string `json:"city"`
	TimeZone     string `json:"time_zone"`
	PostalCode   string `json:"postal_code"`
	Subdivisions string `json:"subdivisions"`
	Coordinates  string `json:"coordinates"`
}
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e Error) Error() string {
	return error(e).Error()
}

// Servers endpoint
type MediaContainerXML struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Servers []Server `xml:"Server"`
	Size    string   `xml:"size,attr"`
}

type Server struct {
	XMLName           xml.Name            `xml:"Server" json:"-" db:"-"`
	AccessToken       string              `xml:"accessToken,attr" db:"-"`
	Address           string              `xml:"address,attr" db:"address"`
	CreatedAt         types.UnixTimestamp `xml:"createdAt,attr" json:"-" db:"-"`
	Host              string              `xml:"host,attr" db:"host"`
	LocalAddresses    string              `xml:"localAddresses,attr" db:"local_addresses"`
	MachineIdentifier string              `xml:"machineIdentifier,attr" db:"machine_identifier" goqu:"skipupdate"`
	Name              string              `xml:"name,attr" db:"name"`
	Owned             bool                `xml:"owned,attr" db:"-"`
	Port              int                 `xml:"port,attr" db:"port"`
	Scheme            string              `xml:"scheme,attr" db:"scheme"`
	Synced            bool                `xml:"synced,attr" db:"-"`
	UpdatedAt         types.UnixTimestamp `xml:"updatedAt,attr" json:"-" db:"-"`
	Version           string              `xml:"version,attr" db:"version"`
	OwnerId           int                 `xml:"ownerId,attr" db:"owner_id"`
}

type Data struct {
	Size            int    `json:"size"`
	AllowSync       bool   `json:"allowSync"`
	Identifier      string `json:"identifier"`
	MediaTagPrefix  string `json:"mediaTagPrefix"`
	MediaTagVersion int    `json:"mediaTagVersion"`
	Title1          string `json:"title1"`
}

// Library Sections Endpoint
type LibraryResponse struct {
	Data LibraryData `json:"MediaContainer"`
}

type LibraryData struct {
	Data
	Sections []Library `json:"Directory"`
}

type Library struct {
	AllowSync        bool              `json:"-" db:"-"`
	Art              string            `json:"-" db:"-"`
	Composite        string            `json:"-" db:"-"`
	Filters          bool              `json:"-" db:"-"`
	Refreshing       bool              `json:"-" db:"-"`
	Thumb            string            `json:"-" db:"-"`
	Key              int               `json:"key,string" db:"key"`
	Type             string            `json:"type" db:"type"`
	Title            string            `json:"title" db:"title"`
	Agent            string            `json:"agent" db:"agent"`
	Scanner          string            `json:"scanner" db:"scanner"`
	Language         string            `json:"-" db:"-"`
	UUID             string            `json:"uuid" db:"uuid" goqu:"skipupdate"`
	UpdatedAt        int               `json:"-" db:"-"`
	CreatedAt        int               `json:"-" db:"-"`
	ScannedAt        int               `json:"-" db:"-"`
	Content          bool              `json:"-" db:"-"`
	Directory        bool              `json:"-" db:"-"`
	ContentChangedAt int               `json:"-" db:"-"`
	Hidden           int               `json:"-" db:"-"`
	Location         []LibraryLocation `json:"-" db:"-"`
}

type LibraryLocation struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
}

// Series endpoint
type SeriesResponse struct {
	Data SeriesData `json:"MediaContainer"`
}

type SeriesData struct {
	Data
	Art                 string   `json:"art"`
	LibrarySectionID    int      `json:"librarySectionID"`
	LibrarySectionTitle string   `json:"librarySectionTitle"`
	LibrarySectionUUID  string   `json:"librarySectionUUID"`
	Nocache             bool     `json:"nocache"`
	Thumb               string   `json:"thumb"`
	Title2              string   `json:"title2"`
	ViewGroup           string   `json:"viewGroup"`
	ViewMode            int      `json:"viewMode"`
	Series              []Series `json:"Metadata"`
}

type Series struct {
	RatingKey             int        `json:"ratingKey,string" db:"rating_key" goqu:"skipupdate"`
	Key                   string     `json:"key" db:"-"`
	SkipChildren          bool       `json:"skipChildren,omitempty" db:"-"`
	GUID                  string     `json:"guid" db:"-"`
	Studio                string     `json:"studio" db:"studio"`
	Type                  string     `json:"type" db:"-"`
	Title                 string     `json:"title" db:"title"`
	ContentRating         string     `json:"contentRating,omitempty" db:"content_rating"`
	Summary               string     `json:"summary" db:"-"`
	Index                 int        `json:"index" db:"-"`
	Rating                float64    `json:"rating,omitempty" db:"-"`
	ViewCount             int        `json:"viewCount,omitempty" db:"-"`
	LastViewedAt          int        `json:"lastViewedAt,omitempty" db:"-"`
	Year                  int        `json:"year" db:"year"`
	Thumb                 string     `json:"thumb" db:"-"`
	Art                   string     `json:"art" db:"-"`
	Banner                string     `json:"banner" db:"-"`
	Theme                 string     `json:"theme,omitempty" db:"-"`
	Duration              int        `json:"duration" db:"-"`
	OriginallyAvailableAt string     `json:"originallyAvailableAt" db:"-"`
	LeafCount             int        `json:"leafCount" db:"-"`
	ViewedLeafCount       int        `json:"viewedLeafCount" db:"-"`
	ChildCount            int        `json:"childCount" db:"-"`
	AddedAt               int        `json:"addedAt" db:"-"`
	UpdatedAt             int        `json:"updatedAt" db:"-"`
	Genre                 []Metadata `json:"Genre" db:"-"`
	Role                  []Metadata `json:"Role,omitempty" db:"-"`
	TitleSort             string     `json:"titleSort,omitempty" db:"-"`
	Collection            []Metadata `json:"Collection,omitempty" db:"-"`
}
type Metadata struct {
	Tag string `json:"tag"`
}

// User endpoint
type User struct {
	ID                      int          `json:"id" db:"id" goqu:"skipupdate"`
	UUID                    string       `json:"uuid" db:"uuid"`
	Username                string       `json:"username" db:"username"`
	Title                   string       `json:"title" db:"-"`
	Email                   string       `json:"email" db:"email"`
	Locale                  string       `json:"locale" db:"-"`
	Confirmed               bool         `json:"confirmed" db:"-"`
	EmailOnlyAuth           bool         `json:"emailOnlyAuth" db:"-"`
	HasPassword             bool         `json:"hasPassword" db:"-"`
	Protected               bool         `json:"protected" db:"-"`
	Thumb                   string       `json:"thumb" db:"thumb"`
	AuthToken               string       `json:"authToken" db:"-"`
	MailingListStatus       string       `json:"mailingListStatus" db:"-"`
	MailingListActive       bool         `json:"mailingListActive" db:"-"`
	ScrobbleTypes           string       `json:"scrobbleTypes" db:"-"`
	Country                 string       `json:"country" db:"-"`
	Pin                     string       `json:"pin" db:"-"`
	Subscription            Subscription `json:"subscription" db:"-"`
	SubscriptionDescription string       `json:"subscriptionDescription" db:"-"`
	Restricted              bool         `json:"restricted" db:"-"`
	Anonymous               interface{}  `json:"anonymous" db:"-"`
	Home                    bool         `json:"home" db:"-"`
	Guest                   bool         `json:"guest" db:"-"`
	HomeSize                int          `json:"homeSize" db:"-"`
	HomeAdmin               bool         `json:"homeAdmin" db:"-"`
	MaxHomeSize             int          `json:"maxHomeSize" db:"-"`
	CertificateVersion      int          `json:"certificateVersion" db:"-"`
	RememberExpiresAt       int          `json:"rememberExpiresAt" db:"-"`
	Profile                 Profile      `json:"profile" db:"-"`
	Entitlements            []string     `json:"entitlements" db:"-"`
	Roles                   []string     `json:"roles" db:"-"`
	Services                []Services   `json:"services" db:"-"`
	AdsConsent              interface{}  `json:"adsConsent" db:"-"`
	AdsConsentSetAt         interface{}  `json:"adsConsentSetAt" db:"-"`
	AdsConsentReminderAt    interface{}  `json:"adsConsentReminderAt" db:"-"`
	ExperimentalFeatures    bool         `json:"experimentalFeatures" db:"-"`
	TwoFactorEnabled        bool         `json:"twoFactorEnabled" db:"-"`
	BackupCodesCreated      bool         `json:"backupCodesCreated" db:"-"`
}

type Subscription struct {
	Active         bool      `json:"active"`
	SubscribedAt   time.Time `json:"subscribedAt"`
	Status         string    `json:"status"`
	PaymentService string    `json:"paymentService"`
	Plan           string    `json:"plan"`
	Features       []string  `json:"features"`
}

type Profile struct {
	AutoSelectAudio              bool   `json:"autoSelectAudio"`
	DefaultAudioLanguage         string `json:"defaultAudioLanguage"`
	DefaultSubtitleLanguage      string `json:"defaultSubtitleLanguage"`
	AutoSelectSubtitle           int    `json:"autoSelectSubtitle"`
	DefaultSubtitleAccessibility int    `json:"defaultSubtitleAccessibility"`
	DefaultSubtitleForced        int    `json:"defaultSubtitleForced"`
}

type Services struct {
	Identifier string `json:"identifier"`
	Endpoint   string `json:"endpoint"`
	Token      string `json:"token,omitempty"`
	Status     string `json:"status"`
	Secret     string `json:"secret,omitempty"`
}
