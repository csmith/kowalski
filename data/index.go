package data

// Index contains a map of collection names to their terms.
var Index = map[string][]string{
	"IACO prefixes":                                        IacoPrefixes,
	"symbols of chemical elements":                         ChemicalElements,
	"ISO-3166 Alpha-2 Country Codes":                       Iso3166Alpha2s,
	"ISO-3166 Alpha-3 Country Codes":                       Iso3166Alpha3s,
	"National Rail CRS Codes":                              CrsCodes,
	"US State Abbreviations":                               StateAbbreviations,
	"vowels":                                               Vowels,
	"consonants":                                           Consonants,
	"letters with ascenders":                               LettersWithAscenders,
	"letters without ascenders":                            LettersWithoutAscenders,
	"letters with descenders":                              LettersWithDescenders,
	"letters without descenders":                           LettersWithoutDescenders,
	"letters from the upper row of a QWERTY keyboard":      QwertyUpperRow,
	"letters from the home row of a QWERTY keyboard":       QwertyHomeRow,
	"letters from the lower row of a QWERTY keyboard":      QwertyLowerRow,
	"letters from the upper two rows of a QWERTY keyboard": QwertyUpperTwoRows,
	"letters from the lower two rows of a QWERTY keyboard": QwertyLowerTwoRows,
	"letters from the outer two rows of a QWERTY keyboard": QwertyOuterRows,
}
