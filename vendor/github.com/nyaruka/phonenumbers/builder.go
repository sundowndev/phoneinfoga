package phonenumbers

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ----------------------------------------------------------------------------
// Golang port of:
// https://github.com/googlei18n/libphonenumber/blob/master/tools/java/common/src/com/google/i18n/phonenumbers/BuildMetadataFromXml.java
// ----------------------------------------------------------------------------

func sp(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func bp(value bool) *bool {
	return &value
}

func ip(value int32) *int32 {
	return &value
}

func BuildPhoneMetadataCollection(inputXML []byte, liteBuild bool, specialBuild bool) (*PhoneMetadataCollection, error) {
	metadata := &PhoneNumberMetadataE{}
	err := xml.Unmarshal(inputXML, metadata)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling XML: %s", err))
	}
	isShortNumberMetadata := false
	isAlternateFormatsMetadata := false
	return buildPhoneMetadataFromElement(metadata, liteBuild, specialBuild, isShortNumberMetadata, isAlternateFormatsMetadata)
}

func buildPhoneMetadataFromElement(document *PhoneNumberMetadataE, liteBuild bool, specialBuild bool, isShortNumberMetadata bool, isAlternateFormatsMetadata bool) (*PhoneMetadataCollection, error) {
	collection := PhoneMetadataCollection{}
	numOfTerritories := len(document.Territories)
	for i := 0; i < numOfTerritories; i++ {
		territoryElement := document.Territories[i]
		regionCode := territoryElement.ID

		metadata := loadCountryMetadata(regionCode, &territoryElement, isShortNumberMetadata, isAlternateFormatsMetadata)
		collection.Metadata = append(collection.Metadata, metadata)
	}
	return &collection, nil
}

// Build a mapping from a country calling code to the region codes which denote the country/region
// represented by that country code. In the case of multiple countries sharing a calling code,
// such as the NANPA countries, the one indicated with "isMainCountryForCode" in the metadata
// should be first.
func BuildCountryCodeToRegionMap(metadataCollection *PhoneMetadataCollection) map[int][]string {
	countryCodeToRegionCodeMap := make(map[int][]string)
	for _, metadata := range metadataCollection.Metadata {
		regionCode := metadata.GetId()
		countryCode := int(metadata.GetCountryCode())
		_, present := countryCodeToRegionCodeMap[countryCode]
		if present {
			phoneList := countryCodeToRegionCodeMap[countryCode]
			if metadata.GetMainCountryForCode() {
				phoneList = append([]string{regionCode}, phoneList...)
			} else {
				phoneList = append(phoneList, regionCode)
			}
			countryCodeToRegionCodeMap[countryCode] = phoneList
		} else {
			// For most countries, there will be only one region code for the country calling code.
			phoneList := []string{}
			if regionCode != "" { // For alternate formats, there are no region codes at all.
				phoneList = append(phoneList, regionCode)
			}
			countryCodeToRegionCodeMap[countryCode] = phoneList
		}
	}
	return countryCodeToRegionCodeMap
}

func validateRE(re string, removeWhitespace bool) string {
	// Removes all the whitespace and newline from the regexp. Not Ming pattern compile options to
	// make it work across programming languages.
	if removeWhitespace {
		re = string(regexp.MustCompile(`\s`).ReplaceAllLiteralString(re, ""))
	}
	_, err := regexp.Compile(re)
	if err != nil {
		panic(err)
	}
	return re
}

func loadTerritoryTagMetadata(regionCode string, territory *TerritoryE, nationalPrefix string) *PhoneMetadata {
	metadata := &PhoneMetadata{}
	metadata.Id = sp(regionCode)

	if territory.CountryCode != 0 {
		metadata.CountryCode = ip(territory.CountryCode)
	}
	if territory.LeadingDigits != "" {
		metadata.LeadingDigits = sp(validateRE(territory.LeadingDigits, false))
	}
	if territory.InternationalPrefix != "" {
		metadata.InternationalPrefix = sp(validateRE(territory.InternationalPrefix, false))
	}
	if territory.PreferredInternationalPrefix != "" {
		metadata.PreferredInternationalPrefix = sp(territory.PreferredInternationalPrefix)
	}
	if territory.NationalPrefixForParsing != "" {
		metadata.NationalPrefixForParsing = sp(validateRE(territory.NationalPrefixForParsing, true))
		if territory.NationalPrefixTransformRule != "" {
			metadata.NationalPrefixTransformRule = sp(validateRE(territory.NationalPrefixTransformRule, false))
		}
	}
	if nationalPrefix != "" {
		metadata.NationalPrefix = sp(nationalPrefix)
		if metadata.NationalPrefixForParsing == nil {
			metadata.NationalPrefixForParsing = sp(nationalPrefix)
		}
	}
	if territory.PreferredExtnPrefix != "" {
		metadata.PreferredExtnPrefix = sp(territory.PreferredExtnPrefix)
	}
	if territory.MainCountryForCode {
		metadata.MainCountryForCode = bp(true)
	}
	if territory.MobileNumberPortableRegion {
		metadata.MobileNumberPortableRegion = bp(true)
	}
	return metadata
}

