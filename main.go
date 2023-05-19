package main

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Constants
const (
	bucketURL = "https://fil-archive.s3.us-east-2.amazonaws.com/mainnet/csv/1/"
)

// Function to download and decompress a gzipped CSV file
func downloadAndDecompressFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer reader.Close()

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

// Function to read CSV and count deals
func countDeals(filePath string) (int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	count := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		// Do your counting logic based on the record
		count++
	}

	return count, nil
}

// Main function
func main() {
	// Calculate dates for last 7 days
	now := time.Now()
	for d := 6; d >= 0; d-- {
		date := now.AddDate(0, 0, -d)
		dateStr := date.Format("2006-01-02")  // The Go time package uses this specific date as a format string

		// Construct file URLs and local file paths
		for _, table := range []string{"market_deal_proposals", "market_deal_states"} {
			url := fmt.Sprintf("%s%s/%s/%s-%s.csv.gz", bucketURL, table, date.Year(), table, dateStr)
			localFile := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s.csv", table, dateStr))

			// Download and decompress the file
			err := downloadAndDecompressFile(url, localFile)
			if err != nil {
				fmt.Printf("Error downloading file %s: %v\n", url, err)
				continue
			}

			// Count deals
			count, err := countDeals(localFile)
			if err != nil {
				fmt.Printf("Error counting deals in file %s: %v\n", localFile, err)
				continue
			}

			fmt.Printf("Counted %d deals in %s for %s\n", count, table, dateStr)
		}
	}
}
