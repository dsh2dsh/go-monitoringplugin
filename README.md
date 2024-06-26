# go-monitoringplugin

[![Go](https://github.com/dsh2dsh/go-monitoringplugin/actions/workflows/go.yml/badge.svg)](https://github.com/dsh2dsh/go-monitoringplugin/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/dsh2dsh/go-monitoringplugin/v2.svg)](https://pkg.go.dev/github.com/dsh2dsh/go-monitoringplugin/v2)

## Description

Golang package for writing monitoring check plugins for
[nagios](https://www.nagios.org/), [icinga2](https://icinga.com/),
[zabbix](https://www.zabbix.com/), [checkmk](https://checkmk.com/), etc.

The package complies with the [Monitoring Plugins Development Guidelines](https://www.monitoring-plugins.org/doc/guidelines.html).

This project is a fork of
[go-monitoringplugin](https://github.com/inexio/go-monitoringplugin).

## Changes from [upstream](https://github.com/inexio/go-monitoringplugin):

  * Use generics, instead of big.Float and big.ParseFloat. See
    [#1](https://github.com/inexio/go-monitoringplugin/pull/1).

## Example / Usage

``` go
package main

import "github.com/inexio/go-monitoringplugin/v2"

func main() {
	// Creating response with a default ok message, that will be displayed when
	// the checks exits with status ok.
	response := monitoringplugin.NewResponse("everything checked!")

	// Set output delimiter (default is \n)
	// response.SetOutputDelimiter(" / ")

	// Updating check plugin status and adding message to the output (status only
	// changes if the new status is worse than the current one).

	// check status stays ok
	response.UpdateStatus(monitoringplugin.OK, "something is ok!")
	// check status updates to critical
	response.UpdateStatus(monitoringplugin.CRITICAL,
		"something else is critical!")
	// check status stays critical, but message will be added to the output
	response.UpdateStatus(monitoringplugin.WARNING, "something else is warning!")

	// adding performance data
	p1 := monitoringplugin.NewPerformanceDataPoint("response_time", 10).
		SetUnit("s").SetMin(0)
	p1.NewThresholds(0, 10, 0, 20)
	if err := response.AddPerformanceDataPoint(p1); err != nil {
		// error handling
	}

	p2 := monitoringplugin.NewPerformanceDataPoint("memory_usage", 50.6).
		SetUnit("%").SetMin(0).SetMax(100)
	p2.NewThresholds(0, 80, 0, 90)
	if err := response.AddPerformanceDataPoint(p2); err != nil {
		// error handling
	}

	err = response.AddPerformanceDataPoint(
		monitoringplugin.NewPerformanceDataPoint("memory_usage", 50.6).
			SetUnit("%").SetMin(0).SetMax(100).
			SetThresholds(monitoringplugin.NewThresholds(0, 80.0, 0, 90.0)))
	if err != nil {
		// error handling
	}

	response.OutputAndExit()
	/* exits program with exit code 2 and outputs:
	CRITICAL: something is ok!
	something else is critical!
	something else is warning! | 'response_time'=10s;10;20;0; 'memory_usage'=50%;80;90;0;100
	*/
}
```
