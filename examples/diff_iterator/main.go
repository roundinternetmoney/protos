// systhesize an order book and iterate over each side
package main

import (
	"fmt"

	pb "roundinternet.money/protos/gen/dex/base/v1"
)

func main() {
	diff := pb.NewBookDiff()

	diff.Asks = append(diff.Asks, pb.NewBookLevel("106", "3"))
	diff.Asks = append(diff.Asks, pb.NewBookLevel("105", "2"))
	diff.Asks = append(diff.Asks, pb.NewBookLevel("104", "1"))
	diff.Bids = append(diff.Bids, pb.NewBookLevel("103", "1"))
	diff.Bids = append(diff.Bids, pb.NewBookLevel("102", "2"))
	diff.Bids = append(diff.Bids, pb.NewBookLevel("101", "3"))

	for p, s := range diff.BidLevels() {
		fmt.Printf("BID | Price: %s Size: %s\n", p, s)
	}
	for p, s := range diff.AskLevels() {
		fmt.Printf("ASK | Price: %s Size: %s\n", p, s)
	}
}
