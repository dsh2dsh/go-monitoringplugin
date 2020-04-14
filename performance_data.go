/* Copyright (c) 2019, inexio GmbH, BSD 2-Clause License */

package monitoringplugin

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
)

/*
PerformanceData is a map where all performanceDataPoints are stored. It assigns labels (string) to performanceDataPoints.
*/
type PerformanceData map[string]PerformanceDataPoint

/*
PerformanceDataPoint contains all information of one PerformanceDataPoint.
*/
type PerformanceDataPoint struct {
	label string
	value float64
	unit  string
	warn  float64 //currently we do not support ranges for warning and critical thresholds, because icinga2 does not support it
	crit  float64
	min   float64
	max   float64

	labelTag string

	hasWarn bool
	hasCrit bool
	hasMin  bool
	hasMax  bool
}

/*
add adds a PerformanceDataPoint to the PerformanceData Map.
The function checks if a PerformanceDataPoint is valid and if there is already another PerformanceDataPoint with the same label in the PerformanceData map.
Usage:
	err := PerformanceData.add(NewPerformanceDataPoint("temperature", 32, "°C").SetWarn(35).SetCrit(40))
	if err != nil {
		...
	}
*/
func (p *PerformanceData) add(point *PerformanceDataPoint) error {
	if err := point.validate(); err != nil {
		return errors.Wrap(err, "given performance data point is not valid")
	}
	if _, ok := (*p)[point.label+point.labelTag]; ok {
		return errors.New("a performance data point with this label does already exist")
	}
	(*p)[point.label+point.labelTag] = *point
	return nil
}

/*
Validates a PerformanceDataPoint.
This function is used to check if a PerformanceDataPoint is compatible with the documentation from 'http://nagios-plugins.org/doc/guidelines.html'(valid name and unit, value is inside the range of min and max etc.)
*/
func (p *PerformanceDataPoint) validate() error {
	if p.label == "" {
		return errors.New("data point label cannot be an empty string")
	}

	match, err := regexp.MatchString("([='])", p.label)
	if err != nil {
		return errors.Wrap(err, "error during regex match")
	}
	if match {
		return errors.New("label can not contain the equal sign or single quote (')")
	}

	match, err = regexp.MatchString("([0-9;'\"])", p.unit)
	if err != nil {
		return errors.Wrap(err, "error during regex match")
	}
	if match {
		return errors.New("unit can not contain numbers, semicolon or quotes")
	}

	if (p.hasMin && p.hasMax) && (p.min > p.max) {
		return errors.New("min cannot be larger than max")
	}
	if p.hasMin && p.value < p.min {
		return errors.New("value cannot be smaller than min")
	}
	if p.hasMax && p.value > p.max {
		return errors.New("value cannot be larger than max")
	}
	return nil
}

/*
NewPerformanceDataPoint creates a new PerformanceDataPoint. Label and value are mandatory but are not checked at this point, the performanceDatePoint's validation is checked later when it is added to the PerformanceData list in the function PerformanceData.add(*PerformanceDataPoint).
It is possible to directly add warning, critical, min and max values with the functions SetWarn(int), SetCrit(int), SetMin(int) and SetMax(int).
Usage:
	PerformanceDataPoint := NewPerformanceDataPoint("memory_usage", 55, "%").SetWarn(80).SetCrit(90)
*/
func NewPerformanceDataPoint(label string, value float64, unit string) *PerformanceDataPoint {
	return &PerformanceDataPoint{
		label:    label,
		value:    value,
		unit:     unit,
		labelTag: "",
	}
}

/*
This function returns the PerformanceDataPoint in the specified string format that will be returned by the check plugin.
*/
func (p *PerformanceDataPoint) outputString(jsonLabel bool) string {
	var outputString string
	if jsonLabel {
		outputString += `'{"metric":"` + p.label + `"`
		if p.labelTag != "" {
			outputString += `,"label":"` + p.labelTag + `"`
		}
		outputString += "}'"
	} else {
		outputString += "'" + p.label
		if p.labelTag != "" {
			outputString += "_" + p.labelTag
		}
		outputString += "'"
	}
	outputString += "=" + fmt.Sprintf("%g", p.value) + p.unit + ";"
	if p.hasWarn {
		outputString += fmt.Sprintf("%g", p.warn)
	}
	outputString += ";"
	if p.hasCrit {
		outputString += fmt.Sprintf("%g", p.crit)
	}
	outputString += ";"
	if p.hasMin {
		outputString += fmt.Sprintf("%g", p.min)
	}
	outputString += ";"
	if p.hasMax {
		outputString += fmt.Sprintf("%g", p.max)
	}

	return outputString
}

/*
SetMin sets minimum value.
*/
func (p *PerformanceDataPoint) SetMin(min float64) *PerformanceDataPoint {
	p.hasMin = true
	p.min = min
	return p
}

/*
SetMax sets maximum value.
*/
func (p *PerformanceDataPoint) SetMax(max float64) *PerformanceDataPoint {
	p.hasMax = true
	p.max = max
	return p
}

/*
SetWarn sets maximum value.
*/
func (p *PerformanceDataPoint) SetWarn(warn float64) *PerformanceDataPoint {
	p.hasWarn = true
	p.warn = warn
	return p
}

/*
SetCrit sets critical value.
*/
func (p *PerformanceDataPoint) SetCrit(crit float64) *PerformanceDataPoint {
	p.hasCrit = true
	p.crit = crit
	return p
}

/*
SetLabelTag adds a tag to the performance data point
If one tag is added more than once, the value before will be overwritten
*/
func (p *PerformanceDataPoint) SetLabelTag(labelTag string) *PerformanceDataPoint {
	p.labelTag = labelTag
	return p
}
