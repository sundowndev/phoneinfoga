package remote

import (
	"github.com/sundowndev/dorkgen"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
)

const Googlesearch = "googlesearch"

type googlesearchScanner struct{}

// GoogleSearchDork is the common format for dork requests
type GoogleSearchDork struct {
	Number string `json:"number" console:"-"`
	Dork   string `json:"dork" console:"-"`
	URL    string `json:"url" console:"URL"`
}

// GoogleSearchResponse is the output of Google search scanner.
// It contains all dorks created ordered by types.
type GoogleSearchResponse struct {
	SocialMedia         []*GoogleSearchDork `json:"social_media" console:"Social media,omitempty"`
	DisposableProviders []*GoogleSearchDork `json:"disposable_providers" console:"Disposable providers,omitempty"`
	Reputation          []*GoogleSearchDork `json:"reputation" console:"Reputation,omitempty"`
	Individuals         []*GoogleSearchDork `json:"individuals" console:"Individuals,omitempty"`
	General             []*GoogleSearchDork `json:"general" console:"General,omitempty"`
}

func NewGoogleSearchScanner() *googlesearchScanner {
	return &googlesearchScanner{}
}

func (s *googlesearchScanner) Identifier() string {
	return Googlesearch
}

func (s *googlesearchScanner) ShouldRun() bool {
	return true
}

func (s *googlesearchScanner) Scan(n *number.Number) (interface{}, error) {
	res := GoogleSearchResponse{
		SocialMedia:         getSocialMediaDorks(n),
		DisposableProviders: getDisposableProvidersDorks(n),
		Reputation:          getReputationDorks(n),
		Individuals:         getIndividualsDorks(n),
		General:             getGeneralDorks(n),
	}

	return res, nil
}

func getDisposableProvidersDorks(number *number.Number) (results []*GoogleSearchDork) {
	var dorks = []*dorkgen.GoogleSearch{
		(&dorkgen.GoogleSearch{}).
			Site("hs3x.com").
			Intext(number.International),
		(&dorkgen.GoogleSearch{}).
			Site("receive-sms-now.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("smslisten.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("smsnumbersonline.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("freesmscode.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("catchsms.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("smstibo.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("smsreceiving.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("getfreesmsnumber.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("sellaite.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("receive-sms-online.info").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("receivesmsonline.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("receive-a-sms.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("sms-receive.net").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("receivefreesms.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("receive-sms.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("receivetxt.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("freephonenum.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("freesmsverification.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("receive-sms-online.com").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("smslive.co").
			Intext(number.International).
			Or().
			Intext(number.RawLocal),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.ToURL(),
		})
	}

	return results
}

func getIndividualsDorks(number *number.Number) (results []*GoogleSearchDork) {
	var dorks = []*dorkgen.GoogleSearch{
		(&dorkgen.GoogleSearch{}).
			Site("numinfo.net").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("sync.me").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("whocallsyou.de").
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("pastebin.com").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("whycall.me").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("locatefamily.com").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("spytox.com").
			Intext(number.RawLocal),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.ToURL(),
		})
	}

	return results
}

func getSocialMediaDorks(number *number.Number) (results []*GoogleSearchDork) {
	var dorks = []*dorkgen.GoogleSearch{
		(&dorkgen.GoogleSearch{}).
			Site("facebook.com").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("twitter.com").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("linkedin.com").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("instagram.com").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("vk.com").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.ToURL(),
		})
	}

	return results
}

func getReputationDorks(number *number.Number) (results []*GoogleSearchDork) {
	var dorks = []*dorkgen.GoogleSearch{
		(&dorkgen.GoogleSearch{}).
			Site("whosenumber.info").
			Intext(number.E164).
			Intitle("who called"),
		(&dorkgen.GoogleSearch{}).
			Intitle("Phone Fraud").
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("findwhocallsme.com").
			Intext(number.E164).
			Or().
			Intext(number.International),
		(&dorkgen.GoogleSearch{}).
			Site("yellowpages.ca").
			Intext(number.E164),
		(&dorkgen.GoogleSearch{}).
			Site("phonenumbers.ie").
			Intext(number.E164),
		(&dorkgen.GoogleSearch{}).
			Site("who-calledme.com").
			Intext(number.E164),
		(&dorkgen.GoogleSearch{}).
			Site("usphonesearch.net").
			Intext(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("whocalled.us").
			Inurl(number.RawLocal),
		(&dorkgen.GoogleSearch{}).
			Site("quinumero.info").
			Intext(number.RawLocal).
			Or().
			Intext(number.International),
		(&dorkgen.GoogleSearch{}).
			Site("uk.popularphotolook.com").
			Inurl(number.RawLocal),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.ToURL(),
		})
	}

	return results
}

func getGeneralDorks(number *number.Number) (results []*GoogleSearchDork) {
	var dorks = []*dorkgen.GoogleSearch{
		(&dorkgen.GoogleSearch{}).
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal).
			Or().
			Intext(number.Local),
		(&dorkgen.GoogleSearch{}).
			Group((&dorkgen.GoogleSearch{}).
				Ext("doc").
				Or().
				Ext("docx").
				Or().
				Ext("odt").
				Or().
				Ext("pdf").
				Or().
				Ext("rtf").
				Or().
				Ext("sxw").
				Or().
				Ext("psw").
				Or().
				Ext("ppt").
				Or().
				Ext("pptx").
				Or().
				Ext("pps").
				Or().
				Ext("csv").
				Or().
				Ext("txt").
				Or().
				Ext("xls")).
			Intext(number.International).
			Or().
			Intext(number.E164).
			Or().
			Intext(number.RawLocal),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.ToURL(),
		})
	}

	return results
}
