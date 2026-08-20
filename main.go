package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

var env = "development"
var version = "0.0.0"
var help = fmt.Sprintf(Help, version)

// Args stores cli arguments
type Args struct {
	Cumulative  bool
	Summary     bool
	SummaryJSON bool
	Quiet       bool
	Delimiter   string
	Help        bool
	Version     bool
}

// CheckSummary stores verification summary data
type CheckSummary struct {
	IsValid      bool `json:"is_valid"`
	ErrosCount   int  `json:"errors"`
	TotalRows    int  `json:"rows"`
	TotalColumns int  `json:"columns"`
}

// String converts CheckSumary into a human-readable string
func (c CheckSummary) String() string {
	status := "VALID"
	if !c.IsValid {
		status = "INVALID"
	}
	return fmt.Sprintf("[%v]\nErrors: %v\nTotal columns: %v\nTotal rows: %v", status, c.ErrosCount, c.TotalColumns, c.TotalRows)
}

// JSON converts CheckSumary into a JSON string
func (c CheckSummary) JSON() string {
	j, _ := json.Marshal(c)
	return string(j)
}

func main() {
	var args = Args{}
	flag.BoolVar(&args.Cumulative, "c", false, "Collect all csv errors and output the list at the end")
	flag.BoolVar(&args.Quiet, "q", false, "Silently terminate with exit(1) upon the first error encountered in the csv")
	flag.BoolVar(&args.Summary, "s", false, "Print summary info (check summary, total columns, total rows)")
	flag.BoolVar(&args.SummaryJSON, "sj", false, "Print summary infoas a JSON string")
	flag.StringVar(&args.Delimiter, "d", ",", "Fields separator (default: comma)")
	flag.BoolVar(&args.Help, "h", false, "Help")
	flag.BoolVar(&args.Help, "help", false, "Alias for -h")
	flag.BoolVar(&args.Version, "v", false, "Version")
	flag.BoolVar(&args.Version, "version", false, "Alias for -v")
	flag.Usage = func() {}
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("checkcsv: ")

	if args.Help {
		fmt.Println(help)
		os.Exit(0)
	}

	if args.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	if args.Summary && args.SummaryJSON {
		log.Fatalln("You cannot use both output formats (-s and -sj) at the same time")
	}

	var errors_list strings.Builder

	var summary = CheckSummary{
		IsValid:      true,
		ErrosCount:   0,
		TotalRows:    0,
		TotalColumns: 0,
	}

	var reader = csv.NewReader(os.Stdin)
	delimiter, _ := utf8.DecodeRuneInString(args.Delimiter)
	reader.Comma = delimiter
	reader.ReuseRecord = true

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			if args.Summary || args.SummaryJSON {
				summary.IsValid = false
				summary.ErrosCount++
			} else {
				if args.Quiet {
					os.Exit(1)
				}

				if args.Cumulative {
					errors_list.WriteString(err.Error() + "\n")
				} else {
					fmt.Print(err.Error())
					os.Exit(1)
				}
			}
		}

		if args.Summary || args.SummaryJSON {
			if summary.TotalColumns == 0 {
				summary.TotalColumns = len(row)
			}
			summary.TotalRows++
		}
	}

	if args.Summary {
		fmt.Println(summary.String())
		os.Exit(0)
	}

	if args.SummaryJSON {
		fmt.Println(summary.JSON())
		os.Exit(0)
	}

	if !args.Quiet && args.Cumulative && errors_list.Len() > 0 {
		fmt.Print(errors_list.String())
		os.Exit(1)
	}
}
