package main

import (
	"os"
	"testing"
	"time"
)

func TestParseWaterLevel(t *testing.T) {
	testData := []struct {
		path string
		expected WaterLevelInformation
	}{
		{
			path: "sample_data/sample_muenchen.html",
			expected: WaterLevelInformation{
				StationName:                     "München Isar",
				LastMeasurementTime:                 time.Date(2021, 8, 12, 17, 15, 0, 0, GetTimeZone()),
				LastLevelCentimeter:              143,
				LastOutFlowCubicmeterPerSecond:    109,
				CurrentWarningLevel:             0,
				WarningLevel1Centimeter:         240,
				WarningLevel2Centimeter:         300,
				WarningLevel3Centimeter:         380,
				WarningLevel4Centimeter:         520,
				HundredYearFloodLevelCentimeter: 510,
			},
		},
		{
			path: "sample_data/sample_lenggries.html",
			expected: WaterLevelInformation{
				StationName:                     "Lenggries Isar",
				LastMeasurementTime:                 time.Date(2021, 8, 12, 20, 30, 0, 0, GetTimeZone()),
				LastLevelCentimeter:              120,
				LastOutFlowCubicmeterPerSecond:    24.4,
				CurrentWarningLevel:             0,
				WarningLevel1Centimeter:         220,
				WarningLevel2Centimeter:         -1,
				WarningLevel3Centimeter:         310,
				WarningLevel4Centimeter:         -1,
				HundredYearFloodLevelCentimeter: 390,
			},
		},
	}

	for _, td := range testData {
		r, err := os.Open(td.path)
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := ParseWaterLevel(r)
		if err != nil {
			t.Error(err)
		}
		if parsed == nil {
			t.Fatal("Expected parsed to be not nil")
		}
		equal, errs := parsed.IsEqual(&td.expected)
		if !equal {
			t.Error("Expected extracted data to be equal")
		}
		for _, err := range errs {
			t.Error(err)
		}
	}
}
