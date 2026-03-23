// compares 2 assets, 2 identical lists of assets, and 2 differing length asset lists
package main

import (
	"fmt"

	basev1 "roundinternet.money/protos/gen/dex/base/v1"
)

// compare succesfuly, then change the lengths of the lists
func main() {
	one := &basev1.Asset{
		Id:          "BTC-USD",
		Ticker:      "BTC",
		MinQuantity: 0.001,
		MaxQuantity: 100,
		LotSize:     0.001,
		TickSize:    0.1,
	}

	two := &basev1.Asset{
		Id:          "BTC-USD",
		Ticker:      "BTC",
		MinQuantity: 0.001,
		MaxQuantity: 100,
		LotSize:     0.001,
		TickSize:    0.1,
	}

	fmt.Printf("Assets equal: %v\n", one.Eq(two))

	three := &basev1.DexAssetResponse{
		A: []*basev1.Asset{one, one, one, one},
	}
	four := &basev1.DexAssetResponse{
		A: []*basev1.Asset{two, two, two, two},
	}

	fmt.Printf("Asset lists equal: %v\n", three.Eq(four))

	five := &basev1.DexAssetResponse{
		A: []*basev1.Asset{one, one, one, one},
	}
	six := &basev1.DexAssetResponse{
		A: []*basev1.Asset{two, two, two, },
	}

	fmt.Printf("Asset lists equal: %v\n", five.Eq(six))	
}
