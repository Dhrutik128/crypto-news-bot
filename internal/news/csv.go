package news

import (
	"encoding/csv"
	"fmt"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"os"
)

func readCsv(fileName string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(fileName)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}


func writeFeed(fileName string, feed *gofeed.Feed) error {
	// Open CSV file
	f, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", feed.FeedLink)); err != nil {
		log.Println(err)
		return err
	}
	return nil
}