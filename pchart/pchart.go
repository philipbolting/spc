// Package pchart provides a p-chart implementation for statistical process control (SPC)
package pchart

import (
	"math"
	"time"
)

// Properties of a sample inspected by quality control
type sample struct {
	start       time.Time // production timestamp of the first item in this sample
	end         time.Time // production timestamp of the last item in this sample
	inspections uint      // number of inspected items in this sample
	defectives  uint      // number of defect items found in this sample
	p           float64   // proportion defective (calculated value)
	lcl         float64   // lower control limit (calculated value)
	ucl         float64   // upper control limit (calculated value)
}

// Properties of a p-chart
type pchart struct {
	part                  string   // identifier of the inspected part
	machine               string   // identifier of the machine which processed the inspected parts
	samples               []sample // quality control samples inspected for this p-chart
	requiresRecalculation bool     // set to true whenever calculated values must be recalculated
	totalInspected        uint
	totalDefectives       uint
	pBar                  float64
}

// NewPChart initializes a new p-chart for a specific part and machine
func NewPChart(part string, machine string) *pchart {
	chart := pchart{part: part, machine: machine, samples: []sample{}, requiresRecalculation: false}
	return &chart
}

// AddSample adds a sample to an existing p-chart
func (chart *pchart) AddSample(start time.Time, end time.Time, inspections uint, defectives uint) {
	sample := sample{start: start, end: end, inspections: inspections, defectives: defectives}
	chart.samples = append(chart.samples, sample)
	chart.requiresRecalculation = true
}

func (chart *pchart) GetTotalInspected() uint {
	if chart.requiresRecalculation {
		chart.reCalculate()
	}
	return chart.totalInspected
}

func (chart *pchart) GetTotalDefectives() uint {
	if chart.requiresRecalculation {
		chart.reCalculate()
	}
	return chart.totalDefectives
}

func (chart *pchart) GetPBar() float64 {
	if chart.requiresRecalculation {
		chart.reCalculate()
	}
	return chart.pBar
}

func (chart *pchart) GetProportionDefectiveForSample(i uint) float64 {
	if chart.requiresRecalculation {
		chart.reCalculate()
	}
	return chart.samples[i].p
}

func (chart *pchart) GetUpperControlLimitForSample(i uint) float64 {
	if chart.requiresRecalculation {
		chart.reCalculate()
	}
	return chart.samples[i].ucl
}

func (chart *pchart) GetLowerControlLimitForSample(i uint) float64 {
	if chart.requiresRecalculation {
		chart.reCalculate()
	}
	return chart.samples[i].lcl
}

func (chart *pchart) reCalculate() {
	chart.totalInspected = 0
	chart.totalDefectives = 0
	for _, sample := range chart.samples {
		chart.totalInspected += sample.inspections
		chart.totalDefectives += sample.defectives
	}

	chart.pBar = float64(chart.totalDefectives) / float64(chart.totalInspected)
	// round pBar to three decimals
	chart.pBar = math.Round(1000*chart.pBar) / 1000

	deltaCL := 0.0
	for i, sample := range chart.samples {
		chart.samples[i].p = float64(sample.defectives) / float64(sample.inspections)
		deltaCL = 3 * math.Sqrt(chart.pBar*float64(1-chart.pBar)/float64(sample.inspections))
		chart.samples[i].ucl = chart.pBar + deltaCL
		chart.samples[i].lcl = chart.pBar - deltaCL

		// round p, ucl and lcl to two decimals
		chart.samples[i].p = math.Round(100*chart.samples[i].p) / 100
		chart.samples[i].ucl = math.Round(100*chart.samples[i].ucl) / 100
		chart.samples[i].lcl = math.Round(100*chart.samples[i].lcl) / 100

		if chart.samples[i].lcl < 0.0 {
			chart.samples[i].lcl = 0.0
		}
	}

	chart.requiresRecalculation = false
}
