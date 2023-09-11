package reports

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
	"github.com/wcharczuk/go-chart/v2/drawing"
)

const STEP_BETWEEN_FLATS = 2

type GraphReports struct {
}

func NewReportsService() ports.ReportsService {
	return &GraphReports{}
}

func (rs GraphReports) GetMonthlyRentalReports(allFlats [][]domain.Flat) {

	if len(allFlats) == 0 {
		log.Println("Rental flats not found")
		return
	}

	rs.printValuesToFile(allFlats, ports.RENTAL_REPORT_MONTHLY)

}

func (rs GraphReports) GetMonthlySaleReports(allFlats [][]domain.Flat) {
	if len(allFlats) == 0 {
		log.Println("Sale flats not found")
		return
	}

	rs.printValuesToFile(allFlats, ports.SALE_REPORT_MONTHLY)

}

func (rs GraphReports) printValuesToFile(allFlats [][]domain.Flat, fileName string) {
	var title string
	var stepValue int
	if strings.Contains(fileName, ports.RENTAL_REPORT_MONTHLY) {
		title = "Montly Rent"
		stepValue = 50
	} else {
		title = "Monthly Sale"
		stepValue = 20000

	}

	colorSlice := []drawing.Color{chart.ColorGreen, chart.ColorRed, chart.ColorYellow, chart.ColorBlue, chart.ColorOrange, chart.ColorBlack, chart.ColorCyan, chart.ColorAlternateGreen, chart.ColorAlternateBlue, chart.ColorAlternateGray, chart.ColorAlternateYellow}

	var chartSeries []chart.Series
	var joinedFlats []domain.Flat
	for i, flatSlice := range allFlats {
		joinedFlats = append(joinedFlats, flatSlice...)

		xValuesToPlot, yValuesToPlot := getAxisValuesToPlot(flatSlice)
		chartSerie := chart.ContinuousSeries{
			Name:    flatSlice[0].Location + " " + strconv.Itoa(flatSlice[0].Size.GetMinSize()) + " to " + strconv.Itoa(flatSlice[0].Size.GetMaxSize()) + "m²",
			XValues: xValuesToPlot,
			YValues: yValuesToPlot,
			Style: chart.Style{
				StrokeColor: colorSlice[i],
				//FillColor:   colorSlice[i].WithAlpha(64),
			},
		}

		chartSeries = append(chartSeries, chartSerie)
	}

	chartValues := getChartValues(allFlats[0])

	sortSlices(joinedFlats)
	tickSlice := getYAxisLabels(joinedFlats, float64(stepValue))

	graph := chart.Chart{
		Title: title,
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: joinedFlats[len(joinedFlats)-1].Price,
				Max: joinedFlats[0].Price,
			},
			Ticks: tickSlice,
		},
		XAxis: chart.XAxis{
			Ticks: chartValues,
		},
		Series: chartSeries,
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
		big = append([]domain.Flat{*domain.NewFlatWithDate(big[i].Location, lowestIndex, 0, flat.Date, *domain.NewFlatSize(big[i].Size.GetMinSize(), big[i].Size.GetMaxSize()))}, big...)
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

	for i := getFirstIndex(flatSlice); i < len(flatSlice); i = i + STEP_BETWEEN_FLATS {
		xValueSlice = append(xValueSlice, float64(i))
		yValueSlice = append(yValueSlice, flatSlice[i].Price)
	}

	return xValueSlice, yValueSlice
}

func getChartValues(flatSlice []domain.Flat) []chart.Tick {

	var chartTick []chart.Tick

	for i := getFirstIndex(flatSlice); i < len(flatSlice); i = i + STEP_BETWEEN_FLATS {
		chartTick = append(chartTick, chart.Tick{Value: float64(i), Label: flatSlice[i].Date})
	}

	return chartTick
}

func sortSlices(flatSlice []domain.Flat) {
	sort.Slice(flatSlice, func(i, j int) bool {
		return flatSlice[i].Price > flatSlice[j].Price
	})
}

func getFirstIndex(flatSlice []domain.Flat) int {
	if math.Mod(float64(len(flatSlice)), STEP_BETWEEN_FLATS) != 0 {
		return 0
	}

	return 1
}
