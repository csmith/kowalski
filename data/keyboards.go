package data

// QwertyUpperRow contains the characters in the upper row of a QWERTY keyboard.
var QwertyUpperRow = []string{
	"q",
	"w",
	"e",
	"r",
	"t",
	"y",
	"u",
	"i",
	"o",
	"p",
}

// QwertyHomeRow contains the characters in the home row of a QWERTY keyboard.
var QwertyHomeRow = []string{
	"a",
	"s",
	"d",
	"f",
	"g",
	"h",
	"j",
	"k",
	"l",
}

// QwertyLowerRow contains the characters in the lower row of a QWERTY keyboard.
var QwertyLowerRow = []string{
	"z",
	"x",
	"c",
	"v",
	"b",
	"n",
	"m",
}

// QwertyUpperTwoRows contains the characters in the top two rows of a QWERTY keyboard.
var QwertyUpperTwoRows = append(QwertyUpperRow, QwertyHomeRow...)
// QwertyLowerTwoRows contains the characters in the lower two rows of a QWERTY keyboard.
var QwertyLowerTwoRows = append(QwertyLowerRow, QwertyHomeRow...)
// QwertyOuterRows contains the characters in the outer two rows of a QWERTY keyboard.
var QwertyOuterRows = append(QwertyLowerRow, QwertyUpperRow...)
