package npchart

import (
	"math"
	"time"
)

// Properties of a sample inspected by quality control
type sample struct {
	start      time.Time // production timestamp of the first item in this sample
	end        time.Time // production timestamp of the last item in this sample
	defectives uint      // number of defect items found in this sample
}

// Properties of an np-chart
type npchart struct {
	part                  string   // identifier of the inspected part
	machine               string   // identifier of the machine which processed the inspected parts
	characteristic        string   // description of the inspected characteristic
	sampleSize            uint     // number of inspected items per sample
	samples               []sample // quality control samples inspected for this p-chart
	requiresRecalculation bool     // set to true whenever calculated values must be recalculated
	totalInspected        uint
	totalDefectives       uint
	npBar                 float64
	lcl                   float64 // lower control limit (calculated value)
	ucl                   float64 // upper control limit (calculated value)
}

// NewNPChart initializes a new np-chart for a specific part, machine, characteristic and sampleSize
func NewNPChart(part string, machine string, characteristic string, sampleSize uint) *npchart {
	chart := npchart{part: part, machine: machine, characteristic: characteristic, sampleSize: sampleSize, requiresRecalculation: false}
	return &chart
}

// AddSample adds a sample to an existing np-chart
func (chart *npchart) AddSample(start time.Time, end time.Time, defectives uint) {
	sample := sample{start: start, end: end, defectives: defectives}
	chart.samples = append(chart.samples, sample)
	chart.requiresRecalculation = true
}

func (chart *npchart) GetTotalInspected() uint {
	if chart.requiresRecalculation {
		chart.recalculate()
	}
	return chart.totalInspected
}

func (chart *npchart) GetTotalDefectives() uint {
	if chart.requiresRecalculation {
		chart.recalculate()
	}
	return chart.totalDefectives
}

func (chart *npchart) GetNPBar() float64 {
	if chart.requiresRecalculation {
		chart.recalculate()
	}
	return chart.npBar
}

func (chart *npchart) GetNumberOfDefectivesForSample(i uint) uint {
	if chart.requiresRecalculation {
		chart.recalculate()
	}
	return chart.samples[i].defectives
}

func (chart *npchart) GetUpperControlLimit() float64 {
	if chart.requiresRecalculation {
		chart.recalculate()
	}
	return chart.ucl
}

func (chart *npchart) GetLowerControlLimit() float64 {
	if chart.requiresRecalculation {
		chart.recalculate()
	}
	return chart.lcl
}

func (chart *npchart) recalculate() {
	chart.totalInspected = 0
	chart.totalDefectives = 0
	totalNumberOfSamples := 0
	for _, sample := range chart.samples {
		chart.totalInspected += chart.sampleSize
		chart.totalDefectives += sample.defectives
		totalNumberOfSamples++
	}

	chart.npBar = float64(chart.totalDefectives) / float64(totalNumberOfSamples)

	pBar := chart.npBar / float64(chart.sampleSize)
	deltaCL := 3 * math.Sqrt(chart.npBar*(1-pBar))
	chart.ucl = chart.npBar + deltaCL
	chart.lcl = chart.npBar - deltaCL

	if chart.lcl < 0.0 {
		chart.lcl = 0.0
	}
}