func setLeadingDigitsPatterns(numberFormatElement *NumberFormatE, format *NumberFormat) {
	if len(numberFormatElement.LeadingDigits) > 0 {
		for i := 0; i < len(numberFormatElement.LeadingDigits); i++ {
			format.LeadingDigitsPattern = append(format.LeadingDigitsPattern, validateRE(numberFormatElement.LeadingDigits[i], true))
		}
	}
}

/**
 * Extracts the pattern for international format. If there is no intlFormat, default to using the
 * national format. If the intlFormat is set to "NA" the intlFormat should be ignored.
 *
 * @throws  RuntimeException if multiple intlFormats have been encountered.
 * @return  whether an international number format is defined.
 */
func loadInternationalFormat(metadata *PhoneMetadata, numberFormatElement *NumberFormatE, nationalFormat *NumberFormat) bool {
	intlFormat := &NumberFormat{}
	intlFormatPattern := numberFormatElement.InternationalFormat
	hasExplicitIntlFormatDefined := false

	if len(intlFormatPattern) > 1 {
		panic("Invalid number of intlFormat patterns for country: " + metadata.GetId())

	} else if len(intlFormatPattern) == 0 {
		// Default to use the same as the national pattern if none is defined.
		intlFormat.merge(nationalFormat)
	} else {
		intlFormat.Pattern = sp(numberFormatElement.Pattern)
		setLeadingDigitsPatterns(numberFormatElement, intlFormat)
		intlFormatPatternValue := intlFormatPattern[0]
		if intlFormatPatternValue != "NA" {
			intlFormat.Format = sp(intlFormatPatternValue)
		}
		hasExplicitIntlFormatDefined = true
	}

	if intlFormat.Format != nil {
		metadata.IntlNumberFormat = append(metadata.IntlNumberFormat, intlFormat)
	}
	return hasExplicitIntlFormatDefined
}

/**
 * Extracts the pattern for the national format.
 *
 * @throws  RuntimeException if multiple or no formats have been encountered.
 */
// @VisibleForTesting
func loadNationalFormat(metadata *PhoneMetadata, numberFormatElement *NumberFormatE, format *NumberFormat) {
	setLeadingDigitsPatterns(numberFormatElement, format)
	format.Pattern = sp(validateRE(numberFormatElement.Pattern, false))
	format.Format = sp(numberFormatElement.Format)
}

func getDomesticCarrierCodeFormattingRule(carrierCodeFormattingRule string, nationalPrefix string) string {
	// Replace $FG with the first group ($1) and $NP with the national prefix.
	carrierCodeFormattingRule = strings.Replace(carrierCodeFormattingRule, "$FG", "$1", 1)
	carrierCodeFormattingRule = strings.Replace(carrierCodeFormattingRule, "$NP", nationalPrefix, 1)
	return carrierCodeFormattingRule
}

func getNationalPrefixFormattingRule(nationalPrefixFormattingRule string, nationalPrefix string) string {
	// Replace $NP with national prefix and $FG with the first group ($1).
	nationalPrefixFormattingRule = strings.Replace(nationalPrefixFormattingRule, "$NP", nationalPrefix, 1)
	nationalPrefixFormattingRule = strings.Replace(nationalPrefixFormattingRule, "$FG", "$1", 1)
	return nationalPrefixFormattingRule
}

/**
 * Extracts the available formats from the provided DOM element. If it does not contain any
 * nationalPrefixFormattingRule, the one passed-in is retained; similarly for
 * nationalPrefixOptionalWhenFormatting. The nationalPrefix, nationalPrefixFormattingRule and
 * nationalPrefixOptionalWhenFormatting values are provided from the parent (territory) element.
 */
