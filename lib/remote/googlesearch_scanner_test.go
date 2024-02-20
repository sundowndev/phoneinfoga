package remote_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"testing"
)

func TestGoogleSearchScanner_Metadata(t *testing.T) {
	scanner := remote.NewGoogleSearchScanner()
	assert.Equal(t, remote.Googlesearch, scanner.Name())
	assert.NotEmpty(t, scanner.Description())
}

func TestGoogleSearchScanner(t *testing.T) {
	testcases := []struct {
		name       string
		number     *number.Number
		expected   map[string]interface{}
		wantErrors map[string]error
	}{
		{
			name: "successful scan",
			number: func() *number.Number {
				n, _ := number.NewNumber("15556661212")
				return n
			}(),
			expected: map[string]interface{}{
				"googlesearch": remote.GoogleSearchResponse{
					SocialMedia: []*remote.GoogleSearchDork{
						{
							Number: "+15556661212",
							Dork:   "site:facebook.com intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Afacebook.com+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:twitter.com intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Atwitter.com+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:linkedin.com intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Alinkedin.com+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:instagram.com intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Ainstagram.com+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:vk.com intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Avk.com+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
					},
					DisposableProviders: []*remote.GoogleSearchDork{
						{
							Number: "+15556661212",
							Dork:   "site:hs3x.com intext:\"15556661212\"",
							URL:    "https://www.google.com/search?q=site%3Ahs3x.com+intext%3A%2215556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:receive-sms-now.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Areceive-sms-now.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:smslisten.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Asmslisten.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:smsnumbersonline.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Asmsnumbersonline.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:freesmscode.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Afreesmscode.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						}, {
							Number: "+15556661212",
							Dork:   "site:catchsms.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Acatchsms.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:smstibo.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Asmstibo.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:smsreceiving.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Asmsreceiving.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:getfreesmsnumber.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Agetfreesmsnumber.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:sellaite.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Asellaite.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:receive-sms-online.info intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Areceive-sms-online.info+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:receivesmsonline.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Areceivesmsonline.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:receive-a-sms.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Areceive-a-sms.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:sms-receive.net intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Asms-receive.net+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:receivefreesms.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Areceivefreesms.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:receive-sms.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Areceive-sms.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:receivetxt.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Areceivetxt.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:freephonenum.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Afreephonenum.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:freesmsverification.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Afreesmsverification.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:receive-sms-online.com intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Areceive-sms-online.com+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:smslive.co intext:\"15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Asmslive.co+intext%3A%2215556661212%22+%7C+intext%3A%225556661212%22",
						},
					},
					Reputation: []*remote.GoogleSearchDork{
						{
							Number: "+15556661212",
							Dork:   "site:whosenumber.info intext:\"+15556661212\" intitle:\"who called\"",
							URL:    "https://www.google.com/search?q=site%3Awhosenumber.info+intext%3A%22%2B15556661212%22+intitle%3A%22who+called%22",
						},
						{
							Number: "+15556661212",
							Dork:   "intitle:\"Phone Fraud\" intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=intitle%3A%22Phone+Fraud%22+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:findwhocallsme.com intext:\"+15556661212\" | intext:\"15556661212\"",
							URL:    "https://www.google.com/search?q=site%3Afindwhocallsme.com+intext%3A%22%2B15556661212%22+%7C+intext%3A%2215556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:yellowpages.ca intext:\"+15556661212\"",
							URL:    "https://www.google.com/search?q=site%3Ayellowpages.ca+intext%3A%22%2B15556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:phonenumbers.ie intext:\"+15556661212\"",
							URL:    "https://www.google.com/search?q=site%3Aphonenumbers.ie+intext%3A%22%2B15556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:who-calledme.com intext:\"+15556661212\"",
							URL:    "https://www.google.com/search?q=site%3Awho-calledme.com+intext%3A%22%2B15556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:usphonesearch.net intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Ausphonesearch.net+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:whocalled.us inurl:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Awhocalled.us+inurl%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:quinumero.info intext:\"5556661212\" | intext:\"15556661212\"",
							URL:    "https://www.google.com/search?q=site%3Aquinumero.info+intext%3A%225556661212%22+%7C+intext%3A%2215556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:uk.popularphotolook.com inurl:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Auk.popularphotolook.com+inurl%3A%225556661212%22",
						},
					},
					Individuals: []*remote.GoogleSearchDork{
						{
							Number: "+15556661212",
							Dork:   "site:numinfo.net intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Anuminfo.net+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:sync.me intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Async.me+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:whocallsyou.de intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Awhocallsyou.de+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:pastebin.com intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Apastebin.com+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:whycall.me intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Awhycall.me+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:locatefamily.com intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Alocatefamily.com+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "site:spytox.com intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=site%3Aspytox.com+intext%3A%225556661212%22",
						},
					},
					General: []*remote.GoogleSearchDork{
						{
							Number: "+15556661212",
							Dork:   "intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\" | intext:\"(555) 666-1212\"",
							URL:    "https://www.google.com/search?q=intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22+%7C+intext%3A%22%28555%29+666-1212%22",
						},
						{
							Number: "+15556661212",
							Dork:   "(ext:doc | ext:docx | ext:odt | ext:pdf | ext:rtf | ext:sxw | ext:psw | ext:ppt | ext:pptx | ext:pps | ext:csv | ext:txt | ext:xls) intext:\"15556661212\" | intext:\"+15556661212\" | intext:\"5556661212\"",
							URL:    "https://www.google.com/search?q=%28ext%3Adoc+%7C+ext%3Adocx+%7C+ext%3Aodt+%7C+ext%3Apdf+%7C+ext%3Artf+%7C+ext%3Asxw+%7C+ext%3Apsw+%7C+ext%3Appt+%7C+ext%3Apptx+%7C+ext%3Apps+%7C+ext%3Acsv+%7C+ext%3Atxt+%7C+ext%3Axls%29+intext%3A%2215556661212%22+%7C+intext%3A%22%2B15556661212%22+%7C+intext%3A%225556661212%22",
						},
					},
				},
			},
			wantErrors: map[string]error{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			scanner := remote.NewGoogleSearchScanner()
			lib := remote.NewLibrary(filter.NewEngine())
			lib.AddScanner(scanner)

			if scanner.DryRun(*tt.number, remote.ScannerOptions{}) != nil {
				t.Fatal("DryRun() should return nil")
			}

			got, errs := lib.Scan(tt.number, remote.ScannerOptions{})
			if len(tt.wantErrors) > 0 {
				assert.Equal(t, tt.wantErrors, errs)
			} else {
				assert.Len(t, errs, 0)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}
