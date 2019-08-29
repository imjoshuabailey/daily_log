package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/extrame/xls"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
	var err error

	execCtx, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.DefaultExecAllocatorOptions[0:1]...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(execCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// run task authenticates user
	err = chromedp.Run(ctx, authenticate(`https://tsdloaner.tsd-inc.com`, os.Getenv(`ACCOUNT`), os.Getenv(`USERNAME`), os.Getenv(`PASS`)))
	if err != nil {
		log.Fatal(err)
	}

	// run task exports units available report
	err = chromedp.Run(ctx, exportUnitsAvailable())
	if err != nil {
		log.Fatal(err)
	}

	// run task exports unit duration report
	err = chromedp.Run(ctx, exportUnitDuration())
	if err != nil {
		log.Fatal(err)
	}

	// run task exports unit utilization report
	err = chromedp.Run(ctx, exportUnitUtilization())
	if err != nil {
		log.Fatal(err)
	}

	// run task exports current units out report
	err = chromedp.Run(ctx, exportUnitsOut())
	if err != nil {
		log.Fatal(err)
	}

	// run task exports full inventory report
	err = chromedp.Run(ctx, exportFullInventory())
	if err != nil {
		log.Fatal(err)
	}

	// run task wait 10 seconds
	err = chromedp.Run(ctx, idleChrome())
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = chromedp.Cancel(ctx)
	if err != nil {
		log.Fatal(err)
	}

	processUnitsAvailable()

	processContractsOpen()

	u, err := processUtilization("/Users/joshbailey/downloads/AdjustedUtilization.csv")
	if err != nil {
		log.Fatalf("error getting Utilization: %v", err)
	}
	os.Remove("/Users/joshbailey/downloads/AdjustedUtilization.csv")

	var _ = u

	d, err := processDuration("/Users/joshbailey/downloads/Service_Advisor_Activity.csv")
	if err != nil {
		log.Fatalf("error getting Duration: %v", err)
	}
	os.Remove("/Users/joshbailey/downloads/Service_Advisor_Activity.csv")
	var _ = d

	processFullInventory()
}

// authenticate and log in user
func authenticate(url string, account string, username string, password string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`//input[@name="txtPassword"]`),
		chromedp.SendKeys(`//input[@name="txtTsdNum"]`, account),
		chromedp.Submit(`//input[@name="txtTsdNum"]`),
		chromedp.SendKeys(`//input[@name="txtUsername"]`, username),
		chromedp.SendKeys(`//input[@name="txtPassword"]`, password),
		chromedp.Click(`//input[@name="LoginBtn"]`),
		chromedp.WaitVisible(`//form[@id="form1"]`),
	}
}

// export UnitsAvailable.xls file
func exportUnitsAvailable() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://tsdloaner.tsd-inc.com/UnitsAvailable.aspx`),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click(`//input[@name="butExport"]`),
	}
}

// export Service_Advisor_Activity.csv file
func exportUnitDuration() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://tsdloaner.tsd-inc.com/ServiceWriters.aspx`),
		chromedp.Sleep(1 * time.Second),
		chromedp.Click(`//input[@name="butRun"]`),
		chromedp.Click(`//input[@name="btnExport"]`),
		chromedp.Click(`//input[@name="butClose"]`),
	}
}

// export AdjustedUtilization.csv file
func exportUnitUtilization() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://tsdloaner.tsd-inc.com/AdjustedUnitUtil.aspx`),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click(`//input[@name="butRun"]`),
		chromedp.Click(`//input[@name="btnExport"]`),
		chromedp.Click(`//input[@name="butClose"]`),
	}
}

// export ClosedContractsOpen.xls file
func exportUnitsOut() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://tsdloaner.tsd-inc.com/CurrOpenRa.aspx`),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click(`//input[@name="butExport"]`),
		chromedp.Sleep(1 * time.Second),
	}
}

//export FullInventory.xls file
func exportFullInventory() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://tsdloaner.tsd-inc.com/FullInventory.aspx`),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click(`//input[@name="butExport"]`),
	}
}

func idleChrome() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Sleep(10 * time.Second),
	}
}

// count the number of 2019 vehicles in units available file
func processUnitsAvailable() {
	if xlFile, err := xls.Open("/Users/joshbailey/downloads/UnitsAvailable.xls", "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			year := sheet1.Row(0).Col(3)
			var count = 0
			for i := 0; i <= (int(sheet1.MaxRow)); i++ {
				row1 := sheet1.Row(i)
				year = row1.Col(3)
				if year == "2019" || year == "19" {
					count = count + 1
				} else if year == "2020" || year == "20" {
					count = count + 1
				}

			}
			fmt.Println(count, "Vehicles Available")
		}
		os.Remove("/Users/joshbailey/downloads/UnitsAvailable.xls")
	} else {
		log.Print(err)
	}
}

func processContractsOpen() {
	if xlFile, err := xls.Open("/Users/joshbailey/downloads/ClosedContractsOpen.xls", "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			days := sheet1.Row(3).Col(15)
			var greaterThanThree = 0
			var doubleDigit = 0
			for i := 0; i <= (int(sheet1.MaxRow)); i++ {
				row1 := sheet1.Row(i)
				days = row1.Col(15)
				if days, err := strconv.Atoi(days); err == nil {
					if days > 3 {
						greaterThanThree = greaterThanThree + 1
						if days > 9 {
							doubleDigit = doubleDigit + 1
						}
					}
				}

			}
			fmt.Println(greaterThanThree, "vehicles out more than 3 days")
			fmt.Println(doubleDigit, "vehicles out more than 9 days")
		}
		os.Remove("/Users/joshbailey/downloads/ClosedContractsOpen.xls")
	} else {
		log.Print(err)
	}
}

func processUtilization(path string) (float64, error) {

	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return 0, err
	}

	removeSybol := strings.Index(records[len(records)-1][7], "%")
	value, err := strconv.ParseFloat(records[len(records)-1][7][:removeSybol], 64)

	if err != nil {
		return 0, err
	}
	fmt.Println("Current utilization percentage is", value)
	return value, nil
}

func processDuration(path string) (float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(records[len(records)-1][7], 64)

	if err != nil {
		return 0, err
	}
	fmt.Println("Current average duration is", value)
	return value, nil

}

func processFullInventory() {
	if xlFile, err := xls.Open("/Users/joshbailey/downloads/FullInventory.xls", "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			year := sheet1.Row(0).Col(3)
			var count2017 = 0
			var count2018 = 0
			var count2019 = 0
			var count2020 = 0
			for i := 0; i <= (int(sheet1.MaxRow)); i++ {
				row1 := sheet1.Row(i)
				year = row1.Col(3)
				if year == "2017" || year == "17" {
					count2017 = count2017 + 1
				} else if year == "2018" || year == "18" {
					count2018 = count2018 + 1
				} else if year == "2019" || year == "19" {
					count2019 = count2019 + 1
				} else if year == "2020" || year == "20" {
					count2020 = count2020 + 1
				}
			}
			fmt.Println(count2017+count2018+count2019+count2020, "Vehicles in fleet")
			fmt.Println(count2017, "2017s")
			fmt.Println(count2018, "2018s")
			fmt.Println(count2019, "2019s")
			fmt.Println(count2020, "2020s")

		}
		os.Remove("/Users/joshbailey/downloads/FullInventory.xls")
	} else {
		log.Print(err)
	}
}
