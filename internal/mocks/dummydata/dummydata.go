package dummydata

import (
	"time"

	"github.com/modell-aachen/gologger/internal/interfaces"
)

var Source interfaces.SourceString = interfaces.SourceString("source")
var Level interfaces.LevelString = interfaces.LevelString("warning")
var Levels []interfaces.LevelString = []interfaces.LevelString{Level}

var Row0 interfaces.LogRow = interfaces.LogRow{
	time.Unix(123, 0),
	Source,
	Level,
	interfaces.LogGeneric{
		"earlykey": "earlyvalue",
	},
	interfaces.LogFields{
		"early",
		"entry",
	},
}
var Row1 interfaces.LogRow = interfaces.LogRow{
	time.Unix(554385600, 0),
	Source,
	Level,
	interfaces.LogGeneric{
		"caller": "Mr. R. A.",
		"key a":  "value a",
		"key b":  "value b",
	},
	interfaces.LogFields{
		"field a",
		"field b",
	},
}
var Row2 interfaces.LogRow = interfaces.LogRow{
	time.Unix(554385610, 0),
	Source,
	Level,
	interfaces.LogGeneric{},
	interfaces.LogFields{
		"late entry",
	},
}