// @VisibleForTesting
func loadAvailableFormats(metadata *PhoneMetadata, element *TerritoryE, nationalPrefix string,
	nationalPrefixFormattingRule string, nationalPrefixOptionalWhenFormatting bool) {
	carrierCodeFormattingRule := ""
	if element.CarrierCodeFormattingRule != "" {
		carrierCodeFormattingRule = validateRE(getDomesticCarrierCodeFormattingRule(element.CarrierCodeFormattingRule, nationalPrefix), false)
	}
	numberFormatElements := element.AvailableFormats
	hasExplicitIntlFormatDefined := false

	if len(numberFormatElements) > 0 {
		for i := 0; i < len(numberFormatElements); i++ {
			numberFormatElement := numberFormatElements[i]
			format := NumberFormat{}

			if numberFormatElement.NationalPrefixFormattingRule != "" {
				format.NationalPrefixFormattingRule = sp(getNationalPrefixFormattingRule(numberFormatElement.NationalPrefixFormattingRule, nationalPrefix))
			} else {
				format.NationalPrefixFormattingRule = sp(nationalPrefixFormattingRule)
			}

			if numberFormatElement.NationalPrefixOptionalWhenFormatting != nil {
				format.NationalPrefixOptionalWhenFormatting = numberFormatElement.NationalPrefixOptionalWhenFormatting
			} else if nationalPrefixOptionalWhenFormatting {
				format.NationalPrefixOptionalWhenFormatting = bp(nationalPrefixOptionalWhenFormatting)
			}

			if numberFormatElement.CarrierCodeFormattingRule != "" {
				format.DomesticCarrierCodeFormattingRule = sp(validateRE(getDomesticCarrierCodeFormattingRule(numberFormatElement.CarrierCodeFormattingRule, nationalPrefix), false))
			} else if carrierCodeFormattingRule != "" {
				format.DomesticCarrierCodeFormattingRule = sp(carrierCodeFormattingRule)
			}
			loadNationalFormat(metadata, &numberFormatElement, &format)
			metadata.NumberFormat = append(metadata.NumberFormat, &format)

			if loadInternationalFormat(metadata, &numberFormatElement, &format) {
				hasExplicitIntlFormatDefined = true
			}
		}
		// Only a small number of regions need to specify the intlFormats in the xml. For the majority
		// of countries the intlNumberFormat metadata is an exact copy of the national NumberFormat
		// metadata. To minimize the size of the metadata file, we only keep intlNumberFormats that
		// actually differ in some way to the national formats.
		if !hasExplicitIntlFormatDefined {
			metadata.IntlNumberFormat = nil
		}
	}
}

/**
 * Checks if the possible lengths provided as a sorted set are equal to the possible lengths
 * stored already in the description pattern. Note that possibleLengths may be empty but must not
 * be null, and the PhoneNumberDesc passed in should also not be null.
 */
func arePossibleLengthsEqual(possibleLengths map[int32]bool, desc *PhoneNumberDesc) bool {
	if len(possibleLengths) != len(desc.PossibleLength) {
		return false
	}

	// check whether the same elements exist
	for _, val := range desc.PossibleLength {
		_, exists := possibleLengths[val]
		if !exists {
			return false
		}
	}
	return true
}

/**
 * Parses a possible length string into a set of the integers that are covered.
 *
 * @param possibleLengthString  a string specifying the possible lengths of phone numbers. Follows
 *     this syntax: ranges or elements are separated by commas, and ranges are specified in
 *     [min-max] notation, inclusive. For example, [3-5],7,9,[11-14] should be parsed to
 *     3,4,5,7,9,11,12,13,14.
 */
