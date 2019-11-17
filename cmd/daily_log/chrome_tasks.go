package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/extrame/xls"
	"github.com/imjoshuabailey/arnold_automator/db"
	"github.com/imjoshuabailey/arnold_automator/lib"
)

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
		chromedp.Sleep(2 * time.Second),
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

// export FullInventory.xls file
func exportFullInventory() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(`https://tsdloaner.tsd-inc.com/FullInventory.aspx`),
		chromedp.Sleep(2 * time.Second),
		chromedp.Click(`//input[@name="butExport"]`),
	}
}

// wait 10 seconds before moving on
func idleChrome() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Sleep(10 * time.Second),
	}
}

// count the number of vehicles units available by year
func processUnitsAvailable(path string) (int, error) {
	var count17 = 0
	var count18 = 0
	var count19 = 0
	var count20 = 0
	var total = 0
	if xlFile, err := xls.Open(path, "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			year := sheet1.Row(0).Col(3)

			for i := 0; i <= (int(sheet1.MaxRow)); i++ {
				row := sheet1.Row(i)
				year = row.Col(3)
				if year == "2017" || year == "17" {
					count17 = count17 + 1
					total = total + 1
				} else if year == "2018" || year == "18" {
					count18 = count18 + 1
					total = total + 1
				} else if year == "2019" || year == "19" {
					count19 = count19 + 1
					total = total + 1
				} else if year == "2020" || year == "20" {
					count20 = count20 + 1
					total = total + 1
				}
			}
		}
	} else {
		return 0, err
	}
	fmt.Println(count17, "2017s available")
	fmt.Println(count18, "2018s available")
	fmt.Println(count19, "2019s available")
	fmt.Println(count20, "2020s available")
	fmt.Println(total, "Vehicles Available")

	return total, nil
}

// count the number of vehicles loaned out >3 days and >9 days
func processContractsOpen(path string) (int, error) {
	var greaterThanThree = 0
	var doubleDigit = 0
	if xlFile, err := xls.Open(path, "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			days := sheet1.Row(3).Col(15)

			for i := 0; i <= (int(sheet1.MaxRow)); i++ {
				row := sheet1.Row(i)
				days = row.Col(15)
				if days, err := strconv.Atoi(days); err == nil {
					if days > 3 {
						greaterThanThree = greaterThanThree + 1
						if days > 9 {
							doubleDigit = doubleDigit + 1
						}
					}
				}
			}
		}

	} else {
		return 0, err
	}
	fmt.Println(greaterThanThree, "vehicles out more than 3 days")
	fmt.Println(doubleDigit, "vehicles out more than 9 days")
	return greaterThanThree, nil
}

// parse vehicle utilization percentage
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

// parse the current loan duration average
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

// count total number of vehicles by year
func processFullInventory(path string) (int, error) {
	var count2017 = 0
	var count2018 = 0
	var count2019 = 0
	var count2020 = 0
	var total = 0

	var (
		host     = os.Getenv(`POSTGRES_HOST`)
		port     = os.Getenv(`POSTGRES_PORT`)
		user     = os.Getenv(`POSTGRES_USER`)
		password = os.Getenv(`POSTGRES_PASSWORD`)
		dbname   = os.Getenv(`POSTGRES_DB`)
	)

	conn, err := db.ConnectDB(host, port, user, password, dbname)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer conn.Close()

	if xlFile, err := xls.Open(path, "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {

			for i := 1; i <= (int(sheet1.MaxRow)); i++ {
				r := sheet1.Row(i)
				date := r.Col(3)

				var v lib.Vehicle

				// assign column value to field on struct
				v.VIN = r.Col(10)
				v.UnitNumber = r.Col(0)
				v.Class = r.Col(2)
				yearStr := r.Col(3)
				year, err := strconv.Atoi(yearStr)
				if err != nil {
					log.Printf("error converting year to integer: %v", err)
					continue
				}
				v.Year = year
				v.Model = r.Col(4)
				v.BodyStyle = r.Col(5)
				v.Color = r.Col(6)
				v.LicenseNumber = r.Col(7)
				milesStr := r.Col(8)
				miles, err := strconv.Atoi(milesStr)
				if err != nil {
					log.Printf("error converting miles to integer: %v", err)
					continue
				}
				v.Miles = miles
				v.Status = sortStatus(r.Col(11))
				if v.Status == 2 {
					v.ContractNumber = getContractNum(r.Col(11))
				} else {
					v.ContractNumber = 0
				}

				v.DateLastUsed = r.Col(12)
				v.StockNumber = r.Col(13)
				// insert into db
				err = db.Insert(conn, v)
				if err != nil {
					log.Fatalf("error inserting into db: %v", err)
				}
				if date == "2017" || date == "17" {
					count2017 = count2017 + 1
					total = total + 1
				} else if date == "2018" || date == "18" {
					count2018 = count2018 + 1
					total = total + 1
				} else if date == "2019" || date == "19" {
					count2019 = count2019 + 1
					total = total + 1
				} else if date == "2020" || date == "20" {
					count2020 = count2020 + 1
					total = total + 1
				}
			}
		}
	} else {
		return 0, err
	}

	fmt.Println(total, "Vehicles in fleet")
	fmt.Println(count2017, "2017s")
	fmt.Println(count2018, "2018s")
	fmt.Println(count2019, "2019s")
	fmt.Println(count2020, "2020s")
	return total, nil
}

func sortStatus(status string) int {
	var statusNum int

	if strings.HasPrefix(status, "Available") {
		statusNum = 1
	} else if strings.HasPrefix(status, "On Loan") {
		statusNum = 2
	} else if strings.HasPrefix(status, "High") {
		statusNum = 3
	} else if strings.HasPrefix(status, "Body") {
		statusNum = 4
	} else {
		statusNum = 0
	}

	return statusNum
}

const Unknown = 0
const Available = 1
const OnLoan = 2
const HighMiles = 3
const Damaged = 4

// switch status {
// case "Available":
//    statusNum = 1
// case strings.HasPrefix(status, "On Loan"):
//     statusNum = 2
// case strings.HasPrefix(status, "High"):
//    statusNum = 3
// case strings.HasPrefix(status,"Body"):
//    statusNum = 4
// default:
//    statusNum = 0
// }

func getContractNum(status string) int {
	var contractNum int

	_, err := fmt.Sscanf(status, "On Loan (%d)", &contractNum)
	if err != nil {
		log.Fatalf("Error scanning contract number: %v", err)
	}

	return contractNum

}
