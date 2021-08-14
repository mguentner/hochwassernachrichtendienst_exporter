package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

func GetTimeZone() *time.Location {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		log.Fatal().Msg("Could not load timezone Europe/Berlin, using time.Local, this might be wrong.")
		location = time.Local
	}
	return location
}

type WaterLevelInformation struct {
	StationName                    string
	StationID                      string
	LastMeasurementTime            time.Time
	LastLevelCentimeter            float64
	LastOutFlowCubicmeterPerSecond float64
	CurrentWarningLevel            int
	WarningLevel1Centimeter        float64
	WarningLevel2Centimeter        float64
	WarningLevel3Centimeter        float64
	WarningLevel4Centimeter        float64
	// Definition:
	// An event that reaches or surpases that level with a probability of 1%
	// per year, statistically this happens 100 times in 10.000 years
	HundredYearFloodLevelCentimeter float64
}

func (a *WaterLevelInformation) IsEqual(b *WaterLevelInformation) (bool, []error) {
	errs := []error{}
	if a.StationName != b.StationName {
		errs = append(errs, errors.New(fmt.Sprintf("StationName not equal: %s != %s", a.StationName, b.StationName)))
	}
	if a.StationID != b.StationID {
		errs = append(errs, errors.New(fmt.Sprintf("StationID not equal: %s != %s", a.StationID, b.StationID)))
	}
	if !a.LastMeasurementTime.Equal(b.LastMeasurementTime) {
		errs = append(errs, errors.New(fmt.Sprintf("LastMeasurement not equal: %v != %v", a.LastMeasurementTime, b.LastMeasurementTime)))
	}
	if a.LastLevelCentimeter != b.LastLevelCentimeter {
		errs = append(errs, errors.New(fmt.Sprintf("LastLevelCentimeter not equal: %f != %f", a.LastLevelCentimeter, b.LastLevelCentimeter)))
	}
	if a.LastOutFlowCubicmeterPerSecond != b.LastOutFlowCubicmeterPerSecond {
		errs = append(errs, errors.New(fmt.Sprintf("LastOutFlowCubicmeterPerSecond not equal: %f != %f", a.LastOutFlowCubicmeterPerSecond, b.LastOutFlowCubicmeterPerSecond)))
	}
	if a.CurrentWarningLevel != b.CurrentWarningLevel {
		errs = append(errs, errors.New(fmt.Sprintf("CurrentWarningLevel not equal: %d != %d", a.CurrentWarningLevel, b.CurrentWarningLevel)))
	}

	if a.WarningLevel1Centimeter != b.WarningLevel1Centimeter {
		errs = append(errs, errors.New(fmt.Sprintf("WarningLevel1Centimeter not equal: %f != %f", a.WarningLevel1Centimeter, b.WarningLevel1Centimeter)))
	}
	if a.WarningLevel2Centimeter != b.WarningLevel2Centimeter {
		errs = append(errs, errors.New(fmt.Sprintf("WarningLevel2Centimeter not equal: %f != %f", a.WarningLevel2Centimeter, b.WarningLevel2Centimeter)))
	}
	if a.WarningLevel3Centimeter != b.WarningLevel3Centimeter {
		errs = append(errs, errors.New(fmt.Sprintf("WarningLevel3Centimeter not equal: %f != %f", a.WarningLevel3Centimeter, b.WarningLevel3Centimeter)))
	}
	if a.WarningLevel4Centimeter != b.WarningLevel4Centimeter {
		errs = append(errs, errors.New(fmt.Sprintf("WarningLevel4Centimeter not equal: %f != %f", a.WarningLevel4Centimeter, b.WarningLevel4Centimeter)))
	}
	if a.HundredYearFloodLevelCentimeter != b.HundredYearFloodLevelCentimeter {
		errs = append(errs, errors.New(fmt.Sprintf("HundredYearFloodLevelCentimeter not equal: %f != %f", a.HundredYearFloodLevelCentimeter, b.HundredYearFloodLevelCentimeter)))
	}
	return len(errs) == 0, errs
}

func PositiveCentimeterToMillimeter(centimeter float64) float64 {
	if centimeter < 0 { // -1 should remain -1
		return centimeter
	} else {
		return centimeter*10
	}
}

