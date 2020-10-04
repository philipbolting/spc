package npchart

import (
	"math"
	"testing"
	"time"
)

func TestNewNPChartSetsPart(t *testing.T) {
	chart := NewNPChart("Test Part ID", "", "", 0)
	expected := "Test Part ID"
	if chart.part != expected {
		t.Errorf("expected %q but got %q", expected, chart.part)
	}
}

func TestNewNPChartSetsMachine(t *testing.T) {
	chart := NewNPChart("", "Test Machine ID", "", 0)
	expected := "Test Machine ID"
	if chart.machine != expected {
		t.Errorf("expected %q but got %q", expected, chart.machine)
	}
}

func TestNewNPChartSetsCharacteristic(t *testing.T) {
	chart := NewNPChart("", "", "Test Characteristic", 0)
	expected := "Test Characteristic"
	if chart.characteristic != expected {
		t.Errorf("expected %q but got %q", expected, chart.characteristic)
	}
}

func TestNewNPChartSetsSampleSize(t *testing.T) {
	chart := NewNPChart("A", "B", "C", 99)
	expected := uint(99)
	if chart.sampleSize != expected {
		t.Errorf("expected %v but got %v", expected, chart.sampleSize)
	}
}

func TestAddSample(t *testing.T) {
	chart := buildTestChart()
	expected := 22
	if len(chart.samples) != expected {
		t.Errorf("expected %v but got %v", expected, len(chart.samples))
	}
}

func TestGetTotalInspected(t *testing.T) {
	chart := buildTestChart()
	expected := uint(1100)
	got := chart.GetTotalInspected()
	if got != expected {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestGetTotalDefectives(t *testing.T) {
	chart := buildTestChart()
	expected := uint(63)
	got := chart.GetTotalDefectives()
	if got != expected {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestGetNPBar(t *testing.T) {
	chart := buildTestChart()
	expected := float64(2.9)
	got := chart.GetNPBar()
	if math.Abs(got-expected) > 0.05 {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestGetNumberOfDefectivesForSample(t *testing.T) {
	chart := buildTestChart()
	for i, p := range []uint{2, 4, 4, 3, 6, 2, 3, 0, 0, 3, 1, 5, 4, 7, 3, 2, 4, 1, 0, 5, 3, 1} {
		expected := p
		got := chart.GetNumberOfDefectivesForSample(uint(i))
		if got != expected {
			t.Errorf("expected %v but got %v (i=%v)", expected, got, i+1)
		}
	}
}

func TestGetUpperControlLimit(t *testing.T) {
	chart := buildTestChart()
	expected := 7.8
	got := chart.GetUpperControlLimit()
	if math.Abs(got-expected) > 0.05 {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestGetLowerControlLimit(t *testing.T) {
	chart := buildTestChart()
	expected := 0.0
	got := chart.GetLowerControlLimit()
	if math.Abs(got-expected) > 0.05 {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func buildTestChart() *npchart {
	defectives := []uint{2, 4, 4, 3, 6, 2, 3, 0, 0, 3, 1, 5, 4, 7, 3, 2, 4, 1, 0, 5, 3, 1}

	chart := NewNPChart("Test Part ID", "Test Machine ID", "Test Characteristic", 50)

	start := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	for _, d := range defectives {
		chart.AddSample(start, end, d)
		start = end
		end = start.Add(time.Hour)
	}

	return chart
}