func parsePossibleLengthStringToSet(possibleLengthString string) map[int32]bool {
	if possibleLengthString == "" {
		panic("Empty possibleLength string found.")
	}
	lengths := strings.Split(possibleLengthString, ",")
	lengthSet := make(map[int32]bool)

	for i := 0; i < len(lengths); i++ {
		lengthSubstring := lengths[i]
		if lengthSubstring == "" {
			panic("Leading, trailing or adjacent commas in possible length string %s, these should only separate numbers or ranges.")
		} else if lengthSubstring[0] == '[' {
			if lengthSubstring[len(lengthSubstring)-1] != ']' {
				panic(fmt.Sprintf("Missing end of range character in possible length string %s.", possibleLengthString))
			}
			// Strip the leading and trailing [], and split on the -.
			minMax := strings.Split(lengthSubstring[1:len(lengthSubstring)-1], "-")
			if len(minMax) != 2 {
				panic(fmt.Sprintf("Ranges must have exactly one - character: missing for %s.", possibleLengthString))
			}
			min, _ := strconv.Atoi(minMax[0])
			max, _ := strconv.Atoi(minMax[1])

			// We don't even accept [6-7] since we prefer the shorter 6,7 variant; for a range to be in
			// use the hyphen needs to replace at least one digit.
			if max-min < 2 {
				panic(fmt.Sprintf("The first number in a range should be two or more digits lower than the second. Culprit possibleLength string: %s", possibleLengthString))
			}

			for j := min; j <= max; j++ {
				lengthSet[int32(j)] = true
			}
		} else {
			length, _ := strconv.Atoi(lengthSubstring)
			lengthSet[int32(length)] = true
		}
	}
	return lengthSet
}

/**
 * Reads the possible lengths present in the metadata and splits them into two sets: one for
 * full-length numbers, one for local numbers.
 *
 * @param data  one or more phone number descriptions, represented as XML nodes
 * @param lengths  a set to which to add possible lengths of full phone numbers
 * @param localOnlyLengths  a set to which to add possible lengths of phone numbers only diallable
 *     locally (e.g. within a province)
 */
func populatePossibleLengthSets(data []*PhoneNumberDescE, lengths map[int32]bool, localOnlyLengths map[int32]bool) {
	for i := 0; i < len(data); i++ {
		desc := data[i]
		if desc == nil || desc.PossibleLengths == nil {
			continue
		}

		element := desc.PossibleLengths
		nationalLengths := element.National

		// We don't add to the phone metadata yet, since we want to sort length elements found under
		// different nodes first, make sure there are no duplicates between them and that the
		// localOnly lengths don't overlap with the others.
		thisElementLengths := parsePossibleLengthStringToSet(nationalLengths)
		if element.LocalOnly != "" {
			thisElementLocalOnlyLengths := parsePossibleLengthStringToSet(element.LocalOnly)

			// intersect our two maps
			intersection := make(map[int32]bool)
			for k := range thisElementLengths {
				if thisElementLocalOnlyLengths[k] {
					intersection[k] = true
				}
			}

			if len(intersection) != 0 {
				panic(fmt.Sprintf("Possible length(s) found specified as a normal and local-only length: %v", intersection))
			}

			// We check again when we set these lengths on the metadata itself in setPossibleLengths
			// that the elements in localOnly are not also in lengths. For e.g. the generalDesc, it
			// might have a local-only length for one type that is a normal length for another type. We
			// don't consider this an error, but we do want to remove the local-only lengths.
			for k := range thisElementLocalOnlyLengths {
				localOnlyLengths[k] = true
			}
		}
		// It is okay if at this time we have duplicates, because the same length might be possible
		// for e.g. fixed-line and for mobile numbers, and this method operates potentially on
		// multiple phoneNumberDesc XML elements.
		for k := range thisElementLengths {
			lengths[k] = true
		}
	}
}

/**
 * Processes a phone number description element from the XML file and returns it as a
 * PhoneNumberDesc. If the description element is a fixed line or mobile number, the parent
 * description will be used to fill in the whole element if necessary, or any components that are
 * missing. For all other types, the parent description will only be used to fill in missing
 * components if the type has a partial definition. For example, if no "tollFree" element exists,
 * we assume there are no toll free numbers for that locale, and return a phone number description
 * with "NA" for both the national and possible number patterns. Note that the parent description
 * must therefore already be processed before this method is called on any child elements.
 *
 * @param parentDesc  a generic phone number description that will be used to fill in missing
 *     parts of the description, or null if this is the root node. This must be processed before
 *     this is run on any child elements.
 * @param countryElement  the XML element representing all the country information
 * @param numberType  the name of the number type, corresponding to the appropriate tag in the XML
 *     file with information about that type
 * @return  complete description of that phone number type
 */
