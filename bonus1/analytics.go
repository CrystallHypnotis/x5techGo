package analytics

type Sale struct {
	Product          string
	Category         string
	Quantity         int
	UnitPriceKopecks int64
	Cancelled        bool
}

type Report struct {
	RevenueKopecks    int64
	UnitsSold         int
	RevenueByCategory map[string]int64
	TopProduct        string
	InvalidSales      int
}

func AnalyzeSales(sales []Sale) Report {
	// TODO: реализуйте расчёт отчёта по правилам задания.
	rep := Report{
		RevenueByCategory: make(map[string]int64),
	}
	if len(sales) == 0 {
		return rep
	}

	best_unit, best_units_sold := "\xff", 0 //немного чёрной магии, самый большой символ в UTF-8
	m := make(map[string]int)

	for _, sale := range sales {
		if sale.Cancelled {
			continue
		}

		if sale.Product == "" || sale.Category == "" ||
			sale.Quantity <= 0 || sale.UnitPriceKopecks <= 0 {
			rep.InvalidSales++
			continue
		}

		revenue := sale.UnitPriceKopecks * int64(sale.Quantity)

		rep.RevenueKopecks += revenue
		rep.UnitsSold += sale.Quantity
		rep.RevenueByCategory[sale.Category] += revenue
		if _, ok := m[sale.Product]; !ok {
			m[sale.Product] = sale.Quantity
		} else {
			m[sale.Product] += sale.Quantity
		}

		if (m[sale.Product] > best_units_sold) || (sale.Product < best_unit &&
			m[sale.Product] == best_units_sold) {

			best_unit = sale.Product
			best_units_sold = m[sale.Product]
		}
	}

	rep.TopProduct = best_unit
	return rep
}
