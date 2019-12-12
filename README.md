# go-monitoringplugin

## Description
Golang package for writing monitoring check plugins for nagios, icinga2, zabix, etc. It provides the check plugin logic specified at https://www.monitoring-plugins.org/doc/guidelines.html

## Example / Usage
	package main

	import (
		monitoringplugin "github.com/inexio/go-monitoringplugin"
	)

	func main() {
		//Creating response with a default ok message that will be displayed when the checks exits with status ok
		response := monitoringplugin.NewResponse("everything checked!")

		//Set output delimiter (default is \n)
		//response.SetOutputDelimiter(" / ")

		//updating check plugin status and adding message to the ouput (status only changes if the new status is worse than the current one)
		response.UpdateStatus(monitoringplugin.OK, "something is ok!") //check status stays ok
		response.UpdateStatus(monitoringplugin.CRITICAL, "something else is critical!") //check status updates to critical
		response.UpdateStatus(monitoringplugin.WARNING, "something else is warning!") //check status stays critical, but message will be added to the output


		//adding performance data
		err := response.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("response_time", 10, "s").SetWarn(10).SetCrit(20).SetMin(0))
		if err != nil {
			//error handling
		}
		err = response.AddPerformanceDataPoint(monitoringplugin.NewPerformanceDataPoint("memory_usage", 50, "%").SetWarn(80).SetCrit(90).SetMin(0).SetMax(100))
		if err != nil {
			//error handling
		}

		response.OutputAndExit()
		/* exits program with exit code 2 and outputs:
		CRITICAL: something is ok!
		something else is critical!
		something else is warning! | 'response_time'=10s;10;20;0; 'memory_usage'=50%;80;90;0;100
		*/
	}