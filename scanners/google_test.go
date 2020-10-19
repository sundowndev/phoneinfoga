package scanners

import (
	"testing"

	assertTest "github.com/stretchr/testify/assert"
	"gopkg.in/sundowndev/phoneinfoga.v2/utils"
)

func TestGoogleSearchScan(t *testing.T) {
	assert := assertTest.New(t)

	number, _ := LocalScan("+33 673421322")

	scan := googlesearchScanCLI(utils.LoggerService, number)
	scanWithFormat := googlesearchScanCLI(utils.LoggerService, number, "06.73.42.13.22")

	t.Run("getDisposableProvidersDorks", func(t *testing.T) {
		t.Run("should generate disposable provider dorks", func(t *testing.T) {
			expectedResult := []*GoogleSearchDork{
				{
					Number: "+33673421322",
					Dork:   "site:hs3x.com intext:\"33673421322\"",
					URL:    "https://www.google.com/search?q=site%3Ahs3x.com+intext%3A%2233673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receive-sms-now.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceive-sms-now.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:smslisten.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asmslisten.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:smsnumbersonline.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asmsnumbersonline.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:freesmscode.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Afreesmscode.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:catchsms.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Acatchsms.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:smstibo.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asmstibo.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:smsreceiving.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asmsreceiving.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:getfreesmsnumber.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Agetfreesmsnumber.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:sellaite.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asellaite.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receive-sms-online.info intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceive-sms-online.info+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receivesmsonline.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceivesmsonline.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receive-a-sms.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceive-a-sms.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:sms-receive.net intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asms-receive.net+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receivefreesms.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceivefreesms.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receive-sms.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceive-sms.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receivetxt.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceivetxt.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:freephonenum.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Afreephonenum.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:freesmsverification.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Afreesmsverification.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receive-sms-online.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceive-sms-online.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:smslive.co intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asmslive.co+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
			}

			assert.Equal(expectedResult, scan.DisposableProviders, "they should be equal")
		})

		t.Run("should generate disposable provider dorks with additional format", func(t *testing.T) {
			expectedResult := []*GoogleSearchDork{
				{
					Number: "+33673421322",
					Dork:   "site:hs3x.com intext:\"33673421322\"",
					URL:    "https://www.google.com/search?q=site%3Ahs3x.com+intext%3A%2233673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receive-sms-now.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceive-sms-now.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:smslisten.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asmslisten.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:smsnumbersonline.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asmsnumbersonline.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:freesmscode.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Afreesmscode.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:catchsms.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Acatchsms.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:smstibo.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asmstibo.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:smsreceiving.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asmsreceiving.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:getfreesmsnumber.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Agetfreesmsnumber.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:sellaite.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asellaite.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receive-sms-online.info intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceive-sms-online.info+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receivesmsonline.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceivesmsonline.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receive-a-sms.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceive-a-sms.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:sms-receive.net intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asms-receive.net+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receivefreesms.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceivefreesms.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receive-sms.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceive-sms.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receivetxt.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceivetxt.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:freephonenum.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Afreephonenum.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:freesmsverification.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Afreesmsverification.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:receive-sms-online.com intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Areceive-sms-online.com+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:smslive.co intext:\"33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Asmslive.co+intext%3A%2233673421322%22+OR+intext%3A%220673421322%22",
				},
			}

			assert.Equal(expectedResult, scanWithFormat.DisposableProviders, "they should be equal")
		})
	})

	t.Run("getSocialMediaDorks", func(t *testing.T) {
		t.Run("should generate social media dorks", func(t *testing.T) {
			expectedResult := []*GoogleSearchDork{
				{
					Number: "+33673421322",
					Dork:   "site:facebook.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Afacebook.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:twitter.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Atwitter.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:linkedin.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Alinkedin.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:instagram.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Ainstagram.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:vk.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Avk.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
			}

			assert.Equal(expectedResult, scan.SocialMedia, "they should be equal")
		})

		t.Run("should generate social media dorks with additional format", func(t *testing.T) {
			expectedResult := []*GoogleSearchDork{
				{
					Number: "+33673421322",
					Dork:   "site:facebook.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Afacebook.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:twitter.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Atwitter.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:linkedin.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Alinkedin.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:instagram.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Ainstagram.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:vk.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Avk.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
			}

			assert.Equal(expectedResult, scanWithFormat.SocialMedia, "they should be equal")
		})
	})

	t.Run("getReputationDorks", func(t *testing.T) {
		t.Run("should generate reputation dorks", func(t *testing.T) {
			expectedResult := []*GoogleSearchDork{
				{
					Number: "+33673421322",
					Dork:   "site:whosenumber.info intext:\"+33673421322\" intitle:\"who called\"",
					URL:    "https://www.google.com/search?q=site%3Awhosenumber.info+intext%3A%22%2B33673421322%22+intitle%3A%22who+called%22",
				},
				{
					Number: "+33673421322",
					Dork:   "intitle:\"Phone Fraud\" intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=intitle%3A%22Phone+Fraud%22+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:findwhocallsme.com intext:\"+33673421322\" OR intext:\"33673421322\"",
					URL:    "https://www.google.com/search?q=site%3Afindwhocallsme.com+intext%3A%22%2B33673421322%22+OR+intext%3A%2233673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:yellowpages.ca intext:\"+33673421322\"",
					URL:    "https://www.google.com/search?q=site%3Ayellowpages.ca+intext%3A%22%2B33673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:phonenumbers.ie intext:\"+33673421322\"",
					URL:    "https://www.google.com/search?q=site%3Aphonenumbers.ie+intext%3A%22%2B33673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:who-calledme.com intext:\"+33673421322\"",
					URL:    "https://www.google.com/search?q=site%3Awho-calledme.com+intext%3A%22%2B33673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:usphonesearch.net intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Ausphonesearch.net+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:whocalled.us inurl:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Awhocalled.us+inurl%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:quinumero.info intext:\"0673421322\" OR intext:\"33673421322\"",
					URL:    "https://www.google.com/search?q=site%3Aquinumero.info+intext%3A%220673421322%22+OR+intext%3A%2233673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:uk.popularphotolook.com inurl:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Auk.popularphotolook.com+inurl%3A%220673421322%22",
				},
			}

			assert.Equal(expectedResult, scan.Reputation, "they should be equal")
		})

		t.Run("should generate reputation dorks with additional format", func(t *testing.T) {
			expectedResult := []*GoogleSearchDork{
				{
					Number: "+33673421322",
					Dork:   "site:whosenumber.info intext:\"+33673421322\" intitle:\"who called\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Awhosenumber.info+intext%3A%22%2B33673421322%22+intitle%3A%22who+called%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "intitle:\"Phone Fraud\" intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=intitle%3A%22Phone+Fraud%22+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:findwhocallsme.com intext:\"+33673421322\" OR intext:\"33673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Afindwhocallsme.com+intext%3A%22%2B33673421322%22+OR+intext%3A%2233673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:yellowpages.ca intext:\"+33673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Ayellowpages.ca+intext%3A%22%2B33673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:phonenumbers.ie intext:\"+33673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Aphonenumbers.ie+intext%3A%22%2B33673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:who-calledme.com intext:\"+33673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Awho-calledme.com+intext%3A%22%2B33673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:usphonesearch.net intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Ausphonesearch.net+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:whocalled.us inurl:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Awhocalled.us+inurl%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:quinumero.info intext:\"0673421322\" OR intext:\"33673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Aquinumero.info+intext%3A%220673421322%22+OR+intext%3A%2233673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:uk.popularphotolook.com inurl:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Auk.popularphotolook.com+inurl%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
			}

			assert.Equal(expectedResult, scanWithFormat.Reputation, "they should be equal")
		})
	})

	t.Run("getIndividualsDorks", func(t *testing.T) {
		t.Run("should generate individual dorks", func(t *testing.T) {
			expectedResult := []*GoogleSearchDork{
				{
					Number: "+33673421322",
					Dork:   "site:numinfo.net intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Anuminfo.net+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:sync.me intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Async.me+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:whocallsyou.de intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Awhocallsyou.de+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:pastebin.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Apastebin.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:whycall.me intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Awhycall.me+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:locatefamily.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Alocatefamily.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:spytox.com intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=site%3Aspytox.com+intext%3A%220673421322%22",
				},
			}

			assert.Equal(expectedResult, scan.Individuals, "they should be equal")
		})

		t.Run("should generate individual dorks with additional format", func(t *testing.T) {
			expectedResult := []*GoogleSearchDork{
				{
					Number: "+33673421322",
					Dork:   "site:numinfo.net intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Anuminfo.net+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:sync.me intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Async.me+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:whocallsyou.de intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Awhocallsyou.de+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:pastebin.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Apastebin.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:whycall.me intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Awhycall.me+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:locatefamily.com intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Alocatefamily.com+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "site:spytox.com intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=site%3Aspytox.com+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
			}

			assert.Equal(expectedResult, scanWithFormat.Individuals, "they should be equal")
		})
	})

	t.Run("getGeneralDorks", func(t *testing.T) {
		t.Run("should generate general dorks", func(t *testing.T) {
			expectedResult := []*GoogleSearchDork{
				{
					Number: "+33673421322",
					Dork:   `intext:"33673421322" OR intext:"+33673421322" OR intext:"0673421322" OR intext:"06 73 42 13 22"`,
					URL:    "https://www.google.com/search?q=intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206+73+42+13+22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "(ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:xls) intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\"",
					URL:    "https://www.google.com/search?q=%28ext%3Adoc+OR+ext%3Adocx+OR+ext%3Aodt+OR+ext%3Apdf+OR+ext%3Artf+OR+ext%3Asxw+OR+ext%3Apsw+OR+ext%3Appt+OR+ext%3Apptx+OR+ext%3Apps+OR+ext%3Acsv+OR+ext%3Atxt+OR+ext%3Axls%29+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22",
				},
			}

			assert.Equal(expectedResult, scan.General, "they should be equal")
		})

		t.Run("should generate general dorks with additional format", func(t *testing.T) {
			expectedResult := []*GoogleSearchDork{
				{
					Number: "+33673421322",
					Dork:   `intext:"33673421322" OR intext:"+33673421322" OR intext:"0673421322" OR intext:"06 73 42 13 22" OR intext:"06.73.42.13.22"`,
					URL:    "https://www.google.com/search?q=intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206+73+42+13+22%22+OR+intext%3A%2206.73.42.13.22%22",
				},
				{
					Number: "+33673421322",
					Dork:   "(ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:xls) intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\" OR intext:\"06.73.42.13.22\"",
					URL:    "https://www.google.com/search?q=%28ext%3Adoc+OR+ext%3Adocx+OR+ext%3Aodt+OR+ext%3Apdf+OR+ext%3Artf+OR+ext%3Asxw+OR+ext%3Apsw+OR+ext%3Appt+OR+ext%3Apptx+OR+ext%3Apps+OR+ext%3Acsv+OR+ext%3Atxt+OR+ext%3Axls%29+intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206.73.42.13.22%22",
				},
			}

			assert.Equal(expectedResult, scanWithFormat.General, "they should be equal")
		})
	})
}
