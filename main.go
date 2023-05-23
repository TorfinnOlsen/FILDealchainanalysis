package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Deal struct {
	Proposal struct {
		StartEpoch int64 `json:"StartEpoch"`
	} `json:"Proposal"`
}

func main() {
	jsonFile, err := os.Open("deals.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened deals.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]Deal
	json.Unmarshal([]byte(byteValue), &result)

	var wg sync.WaitGroup
	var dealCount int
	var recentDealEpoch int64
	var weekDealCount int

	currentTime := time.Now()
	weekAgo := currentTime.AddDate(0, 0, -7)

	for _, deal := range result {
		wg.Add(1)
		go func(deal Deal) {
			defer wg.Done()

			dealTime := time.Unix(deal.Proposal.StartEpoch, 0)

			// Check if deal is in the past 7 days
			if dealTime.After(weekAgo) {
				weekDealCount++
			}

			// Check if this deal is the most recent
			if deal.Proposal.StartEpoch > recentDealEpoch {
				recentDealEpoch = deal.Proposal.StartEpoch
			}

			dealCount++
		}(deal)
	}

	wg.Wait()

	fmt.Println("Total number of deals:", dealCount)
	fmt.Println("Number of deals in the past 7 days:", weekDealCount)
	fmt.Println("Date of the most recent deal:", time.Unix(recentDealEpoch, 0))
}
