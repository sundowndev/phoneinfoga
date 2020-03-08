package phonenumbers

// merge merges two number formats
func (nf *NumberFormat) merge(other *NumberFormat) {
	if other.Pattern != nil {
		nf.Pattern = other.Pattern
	}
	if other.Format != nil {
		nf.Format = other.Format
	}
	for i := 0; i < len(other.LeadingDigitsPattern); i++ {
		nf.LeadingDigitsPattern = append(nf.LeadingDigitsPattern, other.LeadingDigitsPattern[i])
	}
	if other.NationalPrefixFormattingRule != nil {
		nf.NationalPrefixFormattingRule = other.NationalPrefixFormattingRule
	}
	if other.DomesticCarrierCodeFormattingRule != nil {
		nf.DomesticCarrierCodeFormattingRule = other.DomesticCarrierCodeFormattingRule
	}
	if other.NationalPrefixOptionalWhenFormatting != nil {
		nf.NationalPrefixOptionalWhenFormatting = other.NationalPrefixOptionalWhenFormatting
	}
}

func (pd *PhoneNumberDesc) hasPossibleLength(length int32) bool {
	for _, l := range pd.PossibleLength {
		if l == length {
			return true
		}
	}

	return false
}

func (pd *PhoneNumberDesc) hasPossibleLengthLocalOnly(length int32) bool {
	for _, l := range pd.PossibleLengthLocalOnly {
		if l == length {
			return true
		}
	}
	return false
}
