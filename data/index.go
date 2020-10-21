package data

// Index contains a map of collection names to their terms.
var Index = map[string][]string{
	"IACO Prefixes":                  IacoPrefixes,
	"Symbols of Chemical Elements":   ChemicalElements,
	"ISO-3166 Alpha-2 Country Codes": Iso3166Alpha2s,
	"ISO-3166 Alpha-3 Country Codes": Iso3166Alpha3s,
	"National Rail CRS Codes":        CrsCodes,
	"US State Abbreviations":         StateAbbreviations,
	"Vowels":                         Vowels,
	"Consonants":                     Consonants,
	"Letters with ascenders":         LettersWithAscenders,
	"Letters without ascenders":      LettersWithoutAscenders,
	"Letters with descenders":        LettersWithDescenders,
	"Letters without descenders":     LettersWithoutDescenders,
}
