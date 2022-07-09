package cmd

import (
	"github.com/emersion/go-vcard"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

type GenVCardCmdOptions struct {
	Output   string
	Tel      string
	Fullname string
	Title    string
	Org      string
	Email    string
	Kind     string
}

func init() {
	opts := &GenVCardCmdOptions{}
	genVCardCmd := NewGenVCardCmd(opts)

	fl := genVCardCmd.Flags()
	fl.StringVarP(&opts.Output, "output", "o", "/dev/stdout", "Output file")
	fl.StringVar(&opts.Tel, "tel", "", "The phone number of the contact (E164 or international format)")
	fl.StringVar(&opts.Fullname, "fullname", "", "Name of the contact")
	fl.StringVar(&opts.Title, "title", "", "Title of the contact")
	fl.StringVar(&opts.Org, "org", "", "Organization of the contact")
	fl.StringVar(&opts.Email, "email", "", "Email of the contact")
	fl.StringVar(&opts.Kind, "kind", string(vcard.KindIndividual), "Kind of the contact: 'application', 'individual', 'group', 'location' or 'organization'; 'x-*' values may be used for experimental purposes.")

	rootCmd.AddCommand(genVCardCmd)
}

func NewGenVCardCmd(opts *GenVCardCmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gen-vcard",
		Example: "phoneinfoga gen-vcard -n +33678345233 --name \"John Doe\" -o cards.vcf",
		Short:   "Generate vcard v4 file for a given phone number",
		Run: func(cmd *cobra.Command, args []string) {
			destFile, err := os.Create(opts.Output)
			if err != nil {
				log.Fatal(err)
			}
			defer destFile.Close()

			enc := vcard.NewEncoder(destFile)

			card := vcard.Card{}
			card.SetValue(vcard.FieldFormattedName, opts.Fullname)
			card.SetValue(vcard.FieldName, strings.Join(strings.Split(opts.Fullname, " "), ";"))
			card.SetValue(vcard.FieldTelephone, opts.Tel)
			card.SetValue(vcard.FieldOrganization, opts.Org)
			card.SetValue(vcard.FieldTitle, opts.Title)
			card.SetValue(vcard.FieldEmail, opts.Email)
			card.SetValue(vcard.FieldKind, opts.Kind)

			vcard.ToV4(card)
			err = enc.Encode(card)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	return cmd
}
