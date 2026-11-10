package service

type DeliveryPricer interface {
	Price(weightGrams, distanceKM int) (int64, error)
}

type BaseTariff struct {
	BasePriceKopecks int64
	PerKMKopecks     int64
	MaxDistanceKM    int
}

func (b BaseTariff) validate(weightGrams, distanceKM int) error {
	// TODO: проверьте вес, положительную дистанцию и максимальную дистанцию.
	if weightGrams <= 0 {
		return ErrInvalidWeight
	}
	if distanceKM <= 0 {
		return ErrInvalidDistance
	}
	if distanceKM >= b.MaxDistanceKM {
		return ErrDistanceTooLarge
	}
	return nil
}

func billableKilograms(weightGrams int) int64 {
	// TODO: округлите граммы вверх до целого оплачиваемого килограмма.
	if weightGrams <= 0 {
		return 0
	}
	return int64(1 + (weightGrams-1)/1000)
}

type Courier struct {
	BaseTariff
	PerKGKopecks int64
}

func NewCourier() Courier {
	return Courier{BaseTariff: BaseTariff{20_000, 1_500, 100}, PerKGKopecks: 5_000}
}

func (t Courier) Price(weightGrams, distanceKM int) (int64, error) {
	// TODO: реализуйте формулу тарифа.
	err := t.validate(weightGrams, distanceKM)

	if err != nil {
		return 0, err
	}

	return billableKilograms(weightGrams)*t.PerKGKopecks +
		t.BaseTariff.PerKMKopecks*int64(distanceKM) +
		t.BaseTariff.BasePriceKopecks, nil
}

type Pickup struct {
	BaseTariff
	PerKGKopecks int64
}

func NewPickup() Pickup {
	return Pickup{BaseTariff: BaseTariff{10_000, 500, 50}, PerKGKopecks: 2_000}
}

func (t Pickup) Price(weightGrams, distanceKM int) (int64, error) {
	// TODO: реализуйте формулу тарифа.
	err := t.validate(weightGrams, distanceKM)

	if err != nil {
		return 0, err
	}

	return billableKilograms(weightGrams)*t.PerKGKopecks +
		t.BaseTariff.PerKMKopecks*int64(distanceKM) +
		t.BaseTariff.BasePriceKopecks, nil
}

type Express struct {
	BaseTariff
	PerKGKopecks int64
}

func NewExpress() Express {
	return Express{BaseTariff: BaseTariff{50_000, 3_000, 300}, PerKGKopecks: 10_000}
}

func (t Express) Price(weightGrams, distanceKM int) (int64, error) {
	// TODO: реализуйте формулу тарифа.
	err := t.validate(weightGrams, distanceKM)

	if err != nil {
		return 0, err
	}

	return billableKilograms(weightGrams)*t.PerKGKopecks +
		t.BaseTariff.PerKMKopecks*int64(distanceKM) +
		t.BaseTariff.BasePriceKopecks, nil
}

var (
	_ DeliveryPricer = Courier{}
	_ DeliveryPricer = Pickup{}
	_ DeliveryPricer = Express{}
)
