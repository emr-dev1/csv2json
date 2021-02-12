package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	var file string
	var output string
	var headers string

	app := cli.NewApp()

	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "file",
			Aliases:     []string{"f"},
			Usage:       "The path to the csv file to convert to JSON",
			Destination: &file,
		},
		&cli.StringFlag{
			Name:        "headers",
			Usage:       "The headers to include in the parsing of the csv file",
			Destination: &headers,
		},
		&cli.StringFlag{
			Name:        "output",
			Aliases:     []string{"o"},
			Usage:       "The path to output the converted JSON file to.",
			Destination: &output,
		},
	}

	app.Action = func(c *cli.Context) error {
		// Check if the filepath for the input csv file was passed without flags
		if c.NArg() > 0 {
			file = c.Args().Get(0)
			fmt.Println("Using file without passing the -f flag: ", file)
		}

		// Get the file info for the input csv file
		fileInfo, err := os.Stat(file)
		if err != nil {
			return errors.New("File does not exist.")
		}

		// If the file exists then get the name for use later if the output
		// path does not include a json file name.
		var inFileName string = fileInfo.Name()
		var outFileName string
		s := strings.Split(output, "/")

		// If the output path includes a json file to write to then capture the name
		// of the file and pop it off the end of the string to get a valid output path
		// that can be checked later.
		if strings.HasSuffix(output, ".json") {
			outFileName = s[len(s)-1]
			s = s[:len(s)-1]
			output = strings.Join(s, "/")
		}

		// Check if the output path is valid.
		fileInfo, err = os.Stat(output)
		if err != nil {
			return errors.New("Invalid output directory")
		}
		outFileName = strings.Replace(inFileName, "csv", "json", 1)

		outFile, err := os.Create(output + "/" + outFileName)
		if err != nil {
			fmt.Println("Failed to create new file")
		}

		jsonObj, err := readCsvFile(file)
		if err != nil {
			return errors.New("Failed to parse csv file")
		}

		outJsonObj, _ := json.Marshal(jsonObj)

		_, err = outFile.Write(outJsonObj)
		if err != nil {
			return errors.New("Failed to write to the out file.")
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// readCsvFile takes in the path of a csv file, opens and reads the contents of the file.
// Returns a map matching the structure of what will be written to the JSON file.
func readCsvFile(filePath string) ([]map[string]string, error) {
	ret := []map[string]string{}

	// Open the csv file to parse
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("Failed to open the file")
	}

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	headers := records[0]

	for _, row := range records[1:] {
		jsonObj := map[string]string{}
		for j, col := range row {
			jsonObj[headers[j]] = string(col)
		}
		ret = append(ret, jsonObj)
	}

	return ret, nil
}
