package dorkgen

// engineFactory is the main interface for
// search engine implementations.
type engineFactory interface {
	Site(string) *GoogleSearch
	ToString() string
	ToURL() string
	Intext(string) *GoogleSearch
	Inurl(string) *GoogleSearch
	Filetype(string) *GoogleSearch
	Cache(string) *GoogleSearch
	Related(string) *GoogleSearch
	Ext(string) *GoogleSearch
	Exclude(string) *GoogleSearch
	Group(string) *GoogleSearch
}