// @VisibleForTesting
func processPhoneNumberDescElement(parentDesc *PhoneNumberDesc, element *PhoneNumberDescE) *PhoneNumberDesc {
	numberDesc := PhoneNumberDesc{}
	if element == nil {
		numberDesc.NationalNumberPattern = sp("NA")
		return &numberDesc
	}
	if parentDesc != nil {
		// New way of handling possible number lengths. We don't do this for the general
		// description, since these tags won't be present; instead we will calculate its values
		// based on the values for all the other number type descriptions (see
		// setPossibleLengthsGeneralDesc).
		lengths := make(map[int32]bool)
		localOnlyLengths := make(map[int32]bool)
		populatePossibleLengthSets([]*PhoneNumberDescE{element}, lengths, localOnlyLengths)
		setPossibleLengths(lengths, localOnlyLengths, parentDesc, &numberDesc)
	}

	validPattern := element.NationalNumberPattern
	if validPattern != "" {
		numberDesc.NationalNumberPattern = sp(validateRE(validPattern, true))
	}

	exampleNumber := element.ExampleNumber
	if exampleNumber != "" {
		numberDesc.ExampleNumber = sp(exampleNumber)
	}

	return &numberDesc
}

/**
 * Sets the possible length fields in the metadata from the sets of data passed in. Checks that
 * the length is covered by the "parent" phone number description element if one is present, and
 * if the lengths are exactly the same as this, they are not filled in for efficiency reasons.
 *
 * @param parentDesc  the "general description" element or null if desc is the generalDesc itself
 * @param desc  the PhoneNumberDesc object that we are going to set lengths for
 */
func setPossibleLengths(lengths map[int32]bool, localOnlyLengths map[int32]bool, parentDesc *PhoneNumberDesc, desc *PhoneNumberDesc) {
	// We clear these fields since the metadata tends to inherit from the parent element for other
	// fields (via a mergeFrom).
	desc.PossibleLength = nil
	desc.PossibleLengthLocalOnly = nil

	// Only add the lengths to this sub-type if they aren't exactly the same as the possible
	// lengths in the general desc (for metadata size reasons).
	if parentDesc == nil || !arePossibleLengthsEqual(lengths, parentDesc) {
		for length := range lengths {
			if parentDesc == nil || parentDesc.hasPossibleLength(length) {
				desc.PossibleLength = append(desc.PossibleLength, length)
			} else {
				// We shouldn't have possible lengths defined in a child element that are not covered by
				// the general description. We check this here even though the general description is
				// derived from child elements because it is only derived from a subset, and we need to
				// ensure *all* child elements have a valid possible length.
				panic(fmt.Sprintf("Out-of-range possible length found (%d), parent lengths %v.", length, parentDesc.PossibleLength))
			}
		}
	}
	// We check that the local-only length isn't also a normal possible length (only relevant for
	// the general-desc, since within elements such as fixed-line we would throw an exception if we
	// saw this) before adding it to the collection of possible local-only lengths.
	for length := range localOnlyLengths {
		if !lengths[length] {
			// We check it is covered by either of the possible length sets of the parent
			// PhoneNumberDesc, because for example 7 might be a valid localOnly length for mobile, but
			// a valid national length for fixedLine, so the generalDesc would have the 7 removed from
			// localOnly.
			if parentDesc == nil || parentDesc.hasPossibleLength(length) || parentDesc.hasPossibleLengthLocalOnly(length) {
				desc.PossibleLengthLocalOnly = append(desc.PossibleLengthLocalOnly, length)
			} else {
				panic(fmt.Sprintf("Out-of-range local-only possible length found (%d), parent length %v.", length, parentDesc.PossibleLengthLocalOnly))
			}
		}
	}

	// Need to sort both lists, possible lengths need to be ordered
	sort.Slice(desc.PossibleLength, func(i, j int) bool { return desc.PossibleLength[i] < desc.PossibleLength[j] })
	sort.Slice(desc.PossibleLengthLocalOnly, func(i, j int) bool { return desc.PossibleLengthLocalOnly[i] < desc.PossibleLengthLocalOnly[j] })
}

/**
 * Sets possible lengths in the general description, derived from certain child elements.
 */
