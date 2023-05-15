package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Deal represents a single deal in the Filecoin network.
type Deal struct {
	Proposal struct {
		Provider string `json:"Provider"`
		// Include other fields from the deal here...
	} `json:"Proposal"`
}

func main() {
	// Read the JSON file
	data, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Printf("Error reading JSON file: %v\n", err)
		os.Exit(1)
	}

	// Unmarshal the JSON data into a slice of Deal objects
	var deals []Deal
	err = json.Unmarshal(data, &deals)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON data: %v\n", err)
		os.Exit(1)
	}

	// Now you can perform your analysis on the deals
	// For example, to count the number of deals per provider:
	providerDealCounts := make(map[string]int)
	for _, deal := range deals {
		providerDealCounts[deal.Proposal.Provider]++
	}

	fmt.Printf("Number of deals per provider:\n")
	for provider, count := range providerDealCounts {
		fmt.Printf("Provider %s: %d deals\n", provider, count)
	}
}
