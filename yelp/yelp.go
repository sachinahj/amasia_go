package yelp

import "fmt"

func Start() {
	fmt.Println("======================= Start ==========================")
	var lg Log
	var zc ZipCode

	lg.InitWithNotDoneZipCodeAnalysis()
	if lg.Id != 0 {

		zc.ZipCode = lg.ZipCode
		zc.InitWithZipCode()
		zc.RunAnalysis()
		lg.IsDone = true
		lg.Update()

	} else {

		lg.InitWithLatestBusinessesSearch()

		if lg.Id != 0 && !lg.IsDone {
			BusinessesSearch(&lg)
			lg.IsDone = true
			lg.Update()
		}

	}

	lg.InitWithNextLog()
	lg.Insert()
	Start()
}