func setPossibleLengthsGeneralDesc(generalDesc *PhoneNumberDesc, metadataId string, data *TerritoryE, isShortNumberMetadata bool) {
	lengths := make(map[int32]bool)
	localOnlyLengths := make(map[int32]bool)

	// The general description node should *always* be present if metadata for other types is
	// present, aside from in some unit tests.
	// (However, for e.g. formatting metadata in PhoneNumberAlternateFormats, no PhoneNumberDesc
	// elements are present).
	generalDescNode := data.GeneralDesc
	populatePossibleLengthSets([]*PhoneNumberDescE{generalDescNode}, lengths, localOnlyLengths)

	if len(lengths) != 0 || len(localOnlyLengths) != 0 {
		// We shouldn't have anything specified at the "general desc" level: we are going to
		// calculate this ourselves from child elements.
		panic(fmt.Sprintf("Found possible lengths specified at general desc: this should be derived from child elements. Affected country: %s", metadataId))
	}

	if !isShortNumberMetadata {
		// Make a copy here since we want to remove some nodes, but we don't want to do that on our actual data.
		// We remove no-international dialing
		trimmedDescs := []*PhoneNumberDescE{data.GeneralDesc, data.FixedLine, data.Mobile, data.Pager,
			data.TollFree, data.PremiumRate, data.SharedCost, data.PersonalNumber, data.VOIP, data.UAN, data.VoiceMail, data.StandardRate,
			data.ShortCode, data.Emergency, data.CarrierSpecific}
		populatePossibleLengthSets(trimmedDescs, lengths, localOnlyLengths)
	} else {
		populatePossibleLengthSets([]*PhoneNumberDescE{data.ShortCode}, lengths, localOnlyLengths)
		if len(localOnlyLengths) > 0 {
			panic(fmt.Errorf("found local-only lengths in short-number metadata"))
		}
	}
	setPossibleLengths(lengths, localOnlyLengths, nil, generalDesc)
}

func loadCountryMetadata(regionCode string, element *TerritoryE, isShortNumberMetadata bool, isAlternateFormatsMetadata bool) *PhoneMetadata {
	nationalPrefix := element.NationalPrefix
	metadata := loadTerritoryTagMetadata(regionCode, element, nationalPrefix)
	nationalPrefixFormattingRule := getNationalPrefixFormattingRule(element.NationalPrefixFormattingRule, nationalPrefix)
	loadAvailableFormats(metadata, element, nationalPrefix, nationalPrefixFormattingRule, element.NationalPrefixOptionalWhenFormatting)

	if !isAlternateFormatsMetadata {
		// The alternate formats metadata does not need most of the patterns to be set.
		setRelevantDescPatterns(metadata, element, isShortNumberMetadata)
	}
	return metadata
}

func setRelevantDescPatterns(metadata *PhoneMetadata, element *TerritoryE, isShortNumberMetadata bool) {
	generalDesc := processPhoneNumberDescElement(nil, element.GeneralDesc)

	// Calculate the possible lengths for the general description. This will be based on the
	// possible lengths of the child elements.
	setPossibleLengthsGeneralDesc(generalDesc, metadata.GetId(), element, isShortNumberMetadata)
	metadata.GeneralDesc = generalDesc

	if !isShortNumberMetadata {
		// Set fields used by regular length phone numbers.
		metadata.FixedLine = processPhoneNumberDescElement(generalDesc, element.FixedLine)
		metadata.Mobile = processPhoneNumberDescElement(generalDesc, element.Mobile)
		metadata.SharedCost = processPhoneNumberDescElement(generalDesc, element.SharedCost)
		metadata.Voip = processPhoneNumberDescElement(generalDesc, element.VOIP)
		metadata.PersonalNumber = processPhoneNumberDescElement(generalDesc, element.PersonalNumber)
		metadata.Pager = processPhoneNumberDescElement(generalDesc, element.Pager)
		metadata.Uan = processPhoneNumberDescElement(generalDesc, element.UAN)
		metadata.Voicemail = processPhoneNumberDescElement(generalDesc, element.VoiceMail)
		metadata.NoInternationalDialling = processPhoneNumberDescElement(generalDesc, element.NoInternationalDialing)

		mobileAndFixedAreSame := *metadata.Mobile.NationalNumberPattern == *metadata.FixedLine.NationalNumberPattern
		if metadata.GetSameMobileAndFixedLinePattern() != mobileAndFixedAreSame {
			metadata.SameMobileAndFixedLinePattern = bp(mobileAndFixedAreSame)
		}

		metadata.TollFree = processPhoneNumberDescElement(generalDesc, element.TollFree)
		metadata.PremiumRate = processPhoneNumberDescElement(generalDesc, element.PremiumRate)
	} else {
		// Set fields used by short numbers.
		metadata.StandardRate = processPhoneNumberDescElement(generalDesc, element.StandardRate)
		metadata.ShortCode = processPhoneNumberDescElement(generalDesc, element.ShortCode)
		metadata.CarrierSpecific = processPhoneNumberDescElement(generalDesc, element.CarrierSpecific)
		metadata.Emergency = processPhoneNumberDescElement(generalDesc, element.Emergency)
		metadata.TollFree = processPhoneNumberDescElement(generalDesc, element.TollFree)
		metadata.PremiumRate = processPhoneNumberDescElement(generalDesc, element.PremiumRate)
	}
}

