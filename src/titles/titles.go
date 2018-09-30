package titles

import "github.com/stock-simulator-server/src/money"

var Titles = make(map[int64]*Title)

type Title struct {
	Level int64  `json:"level"`
	Name  string `json:"name"`
	Cost  int64  `json:"cost"`
}

func makeTitle(level, cost int64, name string) {
	Titles[level] = &Title{
		Level: level,
		Name:  name,
		Cost:  cost,
	}
}

func PopulateTitles() {
	makeTitle(0, 0, "Noob")
	makeTitle(1, 2*money.Thousand, "l33t N00b")
	makeTitle(2, 10*money.Thousand, "Im trying")
	makeTitle(3, 50*money.Thousand, "Road to T500")
	makeTitle(3, 100*money.Thousand, "Wall Street Bro")
}
