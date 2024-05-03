// Package products contains details relating to the loading of products
package products

type ProductSubscriptions map[ProductName][]Instrument

type ProductName string

type Instrument struct {
	Market     string
	Instrument string
}

const CADLITick ProductName = "CADLITick"

var exampleProducts = ProductSubscriptions{
	CADLITick: []Instrument{
		{
			Market:     "cadli",
			Instrument: "BTC-USD",
		},
	},
}

// Load returns the required product subscriptions
func Load() ProductSubscriptions {
	return exampleProducts
}