// <!ELEMENT phoneNumberMetadata (territories)>
type PhoneNumberMetadataE struct {
	// <!ELEMENT territories (territory+)>
	Territories []TerritoryE `xml:"territories>territory"`
}

// <!ELEMENT territory (references?, availableFormats?, generalDesc, noInternationalDialling?,
//        fixedLine?, mobile?, pager?, tollFree?, premiumRate?,
//        sharedCost?, personalNumber?, voip?, uan?, voicemail?)>
type TerritoryE struct {
	// <!ATTLIST territory id CDATA #REQUIRED>
	ID string `xml:"id,attr"`

	// <!ATTLIST territory mainCountryForCode (true) #IMPLIED>
	MainCountryForCode bool `xml:"mainCountryForCode,attr"`

	// <!ATTLIST territory leadingDigits CDATA #IMPLIED>
	LeadingDigits string `xml:"leadingDigits,attr"`

	// <!ATTLIST territory countryCode CDATA #REQUIRED>
	CountryCode int32 `xml:"countryCode,attr"`

	// <!ATTLIST territory nationalPrefix CDATA #IMPLIED>
	NationalPrefix string `xml:"nationalPrefix,attr"`

	// <!ATTLIST territory internationalPrefix CDATA #IMPLIED>
	InternationalPrefix string `xml:"internationalPrefix,attr"`

	// <!ATTLIST territory preferredInternationalPrefix CDATA #IMPLIED>
	PreferredInternationalPrefix string `xml:"preferredInternationalPrefix,attr"`

	// <!ATTLIST territory nationalPrefixFormattingRule CDATA #IMPLIED>
	NationalPrefixFormattingRule string `xml:"nationalPrefixFormattingRule,attr"`

	// <!ATTLIST territory mobileNumberPortableRegion (true) #IMPLIED>
	MobileNumberPortableRegion bool `xml:"mobileNumberPortableRegion,attr"`

	// <!ATTLIST territory nationalPrefixForParsing CDATA #IMPLIED>
	NationalPrefixForParsing string `xml:"nationalPrefixForParsing,attr"`

	// <!ATTLIST territory nationalPrefixTransformRule CDATA #IMPLIED>
	NationalPrefixTransformRule string `xml:"nationalPrefixTransformRule,attr"`

	// <!ATTLIST territory preferredExtnPrefix CDATA #IMPLIED>
	PreferredExtnPrefix string `xml:"PreferredExtnPrefix"`

	// <!ATTLIST territory nationalPrefixOptionalWhenFormatting (true) #IMPLIED>
	NationalPrefixOptionalWhenFormatting bool `xml:"nationalPrefixOptionalWhenFormatting,attr"`

	// <!ATTLIST territory carrierCodeFormattingRule CDATA #IMPLIED>
	CarrierCodeFormattingRule string `xml:"carrierCodeFormattingRule,attr"`

	// <!ELEMENT references (sourceUrl+)>
	// <!ELEMENT sourceUrl (#PCDATA)>
	References []string `xml:"references>sourceUrl"`

	// <!ELEMENT availableFormats (numberFormat+)>
	AvailableFormats []NumberFormatE `xml:"availableFormats>numberFormat"`

	// <!ELEMENT generalDesc (nationalNumberPattern)>
	GeneralDesc *PhoneNumberDescE `xml:"generalDesc"`

	// <!ELEMENT noInternationalDialling (nationalNumberPattern, possibleLengths, exampleNumber)>
	NoInternationalDialing *PhoneNumberDescE `xml:"noInternationalDialing"`

	// <!ELEMENT fixedLine (nationalNumberPattern, possibleLengths, exampleNumber)>
	FixedLine *PhoneNumberDescE `xml:"fixedLine"`

	// <!ELEMENT mobile (nationalNumberPattern, possibleLengths, exampleNumber)>
	Mobile *PhoneNumberDescE `xml:"mobile"`

	// <!ELEMENT pager (nationalNumberPattern, possibleLengths, exampleNumber)>
	Pager *PhoneNumberDescE `xml:"pager"`

	// <!ELEMENT tollFree (nationalNumberPattern, possibleLengths, exampleNumber)>
	TollFree *PhoneNumberDescE `xml:"tollFree"`

	// <!ELEMENT premiumRate (nationalNumberPattern, possibleLengths, exampleNumber)>
	PremiumRate *PhoneNumberDescE `xml:"premiumRate"`

	// <!ELEMENT sharedCost (nationalNumberPattern, possibleLengths, exampleNumber)>
	SharedCost *PhoneNumberDescE `xml:"sharedCost"`

	// <!ELEMENT personalNumber (nationalNumberPattern, possibleLengths, exampleNumber)>
	PersonalNumber *PhoneNumberDescE `xml:"personalNumber"`

	// <!ELEMENT voip (nationalNumberPattern, possibleLengths, exampleNumber)>
	VOIP *PhoneNumberDescE `xml:"voip"`

	// <!ELEMENT uan (nationalNumberPattern, possibleLengths, exampleNumber)>
	UAN *PhoneNumberDescE `xml:"uan"`

	// <!ELEMENT voicemail (nationalNumberPattern, possibleLengths, exampleNumber)>
	VoiceMail *PhoneNumberDescE `xml:"voicemail"`

	// <!ELEMENT uan (nationalNumberPattern, possibleLengths, exampleNumber)>
	StandardRate *PhoneNumberDescE `xml:"standardRate"`

	// <!ELEMENT voicemail (nationalNumberPattern, possibleLengths, exampleNumber)>
	ShortCode *PhoneNumberDescE `xml:"shortCode"`

	// <!ELEMENT uan (nationalNumberPattern, possibleLengths, exampleNumber)>
	Emergency *PhoneNumberDescE `xml:"Emergency"`

	// <!ELEMENT voicemail (nationalNumberPattern, possibleLengths, exampleNumber)>
	CarrierSpecific *PhoneNumberDescE `xml:"carrierSpecific"`
}

