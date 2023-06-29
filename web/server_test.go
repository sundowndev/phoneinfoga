package web

import (
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func performRequest(r http.Handler, method, path string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, err
}

func BenchmarkAPI(b *testing.B) {
	srv, err := NewServer(true)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("localScan - /api/numbers/:number/scan/local", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			res, err := performRequest(srv, http.MethodGet, "/api/numbers/3312345253/scan/local")
			assert.Equal(b, nil, err)
			assert.Equal(b, res.Result().StatusCode, 200)
		}
	})
}

func TestApi(t *testing.T) {
	srv, err := NewServer(false)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Serve", func(t *testing.T) {
		t.Run("detectContentType", func(t *testing.T) {
			contentType := detectContentType("/file.hash.css", []byte{})
			assert.Equal(t, "text/css; charset=utf-8", contentType)

			contentType = detectContentType("/file.hash.js", []byte{})
			assert.Equal(t, "text/javascript; charset=utf-8", contentType)

			contentType = detectContentType("/file.hash.svg", []byte{})
			assert.Equal(t, "image/svg+xml", contentType)

			contentType = detectContentType("/file.html", []byte("<html></html>"))
			assert.Equal(t, "text/html; charset=utf-8", contentType)
		})

		t.Run("getAllNumbers - /api/numbers", func(t *testing.T) {
			res, err := performRequest(srv, http.MethodGet, "/api/numbers")

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(t, err, nil)
			assert.Equal(t, res.Result().StatusCode, 200)
			assert.Equal(t, string(body), "{\"success\":true,\"numbers\":[]}")
		})

		t.Run("validate - /api/numbers/:number/validate", func(t *testing.T) {
			t.Run("valid number", func(t *testing.T) {
				res, err := performRequest(srv, http.MethodGet, "/api/numbers/3312345253/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(t, err, nil)
				assert.Equal(t, res.Result().StatusCode, 200)
				assert.Equal(t, string(body), "{\"success\":true,\"message\":\"The number is valid\"}")
			})

			t.Run("invalid number", func(t *testing.T) {
				res, err := performRequest(srv, http.MethodGet, "/api/numbers/azerty/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(t, err, nil)
				assert.Equal(t, res.Result().StatusCode, 400)
				assert.Equal(t, `{"success":false,"error":"the given phone number is not valid"}`, string(body))
			})

			t.Run("invalid country code", func(t *testing.T) {
				res, err := performRequest(srv, http.MethodGet, "/api/numbers/09880/validate")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(t, err, nil)
				assert.Equal(t, res.Result().StatusCode, 400)
				assert.Equal(t, string(body), "{\"success\":false,\"error\":\"invalid country code\"}")
			})
		})

		t.Run("localScan - /api/numbers/:number/scan/local", func(t *testing.T) {
			t.Run("valid number", func(t *testing.T) {
				res, err := performRequest(srv, http.MethodGet, "/api/numbers/3312345253/scan/local")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(t, err, nil)
				assert.Equal(t, res.Result().StatusCode, 200)
				assert.Equal(t, string(body), `{"success":true,"result":{"raw_local":"12345253","local":"12345253","e164":"+3312345253","international":"3312345253","country_code":33,"country":"FR"}}`)
			})

			t.Run("invalid number", func(t *testing.T) {
				res, err := performRequest(srv, http.MethodGet, "/api/numbers/9999999999/scan/local")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(t, err, nil)
				assert.Equal(t, res.Result().StatusCode, 400)
				assert.Equal(t, string(body), `{"success":false,"error":"invalid country code"}`)
			})
		})

		t.Run("numverifyScan - /api/numbers/:number/scan/numverify", func(t *testing.T) {
			t.Run("should succeed", func(t *testing.T) {
				defer gock.Off() // Flush pending mocks after test execution

				_ = os.Setenv("NUMVERIFY_API_KEY", "5ad5554ac240e4d3d31107941b35a5eb")
				defer os.Unsetenv("NUMVERIFY_API_KEY")

				number := "79516566591"

				expectedResult := suppliers.NumverifyValidateResponse{
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

				gock.New("https://api.apilayer.com").
					Get("/number_verification/validate").
					MatchHeader("Apikey", "5ad5554ac240e4d3d31107941b35a5eb").
					MatchParam("number", number).
					Reply(200).
					JSON(expectedResult)

				res, err := performRequest(srv, http.MethodGet, "/api/numbers/79516566591/scan/numverify")
				assert.Equal(t, nil, err)

				body, err := ioutil.ReadAll(res.Body)

				assert.Equal(t, nil, err)
				assert.Equal(t, 200, res.Result().StatusCode)
				assert.Equal(t, `{"success":true,"result":{"valid":true,"number":"79516566591","local_format":"9516566591","international_format":"+79516566591","country_prefix":"+7","country_code":"RU","country_name":"Russian Federation","location":"Saint Petersburg and Leningrad Oblast","carrier":"OJSC St. Petersburg Telecom (OJSC Tele2-Saint-Petersburg)","line_type":"mobile"}}`, string(body))

				assert.Equal(t, gock.IsDone(), true, "there should have no pending mocks")
			})

			t.Run("should handle error", func(t *testing.T) {
				defer gock.Off() // Flush pending mocks after test execution

				_ = os.Setenv("NUMVERIFY_API_KEY", "5ad5554ac240e4d3d31107941b35a5eb")
				defer os.Unsetenv("NUMVERIFY_API_KEY")

				number := "79516566591"

				expectedResult := &suppliers.NumverifyErrorResponse{
					Message: "You have exceeded your daily\\/monthly API rate limit. Please review and upgrade your subscription plan at https:\\/\\/apilayer.com\\/subscriptions to continue.",
				}

				gock.New("https://api.apilayer.com").
					Get("/number_verification/validate").
					MatchHeader("Apikey", "5ad5554ac240e4d3d31107941b35a5eb").
					MatchParam("number", number).
					Reply(429).
					JSON(expectedResult)

				res, err := performRequest(srv, http.MethodGet, "/api/numbers/79516566591/scan/numverify")
				assert.Equal(t, nil, err)

				body, err := ioutil.ReadAll(res.Body)

				assert.Equal(t, nil, err)
				assert.Equal(t, 500, res.Result().StatusCode)
				assert.Equal(t, `{"success":false,"error":"You have exceeded your daily\\/monthly API rate limit. Please review and upgrade your subscription plan at https:\\/\\/apilayer.com\\/subscriptions to continue."}`, string(body))

				assert.Equal(t, gock.IsDone(), true, "there should have no pending mocks")
			})
		})

		t.Run("googleSearchScan - /api/numbers/:number/scan/googlesearch", func(t *testing.T) {
			t.Run("should return google search dorks", func(t *testing.T) {
				res, err := performRequest(srv, http.MethodGet, "/api/numbers/330365179268/scan/googlesearch")
				assert.NoError(t, err)

				body, err := ioutil.ReadAll(res.Body)
				assert.NoError(t, err)

				assert.Equal(t, 200, res.Result().StatusCode)
				assert.Equal(t, `{"success":true,"result":{"social_media":[{"number":"+33365179268","dork":"site:facebook.com intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Afacebook.com+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:twitter.com intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Atwitter.com+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:linkedin.com intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Alinkedin.com+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:instagram.com intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Ainstagram.com+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:vk.com intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Avk.com+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"}],"disposable_providers":[{"number":"+33365179268","dork":"site:hs3x.com intext:\"33365179268\"","url":"https://www.google.com/search?q=site%3Ahs3x.com+intext%3A%2233365179268%22"},{"number":"+33365179268","dork":"site:receive-sms-now.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Areceive-sms-now.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:smslisten.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Asmslisten.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:smsnumbersonline.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Asmsnumbersonline.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:freesmscode.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Afreesmscode.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:catchsms.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Acatchsms.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:smstibo.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Asmstibo.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:smsreceiving.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Asmsreceiving.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:getfreesmsnumber.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Agetfreesmsnumber.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:sellaite.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Asellaite.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receive-sms-online.info intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Areceive-sms-online.info+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receivesmsonline.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Areceivesmsonline.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receive-a-sms.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Areceive-a-sms.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:sms-receive.net intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Asms-receive.net+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receivefreesms.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Areceivefreesms.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receive-sms.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Areceive-sms.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receivetxt.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Areceivetxt.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:freephonenum.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Afreephonenum.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:freesmsverification.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Afreesmsverification.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:receive-sms-online.com intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Areceive-sms-online.com+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:smslive.co intext:\"33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Asmslive.co+intext%3A%2233365179268%22+%7C+intext%3A%220365179268%22"}],"reputation":[{"number":"+33365179268","dork":"site:whosenumber.info intext:\"+33365179268\" intitle:\"who called\"","url":"https://www.google.com/search?q=site%3Awhosenumber.info+intext%3A%22%2B33365179268%22+intitle%3A%22who+called%22"},{"number":"+33365179268","dork":"intitle:\"Phone Fraud\" intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=intitle%3A%22Phone+Fraud%22+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:findwhocallsme.com intext:\"+33365179268\" | intext:\"33365179268\"","url":"https://www.google.com/search?q=site%3Afindwhocallsme.com+intext%3A%22%2B33365179268%22+%7C+intext%3A%2233365179268%22"},{"number":"+33365179268","dork":"site:yellowpages.ca intext:\"+33365179268\"","url":"https://www.google.com/search?q=site%3Ayellowpages.ca+intext%3A%22%2B33365179268%22"},{"number":"+33365179268","dork":"site:phonenumbers.ie intext:\"+33365179268\"","url":"https://www.google.com/search?q=site%3Aphonenumbers.ie+intext%3A%22%2B33365179268%22"},{"number":"+33365179268","dork":"site:who-calledme.com intext:\"+33365179268\"","url":"https://www.google.com/search?q=site%3Awho-calledme.com+intext%3A%22%2B33365179268%22"},{"number":"+33365179268","dork":"site:usphonesearch.net intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Ausphonesearch.net+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:whocalled.us inurl:\"0365179268\"","url":"https://www.google.com/search?q=site%3Awhocalled.us+inurl%3A%220365179268%22"},{"number":"+33365179268","dork":"site:quinumero.info intext:\"0365179268\" | intext:\"33365179268\"","url":"https://www.google.com/search?q=site%3Aquinumero.info+intext%3A%220365179268%22+%7C+intext%3A%2233365179268%22"},{"number":"+33365179268","dork":"site:uk.popularphotolook.com inurl:\"0365179268\"","url":"https://www.google.com/search?q=site%3Auk.popularphotolook.com+inurl%3A%220365179268%22"}],"individuals":[{"number":"+33365179268","dork":"site:numinfo.net intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Anuminfo.net+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:sync.me intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Async.me+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:whocallsyou.de intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Awhocallsyou.de+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:pastebin.com intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Apastebin.com+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:whycall.me intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Awhycall.me+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:locatefamily.com intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Alocatefamily.com+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"},{"number":"+33365179268","dork":"site:spytox.com intext:\"0365179268\"","url":"https://www.google.com/search?q=site%3Aspytox.com+intext%3A%220365179268%22"}],"general":[{"number":"+33365179268","dork":"intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\" | intext:\"03 65 17 92 68\"","url":"https://www.google.com/search?q=intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22+%7C+intext%3A%2203+65+17+92+68%22"},{"number":"+33365179268","dork":"(ext:doc | ext:docx | ext:odt | ext:pdf | ext:rtf | ext:sxw | ext:psw | ext:ppt | ext:pptx | ext:pps | ext:csv | ext:txt | ext:xls) intext:\"33365179268\" | intext:\"+33365179268\" | intext:\"0365179268\"","url":"https://www.google.com/search?q=%28ext%3Adoc+%7C+ext%3Adocx+%7C+ext%3Aodt+%7C+ext%3Apdf+%7C+ext%3Artf+%7C+ext%3Asxw+%7C+ext%3Apsw+%7C+ext%3Appt+%7C+ext%3Apptx+%7C+ext%3Apps+%7C+ext%3Acsv+%7C+ext%3Atxt+%7C+ext%3Axls%29+intext%3A%2233365179268%22+%7C+intext%3A%22%2B33365179268%22+%7C+intext%3A%220365179268%22"}]}}`, string(body))
			})
		})

		t.Run("ovhScan - /api/numbers/:number/scan/ovh", func(t *testing.T) {
			t.Run("should find number on OVH", func(t *testing.T) {
				defer gock.Off() // Flush pending mocks after test execution

				gock.New("https://api.ovh.com").
					Get("/1.0/telephony/number/detailedZones").
					MatchParam("country", "fr").
					Reply(200).
					JSON([]suppliers.OVHAPIResponseNumber{
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

				res, err := performRequest(srv, http.MethodGet, "/api/numbers/330365179268/scan/ovh")

				body, _ := ioutil.ReadAll(res.Body)

				assert.Equal(t, err, nil)
				assert.Equal(t, res.Result().StatusCode, 200)
				assert.Equal(t, string(body), `{"success":true,"result":{"found":true,"number_range":"036517xxxx","city":"Abbeville"}}`)

				assert.Equal(t, gock.IsDone(), true, "there should have no pending mocks")
			})
		})

		t.Run("healthHandler - /api/", func(t *testing.T) {
			res, err := performRequest(srv, http.MethodGet, "/api/")
			assert.Equal(t, nil, err)

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(t, 200, res.Result().StatusCode)
			assert.Equal(t, `{"success":true,"version":"dev","commit":"dev","demo":false}`, string(body))
		})

		t.Run("404 error - /api/notfound", func(t *testing.T) {
			res, err := performRequest(srv, http.MethodGet, "/api/notfound")
			assert.Equal(t, err, nil)

			body, _ := ioutil.ReadAll(res.Body)

			assert.Equal(t, res.Result().StatusCode, 404)
			assert.Equal(t, string(body), "{\"success\":false,\"error\":\"resource not found\"}")
		})

		t.Run("Client - /", func(t *testing.T) {
			res, err := performRequest(srv, http.MethodGet, "/")

			assert.Equal(t, nil, err)
			assert.Equal(t, 200, res.Result().StatusCode)
			assert.Equal(t, http.Header{"Content-Type": []string{"text/html; charset=utf-8"}}, res.Header())
		})
	})
}
