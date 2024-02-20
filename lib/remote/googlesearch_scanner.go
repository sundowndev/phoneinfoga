package remote

import (
	"github.com/sundowndev/dorkgen"
	"github.com/sundowndev/dorkgen/googlesearch"
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

func NewGoogleSearchScanner() Scanner {
	return &googlesearchScanner{}
}

func (s *googlesearchScanner) Name() string {
	return Googlesearch
}

func (s *googlesearchScanner) Description() string {
	return "Generate several Google dork requests for a given phone number."
}

func (s *googlesearchScanner) DryRun(_ number.Number, _ ScannerOptions) error {
	return nil
}

func (s *googlesearchScanner) Run(n number.Number, _ ScannerOptions) (interface{}, error) {
	res := GoogleSearchResponse{
		SocialMedia:         getSocialMediaDorks(n),
		DisposableProviders: getDisposableProvidersDorks(n),
		Reputation:          getReputationDorks(n),
		Individuals:         getIndividualsDorks(n),
		General:             getGeneralDorks(n),
	}

	return res, nil
}

func getDisposableProvidersDorks(number number.Number) (results []*GoogleSearchDork) {
	var dorks = []*googlesearch.GoogleSearch{
		dorkgen.NewGoogleSearch().
			Site("hs3x.com").
			InText(number.International),
		dorkgen.NewGoogleSearch().
			Site("receive-sms-now.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("smslisten.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("smsnumbersonline.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("freesmscode.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("catchsms.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("smstibo.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("smsreceiving.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("getfreesmsnumber.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("sellaite.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("receive-sms-online.info").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("receivesmsonline.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("receive-a-sms.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("sms-receive.net").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("receivefreesms.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("receive-sms.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("receivetxt.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("freephonenum.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("freesmsverification.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("receive-sms-online.com").
			InText(number.International).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("smslive.co").
			InText(number.International).
			Or().
			InText(number.RawLocal),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.URL(),
		})
	}

	return results
}

func getIndividualsDorks(number number.Number) (results []*GoogleSearchDork) {
	var dorks = []*googlesearch.GoogleSearch{
		dorkgen.NewGoogleSearch().
			Site("numinfo.net").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("sync.me").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("whocallsyou.de").
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("pastebin.com").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("whycall.me").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("locatefamily.com").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("spytox.com").
			InText(number.RawLocal),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.URL(),
		})
	}

	return results
}

func getSocialMediaDorks(number number.Number) (results []*GoogleSearchDork) {
	var dorks = []*googlesearch.GoogleSearch{
		dorkgen.NewGoogleSearch().
			Site("facebook.com").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("twitter.com").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("linkedin.com").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("instagram.com").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("vk.com").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.URL(),
		})
	}

	return results
}

func getReputationDorks(number number.Number) (results []*GoogleSearchDork) {
	var dorks = []*googlesearch.GoogleSearch{
		dorkgen.NewGoogleSearch().
			Site("whosenumber.info").
			InText(number.E164).
			InTitle("who called"),
		dorkgen.NewGoogleSearch().
			InTitle("Phone Fraud").
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("findwhocallsme.com").
			InText(number.E164).
			Or().
			InText(number.International),
		dorkgen.NewGoogleSearch().
			Site("yellowpages.ca").
			InText(number.E164),
		dorkgen.NewGoogleSearch().
			Site("phonenumbers.ie").
			InText(number.E164),
		dorkgen.NewGoogleSearch().
			Site("who-calledme.com").
			InText(number.E164),
		dorkgen.NewGoogleSearch().
			Site("usphonesearch.net").
			InText(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("whocalled.us").
			InURL(number.RawLocal),
		dorkgen.NewGoogleSearch().
			Site("quinumero.info").
			InText(number.RawLocal).
			Or().
			InText(number.International),
		dorkgen.NewGoogleSearch().
			Site("uk.popularphotolook.com").
			InURL(number.RawLocal),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.URL(),
		})
	}

	return results
}

func getGeneralDorks(number number.Number) (results []*GoogleSearchDork) {
	var dorks = []*googlesearch.GoogleSearch{
		dorkgen.NewGoogleSearch().
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal).
			Or().
			InText(number.Local),
		dorkgen.NewGoogleSearch().
			Group(dorkgen.NewGoogleSearch().
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
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.URL(),
		})
	}

	return results
}