// <!ELEMENT numberFormat (leadingDigits*, format, intlFormat*)>
type NumberFormatE struct {
	// <!ELEMENT leadingDigits (#PCDATA)>
	LeadingDigits []string `xml:"leadingDigits"`

	// <!ELEMENT format (#PCDATA)>
	Format string `xml:"format"`

	// <!ELEMENT intlFormat (#PCDATA)>
	InternationalFormat []string `xml:"intlFormat"`

	// <!ATTLIST numberFormat nationalPrefixFormattingRule CDATA #IMPLIED>
	NationalPrefixFormattingRule string `xml:"nationalPrefixFormattingRule,attr"`

	// <!ATTLIST numberFormat nationalPrefixOptionalWhenFormatting (true) #IMPLIED>
	NationalPrefixOptionalWhenFormatting *bool `xml:"nationalPrefixOptionalWhenFormatting,attr"`

	// <!ATTLIST numberFormat carrierCodeFormattingRule CDATA #IMPLIED>
	CarrierCodeFormattingRule string `xml:"carrierCodeFormattingRule,attr"`

	// <!ATTLIST numberFormat pattern CDATA #REQUIRED>
	Pattern string `xml:"pattern,attr" validate:"required"`
}

type PossibleLengthE struct {
	// <!ATTLIST possibleLengths national CDATA #REQUIRED>
	National string `xml:"national,attr"`

	// <!ATTLIST possibleLengths localOnly CDATA #IMPLIED>
	LocalOnly string `xml:"localOnly,attr"`
}

type PhoneNumberDescE struct {
	// <!ELEMENT nationalNumberPattern (#PCDATA)>
	NationalNumberPattern string `xml:"nationalNumberPattern"`

	// <!ELEMENT possibleLengths EMPTY>
	PossibleLengths *PossibleLengthE `xml:"possibleLengths"`

	// <!ELEMENT exampleNumber (#PCDATA)>
	ExampleNumber string `xml:"exampleNumber"`
}
