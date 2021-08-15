package infrastructure

import (
	"idealista/domain"
	"idealista/domain/ports"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/wcharczuk/go-chart/v2"
)

type GraphReports struct {
}

func NewReportsService() ports.ReportsService {
	return &GraphReports{}
}

func (rs GraphReports) GetMonthlyRentalReports(smallFlatSlice []domain.Flat, bigFlatSlice []domain.Flat) {
	if len(smallFlatSlice) == 0 || len(bigFlatSlice) == 0 {
		log.Println("Rental flats not found: ", len(smallFlatSlice), len(bigFlatSlice))
		return
	}

	rs.printValuesToFile(smallFlatSlice, bigFlatSlice, ports.RENTAL_REPORT_MONTHLY)
}

func (rs GraphReports) GetMonthlySaleReports(smallFlatSlice []domain.Flat, bigFlatSlice []domain.Flat) {
	if len(smallFlatSlice) == 0 || len(bigFlatSlice) == 0 {
		log.Println("Sale flats not found: ", len(smallFlatSlice), len(bigFlatSlice))
		return
	}

	rs.printValuesToFile(smallFlatSlice, bigFlatSlice, ports.SALE_REPORT_MONTHLY)
}

func (rs GraphReports) printValuesToFile(smallFlatSlice []domain.Flat, bigFlatSlice []domain.Flat, fileName string) {
	var title string
	var stepValue int
	if strings.Contains(fileName, ports.RENTAL_REPORT_MONTHLY) {
		title = "Rent in Granollers"
		stepValue = 50
	} else {
		title = "Sale in Granollers"
		stepValue = 10000

	}
	if len(smallFlatSlice) > len(bigFlatSlice) {
		smallFlatSlice, bigFlatSlice = equalizeFlatSlices(smallFlatSlice, bigFlatSlice, float64(stepValue))
	}

	smallXValuesToPlot, smallYValuesToPlot := getAxisValuesToPlot(smallFlatSlice)
	bigXValuesToPlot, bigYValuesToPlot := getAxisValuesToPlot(bigFlatSlice)
	chartValues := getChartValues(smallFlatSlice)

	sortSlices(smallFlatSlice, bigFlatSlice)

	tickSlice := getYAxisLabels(bigFlatSlice, float64(stepValue))

	graph := chart.Chart{
		Title: title,
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: bigFlatSlice[len(bigFlatSlice)-1].Price,
				Max: bigFlatSlice[0].Price,
			},
			Ticks: tickSlice,
		},
		XAxis: chart.XAxis{
			Ticks: chartValues,
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "75 to 90m2",
				XValues: smallXValuesToPlot,
				YValues: smallYValuesToPlot,
				Style: chart.Style{
					StrokeColor: chart.ColorGreen,
					FillColor:   chart.ColorGreen.WithAlpha(64),
				},
			},
			chart.ContinuousSeries{
				Name:    "90 to 110m2",
				XValues: bigXValuesToPlot,
				YValues: bigYValuesToPlot,
				Style: chart.Style{
					StrokeColor: chart.ColorRed,
					FillColor:   chart.ColorRed.WithAlpha(64),
				},
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	f, _ := os.Create(fileName)
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func equalizeFlatSlices(small []domain.Flat, big []domain.Flat, stepValue float64) ([]domain.Flat, []domain.Flat) {

	sortedSmallFlat := make([]domain.Flat, len(small))
	copy(sortedSmallFlat, small)

	sort.Slice(sortedSmallFlat, func(i, j int) bool {
		return sortedSmallFlat[i].Price > sortedSmallFlat[j].Price
	})

	rest := math.Mod(float64(sortedSmallFlat[len(sortedSmallFlat)-1].Price), stepValue)
	lowestIndex := float64(sortedSmallFlat[len(sortedSmallFlat)-1].Price) - rest
	i := 0
	for len(small) > len(big) {
		flat := small[i]
		big = append([]domain.Flat{*domain.NewFlatWithDate(lowestIndex, 0, flat.Size, flat.Date)}, big...)
		i++
	}

	return small, big
}

func getYAxisLabels(flatSlice []domain.Flat, stepValue float64) []chart.Tick {
	var tickSlice []chart.Tick

	rest := math.Mod(float64(flatSlice[len(flatSlice)-1].Price), stepValue)
	lowestIndex := float64(flatSlice[len(flatSlice)-1].Price) - rest

	tickSlice = append(tickSlice, chart.Tick{Value: lowestIndex, Label: strconv.Itoa(int(lowestIndex))})

	for lowestIndex < flatSlice[0].Price {
		lowestIndex += stepValue
		tickSlice = append(tickSlice, chart.Tick{Value: lowestIndex, Label: strconv.Itoa(int(lowestIndex))})
	}

	return tickSlice
}

func getAxisValuesToPlot(flatSlice []domain.Flat) ([]float64, []float64) {

	var xValueSlice []float64
	var yValueSlice []float64
	for i, v := range flatSlice {
		xValueSlice = append(xValueSlice, float64(i))
		yValueSlice = append(yValueSlice, v.Price)
	}

	return xValueSlice, yValueSlice
}

func getChartValues(flatSlice []domain.Flat) []chart.Tick {

	var chartTick []chart.Tick
	for i, v := range flatSlice {
		chartTick = append(chartTick, chart.Tick{Value: float64(i), Label: v.Date})
	}

	return chartTick
}

func sortSlices(smallFlatSlice []domain.Flat, bigFlatSlice []domain.Flat) {
	sort.Slice(smallFlatSlice, func(i, j int) bool {
		return smallFlatSlice[i].Price > smallFlatSlice[j].Price
	})

	sort.Slice(bigFlatSlice, func(i, j int) bool {
		return bigFlatSlice[i].Price > bigFlatSlice[j].Price
	})
}