func (wl *WaterLevelInformation) AddMetricsToRegistry(registry *prometheus.Registry) error {
	addAndSetGauge := func(name string, help string, value float64) error {
		g := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:  "hnd", // hochwassernachrichtendienst
			Name:        name,
			Help:        help,
			ConstLabels: map[string]string{
				"id": wl.StationID,
				"name": wl.StationName,
			},
		})
		err := registry.Register(g)
		if err != nil {
			return err
		}
		g.Set(value)
		return nil
	}
	err := addAndSetGauge("last_level_millimeters", "last level in millimeters", PositiveCentimeterToMillimeter(wl.LastLevelCentimeter))
	if err != nil {
		return err
	}
	err = addAndSetGauge("last_outflow_cubicmeters_per_second", "last outflow in cubicmeters/second, -1 if not present", wl.LastOutFlowCubicmeterPerSecond)
	if err != nil {
		return err
	}
	err = addAndSetGauge("current_warning_level", "current reported warning level", float64(wl.CurrentWarningLevel))
	if err != nil {
		return err
	}
	err = addAndSetGauge("warning_level_1_millimeters", "warning level 1 in millimeters, -1 if not present", PositiveCentimeterToMillimeter(wl.WarningLevel1Centimeter))
	if err != nil {
		return err
	}
	err = addAndSetGauge("warning_level_2_millimeters", "warning level 2 in millimeters, -1 if not present", PositiveCentimeterToMillimeter(wl.WarningLevel2Centimeter))
	if err != nil {
		return err
	}
	err = addAndSetGauge("warning_level_3_millimeters", "warning level 3 in millimeters, -1 if not present", PositiveCentimeterToMillimeter(wl.WarningLevel3Centimeter))
	if err != nil {
		return err
	}
	err = addAndSetGauge("warning_level_4_millimeters", "warning level 4 in millimeters, -1 if not present", PositiveCentimeterToMillimeter(wl.WarningLevel4Centimeter))
	if err != nil {
		return err
	}
	err = addAndSetGauge("hundered_year_flood_level_millimeters", "an event that reaches or surpasses that level with a probability of 1% per year, -1 if not present", PositiveCentimeterToMillimeter(wl.HundredYearFloodLevelCentimeter))
	if err != nil {
		return err
	}
	return nil
}

type ParserState int

const (
	Initial             ParserState = iota
	Name                ParserState = iota
	CurrentWarningLevel ParserState = iota
	HundredYearFlood1   ParserState = iota
	HundredYearFlood2   ParserState = iota
	HundredYearFlood3   ParserState = iota
	LastLevelCentimeter ParserState = iota
	LastOutFlowCubic    ParserState = iota
	LastValueTime       ParserState = iota
)

func ParseWaterLevelTime(value string) (*time.Time, error) {
	splitted := strings.Split(value, ",")
	if len(splitted) != 2 {
		return nil, errors.New("Wrong split value")
	}
	dateElements := strings.Split(splitted[0], ".")
	if len(dateElements) != 3 {
		return nil, errors.New("Wrong amount of date values")
	}
	timeElements := strings.Split(splitted[1], ":")
	if len(timeElements) != 2 {
		return nil, errors.New("Wrong amount of time values")
	}
	year, err := strconv.ParseInt(strings.TrimSpace(dateElements[2]), 10, 32)
	if err != nil || year > 99 {
		return nil, errors.New("Failed to parse year")
	}
	month, err := strconv.ParseInt(strings.TrimSpace(dateElements[1]), 10, 32)
	if err != nil || month > 12 {
		return nil, errors.New("Failed to parse month")
	}
	day, err := strconv.ParseInt(strings.TrimSpace(dateElements[0]), 10, 32)
	if err != nil || day > 31 {
		return nil, errors.New("Failed to parse day")
	}
	hour, err := strconv.ParseInt(strings.TrimSpace(timeElements[0]), 10, 32)
	if err != nil || hour > 23 {
		return nil, errors.New("Failed to parse hour")
	}
	minute, err := strconv.ParseInt(strings.TrimSpace(timeElements[1]), 10, 32)
	if err != nil || minute > 59 {
		return nil, errors.New("Failed to parse minutes")
	}
	t := time.Date(int(2000+year), time.Month(month), int(day), int(hour), int(minute), 0, 0, GetTimeZone())
	return &t, nil
}

