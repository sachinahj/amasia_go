package yelp

import "fmt"

func Start() {
	fmt.Println("=================================================")
	var lg Log
	var zc ZipCode

	lg.InitWithNotDoneZipCodeAnalysis()
	if lg.Id != 0 {
		zc.InitWithOldestBusinessesSearch()

		if zc.ZipCode != 0 {
			zc.ZipCode = lg.ZipCode
			zc.Analyze()
			lg.IsDone = true
			lg.Update()
		}

	} else {

		zc.InitWithForceBusinessesSearch()

		if zc.ZipCode != 0 {
			fmt.Println("---Forced ZipCode \n", zc)
			lg.ZipCode = zc.ZipCode
			lg.InitWithNewBusinessesSearch()
			lg.Insert()
		} else {
			lg.InitWithLatestBusinessesSearch()
			fmt.Println("---No Forced ZipCode, Latest Log Businesses Search \n", lg)
		}

		if lg.ZipCode == 0 {
			zc.InitWithOldestBusinessesSearch()

			if zc.ZipCode != 0 {
				lg.ZipCode = zc.ZipCode
				lg.InitWithNewBusinessesSearch()
				lg.Insert()
				fmt.Println("---No Latest Log, Oldest ZipCode Log \n", lg)
			}
		}

		if lg.Id == 0 {
			fmt.Println("---Nothing to do \n")
			return
		}

		if !lg.IsDone {
			BusinessesSearch(&lg)
			lg.IsDone = true
			lg.Update()
		}

	}

	lg.InitWithNextLog()
	lg.Insert()
	Start()
}
