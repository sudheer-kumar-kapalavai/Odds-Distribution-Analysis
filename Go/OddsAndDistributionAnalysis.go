package main

import (
	"fmt"
    "io"
    "net/http"
    "os"
    "time"
    "encoding/csv"
    "log"
    "strconv"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
)

func main() {

	var url,filename, stockName, startDate, endDate, interval string
	var start, end int64

	filename = "data.csv"

	
	fmt.Println("Enter the Stock Name as per Yahoo finance: ")
    fmt.Scanln(&stockName)

    fmt.Println("Enter the Start date in YYYY-MM-DD: ")
    fmt.Scanln(&startDate)
    start = getEpoc(startDate)

    fmt.Println("Enter the End date in YYYY-MM-DD: ")
    fmt.Scanln(&endDate)
    end = getEpoc(endDate)

    fmt.Println("Select interval from [1d, 1wk, 1mo]: ")
    fmt.Scanln(&interval)



	url = "https://query1.finance.yahoo.com/v7/finance/download/" + stockName + "?period1=" + strconv.FormatInt(start,10) + "&period2=" + strconv.FormatInt(end,10) + "&interval=" + interval + "&events=history"
    err := DownloadFile(url, filename)
    if err != nil {
        panic(err)
    }



	var data[][]string = readSample(filename)

	returns := make([]float64, len(data) - 1)

	for i:= 1; i < len(data); i++{
		close, err := strconv.ParseFloat(data[i][5], 64)
		if err != nil {
            continue
        }
		open,  err := strconv.ParseFloat(data[i][1], 64)
		if err != nil {
            continue
        }
		returns[i-1] = (close - open)/open
	}
	histPlot(returns, "histogram of " + stockName)

}


func getEpoc(date string) int64 {
	//Convert the given date to epoc value

	thetime, e := time.Parse(time.RFC3339, date + "T00:00:00+00:00")
	if e != nil {
		panic("Can't parse time format")
	}
	epoch := thetime.Unix()
	return epoch
}


func readSample(filename string) [][]string {
	//Read the given filename

    f, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    rows, err := csv.NewReader(f).ReadAll()
    f.Close()
    if err != nil {
        log.Fatal(err)
    }
    return rows
}

func histPlot(values plotter.Values, name string) {
	//Draw the plot and store it in current folder

    p := plot.New()
    p.Title.Text = name
 
    hist, err := plotter.NewHist(values, 20)
    if err != nil {
        panic(err)
    }
    p.Add(hist)
 
    if err := p.Save(3*vg.Inch, 3*vg.Inch, "hist.png"); err != nil {
        panic(err)
    }
}

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func DownloadFile(url string, filepath string) error {
    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}