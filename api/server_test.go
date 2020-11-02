package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	assertTest "github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"gopkg.in/sundowndev/phoneinfoga.v2/scanners"
)

var r *gin.Engine

func performRequest(r http.Handler, method, path string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, err
}

func BenchmarkAPI(b *testing.B) {
	assert := assertTest.New(b)
	r = gin.Default()
	r = Serve(r, true)

	b.Run("localScan - /api/numbers/:number/scan/local", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res, err := performRequest(r, "GET", "/api/numbers/3312345253/scan/local")
			assert.Equal(nil, err)
			assert.Equal(res.Result().StatusCode, 200, "should be equal")
		}
	})
}

func TestApi(t *testing.T) {
	assert := assertTest.New(t)
	r = gin.Default()
	r = Serve(r, false)

	t.Run("detectContentType", func(t *testing.T) {
		contentType := detectContentType("/file.hash.css", []byte{})
		assert.Equal("text/css", contentType, "should be equal")

		contentType = detectContentType("/file.hash.js", []byte{})
		assert.Equal("application/javascript", contentType, "should be equal")

		contentType = detectContentType("/file.hash.svg", []byte{})
		assert.Equal("image/svg+xml", contentType, "should be equal")

		contentType = detectContentType("/file.html", []byte("<html></html>"))
		assert.Equal("text/html; charset=utf-8", contentType, "should be equal")
	})

	t.Run("Serve", func(t *testing.T) {
		t.Run("getAllNumbers - /api/numbers", func(t *testing.T) {
			res, err := performRequest(r, "GET", "/api/numbers")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(err, nil, "should be equal")
			assert.Equal(res.Result().StatusCode, 200, "should be equal")
			assert.Equal(string(body), "{\"success\":true,\"error\":\"\",\"numbers\":[]}", "should be equal")
		})

		t.Run("validate - /api/numbers/:number/validate", func(t *testing.T) {
			t.Run("valid number", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/3312345253/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 200, "should be equal")
				assert.Equal(string(body), "{\"success\":true,\"error\":\"The number is valid\"}", "should be equal")
			})

			t.Run("invalid number", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/azerty/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 400, "should be equal")
				assert.Equal(string(body), "{\"success\":false,\"error\":\"Parameter 'number' must be a valid integer.\"}", "should be equal")
			})

			t.Run("invalid country code", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/09880/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 400, "should be equal")
				assert.Equal(string(body), "{\"success\":false,\"error\":\"invalid country code\"}", "should be equal")
			})
		})

		t.Run("localScan - /api/numbers/:number/scan/local", func(t *testing.T) {
			t.Run("valid number", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/3312345253/scan/local")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 200, "should be equal")
				assert.Equal(string(body), `{"success":true,"error":"","result":{"rawLocal":"12345253","local":"12345253","E164":"+3312345253","international":"3312345253","countryCode":33,"country":"FR","carrier":""}}`, "should be equal")
			})

			t.Run("invalid number", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/9999999999/scan/local")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 400, "should be equal")
				assert.Equal(string(body), `{"success":false,"error":"invalid country code"}`, "should be equal")
			})
		})

		t.Run("numverifyScan - /api/numbers/:number/scan/numverify", func(t *testing.T) {
			t.Run("should succeed", func(t *testing.T) {
				defer gock.Off() // Flush pending mocks after test execution

				expectedResult := scanners.NumverifyScannerResponse{
					Valid:               true,
					Number:              "79516566591",
					LocalFormat:         "9516566591",
					InternationalFormat: "+79516566591",
					CountryPrefix:       "+7",
					CountryCode:         "RU",
					CountryName:         "Russian Federation",
					Location:            "Saint Petersburg and Leningrad Oblast",
					Carrier:             "OJSC St. Petersburg Telecom (OJSC Tele2-Saint-Petersburg)",
					LineType:            "mobile",
				}

				gock.New("http://numverify.com").
					Get("/").
					Reply(200).BodyString(`<html><body><input type="hidden" name="scl_request_secret" value="secret"/></body></html>`)

				gock.New("https://numverify.com").
					Get("/php_helper_scripts/phone_api.php").
					MatchParam("secret_key", "5ad5554ac240e4d3d31107941b35a5eb").
					MatchParam("number", "79516566591").
					Reply(200).
					JSON(expectedResult)

				res, err := performRequest(r, "GET", "/api/numbers/79516566591/scan/numverify")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 200, "should be equal")
				assert.Equal(string(body), `{"success":true,"error":"","result":{"valid":true,"number":"79516566591","local_format":"9516566591","international_format":"+79516566591","country_prefix":"+7","country_code":"RU","country_name":"Russian Federation","location":"Saint Petersburg and Leningrad Oblast","carrier":"OJSC St. Petersburg Telecom (OJSC Tele2-Saint-Petersburg)","line_type":"mobile"}}`, "should be equal")

				assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
			})
		})

		t.Run("googleSearchScan - /api/numbers/:number/scan/googlesearch", func(t *testing.T) {
			t.Run("should return google search dorks", func(t *testing.T) {
				res, err := performRequest(r, "GET", "/api/numbers/330365179268/scan/googlesearch")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 200, "should be equal")
				assert.Equal(string(body), `{"success":true,"error":"","result":{"socialMedia":[{"number":"+33365179268","dork":"site:facebook.com intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Afacebook.com+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:twitter.com intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Atwitter.com+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:linkedin.com intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Alinkedin.com+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:instagram.com intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Ainstagram.com+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:vk.com intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Avk.com+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"}],"disposableProviders":[{"number":"+33365179268","dork":"site:hs3x.com intext:\"33365179268\"","URL":"https://www.google.com/search?q=site%3Ahs3x.com+intext%3A%2233365179268%22"},{"number":"+33365179268","dork":"site:receive-sms-now.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Areceive-sms-now.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:smslisten.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Asmslisten.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:smsnumbersonline.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Asmsnumbersonline.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:freesmscode.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Afreesmscode.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:catchsms.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Acatchsms.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:smstibo.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Asmstibo.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:smsreceiving.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Asmsreceiving.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:getfreesmsnumber.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Agetfreesmsnumber.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:sellaite.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Asellaite.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receive-sms-online.info intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Areceive-sms-online.info+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receivesmsonline.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Areceivesmsonline.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receive-a-sms.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Areceive-a-sms.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:sms-receive.net intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Asms-receive.net+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receivefreesms.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Areceivefreesms.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receive-sms.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Areceive-sms.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receivetxt.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Areceivetxt.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:freephonenum.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Afreephonenum.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:freesmsverification.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Afreesmsverification.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receive-sms-online.com intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Areceive-sms-online.com+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:smslive.co intext:\"33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Asmslive.co+intext%3A%2233365179268%22+OR+intext%3A%220365179268%22"}],"reputation":[{"number":"+33365179268","dork":"site:whosenumber.info intext:\"+33365179268\" intitle:\"who called\"","URL":"https://www.google.com/search?q=site%3Awhosenumber.info+intext%3A%22%2B33365179268%22+intitle%3A%22who+called%22"},{"number":"+33365179268","dork":"intitle:\"Phone Fraud\" intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=intitle%3A%22Phone+Fraud%22+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:findwhocallsme.com intext:\"+33365179268\" OR intext:\"33365179268\"","URL":"https://www.google.com/search?q=site%3Afindwhocallsme.com+intext%3A%22%2B33365179268%22+OR+intext%3A%2233365179268%22"},{"number":"+33365179268","dork":"site:yellowpages.ca intext:\"+33365179268\"","URL":"https://www.google.com/search?q=site%3Ayellowpages.ca+intext%3A%22%2B33365179268%22"},{"number":"+33365179268","dork":"site:phonenumbers.ie intext:\"+33365179268\"","URL":"https://www.google.com/search?q=site%3Aphonenumbers.ie+intext%3A%22%2B33365179268%22"},{"number":"+33365179268","dork":"site:who-calledme.com intext:\"+33365179268\"","URL":"https://www.google.com/search?q=site%3Awho-calledme.com+intext%3A%22%2B33365179268%22"},{"number":"+33365179268","dork":"site:usphonesearch.net intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Ausphonesearch.net+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:whocalled.us inurl:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Awhocalled.us+inurl%3A%220365179268%22"},{"number":"+33365179268","dork":"site:quinumero.info intext:\"0365179268\" OR intext:\"33365179268\"","URL":"https://www.google.com/search?q=site%3Aquinumero.info+intext%3A%220365179268%22+OR+intext%3A%2233365179268%22"},{"number":"+33365179268","dork":"site:uk.popularphotolook.com inurl:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Auk.popularphotolook.com+inurl%3A%220365179268%22"}],"individuals":[{"number":"+33365179268","dork":"site:numinfo.net intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Anuminfo.net+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:sync.me intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Async.me+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:whocallsyou.de intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Awhocallsyou.de+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:pastebin.com intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Apastebin.com+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:whycall.me intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Awhycall.me+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:locatefamily.com intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Alocatefamily.com+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:spytox.com intext:\"0365179268\"","URL":"https://www.google.com/search?q=site%3Aspytox.com+intext%3A%220365179268%22"}],"general":[{"number":"+33365179268","dork":"intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\" OR intext:\"03 65 17 92 68\"","URL":"https://www.google.com/search?q=intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22+OR+intext%3A%2203+65+17+92+68%22"},{"number":"+33365179268","dork":"(ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:xls) intext:\"33365179268\" OR intext:\"+33365179268\" OR intext:\"0365179268\"","URL":"https://www.google.com/search?q=%28ext%3Adoc+OR+ext%3Adocx+OR+ext%3Aodt+OR+ext%3Apdf+OR+ext%3Artf+OR+ext%3Asxw+OR+ext%3Apsw+OR+ext%3Appt+OR+ext%3Apptx+OR+ext%3Apps+OR+ext%3Acsv+OR+ext%3Atxt+OR+ext%3Axls%29+intext%3A%2233365179268%22+OR+intext%3A%22%2B33365179268%22+OR+intext%3A%220365179268%22"}]}}`, "should be equal")
			})
		})

		t.Run("ovhScan - /api/numbers/:number/scan/ovh", func(t *testing.T) {
			t.Run("should find number on OVH", func(t *testing.T) {
				defer gock.Off() // Flush pending mocks after test execution

				gock.New("https://api.ovh.com").
					Get("/1.0/telephony/number/detailedZones").
					MatchParam("country", "fr").
					Reply(200).
					JSON([]scanners.OVHAPIResponseNumber{
						{
							ZneList:             []string{},
							MatchingCriteria:    "",
							Prefix:              33,
							InternationalNumber: "003336517xxxx",
							Country:             "fr",
							ZipCode:             "",
							Number:              "036517xxxx",
							City:                "Abbeville",
							AskedCity:           "",
						},
					})

				res, err := performRequest(r, "GET", "/api/numbers/330365179268/scan/ovh")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(err, nil, "should be equal")
				assert.Equal(res.Result().StatusCode, 200, "should be equal")
				assert.Equal(string(body), `{"success":true,"error":"","result":{"found":true,"numberRange":"036517xxxx","city":"Abbeville","zipCode":""}}`, "should be equal")

				assert.Equal(gock.IsDone(), true, "there should have no pending mocks")
			})
		})

		t.Run("healthHandler - /api/", func(t *testing.T) {
			res, err := performRequest(r, "GET", "/api/")
			assert.Equal(nil, err, "should be equal")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(200, res.Result().StatusCode, "should be equal")
			assert.Equal("{\"success\":true,\"version\":\"unknown\",\"commit\":\"unknown\"}", string(body), "should be equal")
		})

		t.Run("404 error - /api/notfound", func(t *testing.T) {
			res, err := performRequest(r, "GET", "/api/notfound")
			assert.Equal(err, nil, "should be equal")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(res.Result().StatusCode, 404, "should be equal")
			assert.Equal(string(body), "{\"success\":false,\"error\":\"Resource not found\"}", "should be equal")
		})

		t.Run("Client - /", func(t *testing.T) {
			res, err := performRequest(r, "GET", "/")

			assert.Equal(nil, err, "should be equal")
			assert.Equal(200, res.Result().StatusCode, "should be equal")
			assert.Equal(http.Header{"Content-Type": []string{"text/html; charset=utf-8"}}, res.Header(), "should be equal")
		})
	})
}
