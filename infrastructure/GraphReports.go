package infrastructure

import (
	"idealista/domain"
	"idealista/domain/ports"
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

func (rs GraphReports) GetMonthlyRentalReports(flatSlice []domain.Flat) {
	rs.printValuesToFile(flatSlice, ports.RENTAL_REPORT_MONTHLY)
}

func (rs GraphReports) GetMonthlySaleReports(flatSlice []domain.Flat) {
	rs.printValuesToFile(flatSlice, ports.SALE_REPORT_MONTHLY)
}

func (rs GraphReports) printValuesToFile(flatSlice []domain.Flat, fileName string) {
	valueSlice := getValuesToPlot(flatSlice)
	var tickSlice []chart.Tick

	var title string
	if strings.Contains(fileName, ports.RENTAL_REPORT_MONTHLY) {
		title = "Rental in Granollers"
		tickSlice = getYAxisLabels(flatSlice, 50)
	} else {
		title = "Sale in Granollers"
		tickSlice = getYAxisLabels(flatSlice, 10000)
	}

	graph := chart.BarChart{
		Title: title,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars:     valueSlice,

		YAxis: chart.YAxis{
			Ticks: tickSlice,
		},
	}

	f, _ := os.Create(fileName)
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func getYAxisLabels(flatSlice []domain.Flat, stepValue float64) []chart.Tick {
	var tickSlice []chart.Tick

	sort.Slice(flatSlice, func(i, j int) bool {
		return flatSlice[i].Price > flatSlice[j].Price
	})

	rest := math.Mod(float64(flatSlice[len(flatSlice)-1].Price), stepValue)
	lowestIndex := float64(flatSlice[len(flatSlice)-1].Price) - rest

	tickSlice = append(tickSlice, chart.Tick{Value: lowestIndex, Label: strconv.Itoa(int(lowestIndex))})

	for lowestIndex < flatSlice[0].Price {
		lowestIndex += stepValue
		tickSlice = append(tickSlice, chart.Tick{Value: lowestIndex, Label: strconv.Itoa(int(lowestIndex))})
	}

	return tickSlice
}

func getValuesToPlot(flatSlice []domain.Flat) chart.Values {

	var valueSlice chart.Values
	for i, v := range flatSlice {
		if i%2 == 0 {
			valueSlice = append(valueSlice, chart.Value{Value: float64(v.Price), Label: v.Date})
		}
	}

	return valueSlice
}