// parse HTML rendered by https://m.hnd.bayern.de/pegel.php
func ParseWaterLevel(r io.Reader) (*WaterLevelInformation, error) {
	t := html.NewTokenizer(r)
	result := WaterLevelInformation{
		StationName:                     "",
		StationID:                     "",
		LastMeasurementTime:             time.Now(),
		LastLevelCentimeter:             -1,
		LastOutFlowCubicmeterPerSecond:  -1,
		CurrentWarningLevel:             -1,
		WarningLevel1Centimeter:         -1,
		WarningLevel2Centimeter:         -1,
		WarningLevel3Centimeter:         -1,
		WarningLevel4Centimeter:         -1,
		HundredYearFloodLevelCentimeter: -1,
	}
	parserState := Initial
	for {
		next := t.Next()
		switch next {
		case html.ErrorToken:
			err := t.Err()
			if err == io.EOF {
				return &result, nil
			}
			return &result, err
		case html.TextToken:
			text := string(t.Text())
			switch parserState {
			case Name:
				result.StationName += text
			case Initial:
				if strings.Contains(text, "Meldestufe:") {
					parserState = CurrentWarningLevel
				} else if strings.Contains(text, "Meldestufe ") {
					fields := strings.Fields(text)
					centimeterValue, err := strconv.ParseFloat(fields[2], 64)
					if err != nil {
						log.Warn().Msg("Could not parse CentimeterValue of WarningLevel")
						centimeterValue = -1
					}
					switch fields[1] {
					case "1:":
						result.WarningLevel1Centimeter = centimeterValue
					case "2:":
						result.WarningLevel2Centimeter = centimeterValue
					case "3:":
						result.WarningLevel3Centimeter = centimeterValue
					case "4:":
						result.WarningLevel4Centimeter = centimeterValue
					default:
						log.Warn().Msgf("Encountered Unknown Warning Level: %s", fields[1])
					}
				} else if strings.Contains(text, "IÜG: HQ") {
					parserState = HundredYearFlood1
				} else if strings.Contains(text, "Wasserstand [cm]:") {
					parserState = LastLevelCentimeter
				} else if strings.Contains(text, "Abfluss [m³/s]") {
					parserState = LastOutFlowCubic
				} else if strings.Contains(text, "Letzter Wert:") {
					parserState = LastValueTime
				}
			case CurrentWarningLevel:
				level, err := strconv.ParseInt(text, 10, 64)
				if err != nil {
					log.Warn().Msg("Could not parse CurrentWarningLevel")
					result.CurrentWarningLevel = -1
				} else {
					result.CurrentWarningLevel = int(level)
				}
				parserState = Initial
			case HundredYearFlood2:
				if text == "100" {
					parserState = HundredYearFlood3
				}
			case HundredYearFlood3:
				fields := strings.Fields(text)
				centimeterValue, err := strconv.ParseFloat(fields[1], 64)
				if err != nil {
					log.Warn().Msg("Could not parse CentimeterValue of HunderedYearFlood")
					centimeterValue = -1
				}
				result.HundredYearFloodLevelCentimeter = centimeterValue
				parserState = Initial
			case LastLevelCentimeter:
				centimeterValue, err := strconv.ParseFloat(text, 64)
				if err != nil {
					log.Warn().Msg("Could not parse CentimeterValue of LastLevelCentimeter")
					centimeterValue = -1
				}
				result.LastLevelCentimeter = centimeterValue
				parserState = Initial
			case LastOutFlowCubic:
				cubicmeterValue, err := strconv.ParseFloat(strings.Replace(text, ",", ".", 1), 64)
				if err != nil {
					log.Warn().Msg("Could not parse cubicMeterValue of LastOutFlowCubicmeterPerSecond")
					cubicmeterValue = -1
				}
				result.LastOutFlowCubicmeterPerSecond = cubicmeterValue
				parserState = Initial
			case LastValueTime:
				parsedTime, err := ParseWaterLevelTime(text)
				if err != nil {
					log.Warn().Msgf("Could not parse date and time: %s", err.Error())
					result.LastMeasurementTime = time.Now()
				} else {
					result.LastMeasurementTime = *parsedTime
				}
				parserState = Initial
			}

		case html.SelfClosingTagToken:
			switch parserState {
			case Name:
				if tagName, _ := t.TagName(); string(tagName) == "br" {
					result.StationName += " "
				}
			}
		case html.StartTagToken:
			attrs := make(map[string]string)
			tagName, hasAttr := t.TagName()
			if hasAttr {
				moreAttrs := false
				for {
					key, val, more := t.TagAttr()
					moreAttrs = more
					attrs[string(key)] = string(val)
					if moreAttrs == false {
						break
					}
				}
			}
			switch parserState {
			case Initial:
				if string(tagName) == "span" {
					if class, ok := attrs["class"]; ok {
						if class == "header" {
							parserState = Name
						}
					}
				}
			case Name:
				if string(tagName) == "span" {
					parserState = Initial
				}
			case HundredYearFlood1:
				if string(tagName) == "sub" {
					parserState = HundredYearFlood2
				}
			}
		case html.EndTagToken:
			switch parserState {
			case Name:
				parserState = Initial
			}
		}
	}
}
