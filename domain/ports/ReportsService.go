package ports

import "idealista/domain"

const RENTAL_REPORT_MONTHLY = "rental-by-month.gif"
const SALE_REPORT_MONTHLY = "sale-by-month.gif"

type ReportsService interface {
	GetMonthlyRentalReports([][]domain.Flat)
	GetMonthlySaleReports([][]domain.Flat)
}
