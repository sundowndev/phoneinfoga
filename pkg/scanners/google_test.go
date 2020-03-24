package scanners

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogleSearchScan(t *testing.T) {
	assert := assert.New(t)

	t.Run("getDisposableProvidersDorks", func(t *testing.T) {})
	t.Run("getIndividualsDorks", func(t *testing.T) {})
	t.Run("getSocialMediaDorks", func(t *testing.T) {})
	t.Run("getReputationDorks", func(t *testing.T) {})

	t.Run("getGeneralDorks", func(t *testing.T) {
		t.Run("should generate general dorks", func(t *testing.T) {
			expectedResult := GoogleSearchResponse{
				General: []*GoogleSearchDork{
					&GoogleSearchDork{
						Number: "+33673421322",
						Dork:   `intext:"33673421322" OR intext:"+33673421322" OR intext:"0673421322" OR intext:"06 73 42 13 22"`,
						URL:    "https://www.google.com/search?q=intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22+OR+intext%3A%2206+73+42+13+22%22",
					},
					&GoogleSearchDork{
						Number: "+33673421322",
						Dork:   "(intext:\"33673421322\" OR intext:\"+33673421322\" OR intext:\"0673421322\") ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:xls",
						URL:    "https://www.google.com/search?q=%28intext%3A%2233673421322%22+OR+intext%3A%22%2B33673421322%22+OR+intext%3A%220673421322%22%29+ext%3Adoc+OR+ext%3Adocx+OR+ext%3Aodt+OR+ext%3Apdf+OR+ext%3Artf+OR+ext%3Asxw+OR+ext%3Apsw+OR+ext%3Appt+OR+ext%3Apptx+OR+ext%3Apps+OR+ext%3Acsv+OR+ext%3Atxt+OR+ext%3Axls",
					},
				},
			}

			number, _ := LocalScan("+33 673421322")

			scan := GoogleSearchScan(number)

			assert.Equal(scan.General, expectedResult.General, "they should be equal")
		})
	})
}
