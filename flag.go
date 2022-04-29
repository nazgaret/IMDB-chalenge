package main

import "flag"

const (
	maxApiRequestsDefault = 100
	maxRunTimeDefault     = 300
)

type Config struct {
	filePath                   string
	maxApiRequests, maxRunTime int
}

func parseFlags(config *Config, filter *Filter) {
	//config part
	flag.StringVar(&config.filePath, "filePath", "title.basics.tsv.gz", "path to IMDB data tar/gz file")
	flag.IntVar(&config.maxApiRequests, "maxApiRequestsDefault", maxApiRequestsDefault, "maxApiRequestsDefault to IMDB API")
	flag.IntVar(&config.maxRunTime, "maxRunTimeDefault", maxRunTimeDefault, "maxRunTimeDefault in seconds")

	//filter part
	flag.StringVar(&filter[FilterTitleTypeIndex], "titleType", "", "titleType filter")
	flag.StringVar(&filter[FilterPrimaryTitleIndex], "primaryTitle", "", "primaryTitle filter")
	flag.StringVar(&filter[FilterOriginalTitleIndex], "originalTitle", "", "originalTitle filter")
	flag.StringVar(&filter[FilterStartYearIndex], "startYear", "", "startYear filter")
	flag.StringVar(&filter[FilterEndYearIndex], "endYear", "", "endYear filter")
	flag.StringVar(&filter[FilterRuntimeMinutesIndex], "runtimeMinutes", "", "runtimeMinutes filter")
	flag.StringVar(&filter[FilterGenresIndex], "genres", "", "genres filter")

	flag.Parse()
}
