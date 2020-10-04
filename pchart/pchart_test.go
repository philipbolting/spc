package pchart

import (
	"math"
	"testing"
	"time"
)

func TestNewPChartSetsPart(t *testing.T) {
	chart := NewPChart("Test Part ID", "")
	expected := "Test Part ID"
	if chart.part != expected {
		t.Errorf("expected %q but got %q", expected, chart.part)
	}
}

func TestNewPChartSetsMachine(t *testing.T) {
	chart := NewPChart("", "Test Machine ID")
	expected := "Test Machine ID"
	if chart.machine != expected {
		t.Errorf("expected %q but got %q", expected, chart.machine)
	}
}

func TestNewPChartSamplesAreEmpty(t *testing.T) {
	chart := NewPChart("", "")
	expected := 0
	if len(chart.samples) != expected {
		t.Errorf("expected %v but got %v", expected, len(chart.samples))
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
	expected := uint(3624)
	got := chart.GetTotalInspected()
	if got != expected {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestGetTotalDefectives(t *testing.T) {
	chart := buildTestChart()
	expected := uint(243)
	got := chart.GetTotalDefectives()
	if got != expected {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestGetPBar(t *testing.T) {
	chart := buildTestChart()
	expected := float64(0.067)
	got := chart.GetPBar()
	if math.Abs(got-expected) > 0.001 {
		t.Errorf("expected %v but got %v", expected, got)
	}
}

func TestGetProportionDefectiveForSample(t *testing.T) {
	chart := buildTestChart()
	for i, p := range []float64{
		0.07, 0.03, 0.19, 0.15, 0.35, 0.09, 0.04, 0.00, 0.01, 0.03,
		0.05, 0.20, 0.02, 0.00, 0.08, 0.06, 0.10, 0.04, 0.00, 0.08,
		0.08, 0.00} {
		expected := float64(p)
		got := chart.GetProportionDefectiveForSample(uint(i))
		if math.Abs(got-expected) > 0.005 {
			t.Errorf("expected %v but got %v (i=%v)", expected, got, i)
		}
	}
}

func TestGetUpperControlLimitForSample(t *testing.T) {
	chart := buildTestChart()
	for i, ucl := range []float64{
		0.15, 0.19, 0.12, 0.15, 0.12, 0.10, 0.10, 0.18, 0.10, 0.10,
		0.12, 0.26, 0.14, 0.19, 0.16, 0.17, 0.15, 0.16, 0.24, 0.17,
		0.16, 0.14} {
		expected := float64(ucl)
		got := chart.GetUpperControlLimitForSample(uint(i))
		if math.Abs(got-expected) > 0.005 {
			t.Errorf("expected %v but got %v (sample %v)", expected, got, i+1)
		}
	}
}

func TestGetLowerControlLimitForSample(t *testing.T) {
	chart := buildTestChart()
	for i, ucl := range []float64{
		0.00, 0.00, 0.02, 0.00, 0.01, 0.03, 0.04, 0.00, 0.04, 0.03,
		0.01, 0.00, 0.00, 0.00, 0.00, 0.00, 0.00, 0.00, 0.00, 0.00,
		0.00, 0.00} {
		expected := float64(ucl)
		got := chart.GetLowerControlLimitForSample(uint(i))
		if math.Abs(got-expected) > 0.005 {
			t.Errorf("expected %v but got %v (sample %v)", expected, got, i+1)
		}
	}
}

func buildTestChart() *pchart {
	defectives := []uint{6, 1, 41, 13, 60, 40, 22, 0, 3, 14, 9, 3, 2, 0, 5, 3, 8, 3, 0, 4, 6, 0}
	inspected := []uint{92, 36, 212, 86, 172, 448, 564, 48, 594, 530, 188, 15, 97, 36, 65, 54, 82, 67, 18, 52, 72, 96}

	chart := NewPChart("Test Part ID", "Test Machine ID")

	start := time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	for i, d := range defectives {
		chart.AddSample(start, end, inspected[i], d)
		start = end
		end = start.Add(time.Hour)
	}

	return chart
}
