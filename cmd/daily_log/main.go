package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	var err error

	// create chrome instance
	c, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	// run task authenticates user
	err = chromedp.Run(c, authenticate(`https://tsdloaner.tsd-inc.com`, `54106`, `jbailey`, `Carwash1`))
	if err != nil {
		log.Fatal(err)
	}

	// run task exports units available report
	err = chromedp.Run(c, exportUnitsAvailable())
	if err != nil {
		log.Fatal(err)
	}

	// run task exports unit duration report
	err = chromedp.Run(c, exportUnitDuration())
	if err != nil {
		log.Fatal(err)
	}

	// run task exports unit utilization report
	err = chromedp.Run(c, exportUnitUtilization())
	if err != nil {
		log.Fatal(err)
	}

	// run task exports current units out report
	err = chromedp.Run(c, exportUnitsOut())
	if err != nil {
		log.Fatal(err)
	}

	// run task wait 10 seconds
	err = chromedp.Run(c, idleChrome())
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = chromedp.Cancel(c)
	if err != nil {
		log.Fatal(err)
	}

	// // wait for chrome to finish
	// err = chromedp.Wait(c)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// processUnitsAvailable()

	// processContractsOpen()

	// u, err := processUtilization("/Users/joshbailey/downloads/AdjustedUtilization.csv")
	// if err != nil {
	// 	log.Fatalf("error getting Utilization: %v", err)
	// }

	// var _ = u

	// d, err := processDuration("/Users/joshbailey/downloads/Service_Advisor_Activity.csv")
	// if err != nil {
	// 	log.Fatalf("error getting Duration: %v", err)
	// }

	// var _ = d
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
		chromedp.WaitVisible(`//input[@name="butExport"]`),
		chromedp.Click(`//input[@name="butExport"]`),
	}
}

// export Service_Advisor_Activity.csv file
func exportUnitDuration() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://tsdloaner.tsd-inc.com/ServiceWriters.aspx`),
		chromedp.Click(`//input[@name="butRun"]`),
		chromedp.Click(`//input[@name="btnExport"]`),
		chromedp.Click(`//input[@name="butClose"]`),
	}
}

// export AdjustedUtilization.csv file
func exportUnitUtilization() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://tsdloaner.tsd-inc.com/AdjustedUnitUtil.aspx`),
		chromedp.Click(`//input[@name="butRun"]`),
		chromedp.Click(`//input[@name="btnExport"]`),
		chromedp.Click(`//input[@name="butClose"]`),
	}
}

// export ClosedContractsOpen.xls file
func exportUnitsOut() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://tsdloaner.tsd-inc.com/CurrOpenRa.aspx`),
		chromedp.Click(`//input[@name="butExport"]`),
	}
}

func idleChrome() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Sleep(10 * time.Second),
	}
}

// // count the number of 2019 vehicles in units available file
// func processUnitsAvailable() {
// 	if xlFile, err := xls.Open("/Users/joshbailey/downloads/UnitsAvailable.xls", "utf-8"); err == nil {
// 		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
// 			year := sheet1.Row(0).Col(3)
// 			var count = 0
// 			for i := 0; i <= (int(sheet1.MaxRow)); i++ {
// 				row1 := sheet1.Row(i)
// 				year = row1.Col(3)
// 				if year == "2019" || year == "19" {
// 					count = count + 1
// 				}

// 			}
// 			fmt.Println(count, "Vehicles Available")
// 		}
// 	} else {
// 		log.Print(err)
// 	}
// }

// func processContractsOpen() {
// 	if xlFile, err := xls.Open("/Users/joshbailey/downloads/ClosedContractsOpen.xls", "utf-8"); err == nil {
// 		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
// 			days := sheet1.Row(3).Col(15)
// 			var greaterThanThree = 0
// 			var doubleDigit = 0
// 			for i := 0; i <= (int(sheet1.MaxRow)); i++ {
// 				row1 := sheet1.Row(i)
// 				days = row1.Col(15)
// 				if days, err := strconv.Atoi(days); err == nil {
// 					if days > 3 {
// 						greaterThanThree = greaterThanThree + 1
// 						if days > 9 {
// 							doubleDigit = doubleDigit + 1
// 						}
// 					}
// 				}

// 			}
// 			fmt.Println(greaterThanThree, "vehicles out more than 3 days")
// 			fmt.Println(doubleDigit, "vehicles out more than 9 days")
// 		}
// 	} else {
// 		log.Print(err)
// 	}
// }

// func processUtilization(path string) (float64, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer f.Close()

// 	r := csv.NewReader(f)
// 	records, err := r.ReadAll()
// 	if err != nil {
// 		return 0, err
// 	}

// 	removeSybol := strings.Index(records[len(records)-1][7], "%")
// 	value, err := strconv.ParseFloat(records[len(records)-1][7][:removeSybol], 64)

// 	if err != nil {
// 		return 0, err
// 	}
// 	fmt.Println(value)
// 	return value, nil
// }

// func processDuration(path string) (float64, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer f.Close()

// 	r := csv.NewReader(f)
// 	records, err := r.ReadAll()
// 	if err != nil {
// 		return 0, err
// 	}

// 	value, err := strconv.ParseFloat(records[len(records)-1][7], 64)

// 	if err != nil {
// 		return 0, err
// 	}
// 	fmt.Println(value)
// 	return value, nil

// }

// // run task exports full inventory report
// err = c.Run(ctxt, exportFullInventory())
// if err != nil {
// 	log.Fatal(err)
// }

// //export FullInventory.xls file
// func exportFullInventory() chromedp.Tasks {
// 	return chromedp.Tasks{
// 		chromedp.Navigate(`https://tsdloaner.tsd-inc.com/FullInventory.aspx`),
// 		chromedp.Click(`//input[@name="butExport"]`),
// 	}
// }
