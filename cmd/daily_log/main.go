package main

import (
	"context"
	"log"
	"os"

	"github.com/chromedp/chromedp"
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

	// run task to process units available table
	a, err := processUnitsAvailable("/Users/joshbailey/downloads/UnitsAvailable.xls")
	if err != nil {
		log.Fatalf("error getting available units: %v", err)
	}
	os.Remove("/Users/joshbailey/downloads/UnitsAvailable.xls")
	var _ = a

	// run task to process contracts open table
	c, err := processContractsOpen("/Users/joshbailey/downloads/ClosedContractsOpen.xls")
	if err != nil {
		log.Fatalf("error getting open contracts: %v", err)
	}
	os.Remove("/Users/joshbailey/downloads/ClosedContractsOpen.xls")
	var _ = c

	// run task to process utilization table
	u, err := processUtilization("/Users/joshbailey/downloads/AdjustedUtilization.csv")
	if err != nil {
		log.Fatalf("error getting Utilization: %v", err)
	}
	os.Remove("/Users/joshbailey/downloads/AdjustedUtilization.csv")
	var _ = u

	// run task to process advisor duration table
	d, err := processDuration("/Users/joshbailey/downloads/Service_Advisor_Activity.csv")
	if err != nil {
		log.Fatalf("error getting Duration: %v", err)
	}
	os.Remove("/Users/joshbailey/downloads/Service_Advisor_Activity.csv")
	var _ = d

	// run task to process full inventory table
	f, err := processFullInventory("/Users/joshbailey/downloads/FullInventory.xls")
	if err != nil {
		log.Fatalf("error getting full inventory: %v", err)
	}
	os.Remove("/Users/joshbailey/downloads/FullInventory.xls")
	var _ = f
}
