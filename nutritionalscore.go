package main

type ScoreType int

const (
	Food ScoreType = iota
	Beverage
	Water
	Cheese
)

var scoreToLetter = []string{"A", "B", "C", "D", "E", "F"}

type NutritionalScore struct {
	Value     int
	Positive  int
	Negative  int
	ScoreType ScoreType
}

func (ns NutritionalScore) GetNutriScore() string {
	scoreType := ns.ScoreType

	if scoreType == Food {
		return scoreToLetter[getPointsFromRange(float64(ns.Value), []float64{18, 10, 2, -1})]
	}

	if scoreType == Water {
		return scoreToLetter[0]
	}

	return scoreToLetter[getPointsFromRange(float64(ns.Value), []float64{9, 5, 1, -2})]
}

type EnergyKJ float64

func (e EnergyKJ) GetPoints(st ScoreType) int {
	if st == Beverage {
		return getPointsFromRange(float64(e), energyLevelsBeverage)
	}

	return getPointsFromRange(float64(e), energyLevels)
}

type SugarGram float64

func (s SugarGram) GetPoints(st ScoreType) int {
	if st == Beverage {
		return getPointsFromRange(float64(s), sugarLevelsBeverage)
	}
	return getPointsFromRange(float64(s), sugarLevels)
}

type SaturatedFattyAcids float64

func (sfa SaturatedFattyAcids) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(sfa), saturatedFattyAcidsLevels)
}

type SodiumMilliGram float64

func (s SodiumMilliGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(s), sodiumLevels)
}

type FruitsPercent float64

func (f FruitsPercent) GetPoints(st ScoreType) int {
	if st == Beverage {
		if f > 80 {
			return 10
		} else if f > 60 {
			return 4
		} else if f > 40 {
			return 2
		}
		return 0
	}

	if f > 80 {
		return 5
	} else if f > 60 {
		return 2
	} else if f > 40 {
		return 1
	}
	return 0
}

type FibreGram float64

func (f FibreGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(f), fiberLevels)
}

type ProteinGram float64

func (p ProteinGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(p), proteinLevels)
}

//convering kcal to kj
func EnergyFromKcal(kcal float64) EnergyKJ {
	return EnergyKJ(kcal * 4.18)
}

//extracting sodium ing mg from salt
func SodiumFroSalt(saltMg float64) SodiumMilliGram {
	return SodiumMilliGram(saltMg / 2.5)
}

type NutritionalData struct {
	Energy              EnergyKJ
	Sugars              SugarGram
	SaturatedFattyAcids SaturatedFattyAcids
	Sodium              SodiumMilliGram
	Fruits              FruitsPercent
	Fibre               FibreGram
	Protein             ProteinGram
	IsWater             bool
}

var (
	energyLevels              = []float64{3350, 3015, 2680, 2345, 2010, 1675, 1340, 1005, 670, 335}
	sugarLevels               = []float64{45, 60, 36, 31, 27, 22.5, 18, 13.5, 9, 4.5}
	saturatedFattyAcidsLevels = []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	sodiumLevels              = []float64{900, 810, 720, 630, 540, 450, 360, 270, 180, 90}

	fiberLevels   = []float64{4.7, 3.7, 2.8, 1.9, 0.9}
	proteinLevels = []float64{8, 6.4, 4.8, 3.2, 1.6}

	energyLevelsBeverage = []float64{270, 240, 210, 180, 150, 120, 90, 60, 30}
	sugarLevelsBeverage  = []float64{13.5, 12, 10.5, 9, 7.5, 6, 4.5, 3, 1.5, 0}
)

func GetNutritionalScore(n NutritionalData, st ScoreType) NutritionalScore {
	var (
		value    int
		positive int
		negative int
	)

	if st != Water {
		fruitPoints := n.Fruits.GetPoints(st)
		fibrePoints := n.Fibre.GetPoints(st)

		negative = n.Energy.GetPoints(st) + n.Sugars.GetPoints(st) + n.Sodium.GetPoints(st) + n.SaturatedFattyAcids.GetPoints(st)
		positive = fruitPoints + fibrePoints + n.Protein.GetPoints(st)

		if st == Cheese {
			value = negative - positive
		} else {
			if negative >= 11 && fruitPoints < 5 {
				value = negative - positive - fruitPoints
			} else {
				value = negative - positive
			}
		}
	}

	return NutritionalScore{
		Value:     value,
		Positive:  positive,
		Negative:  negative,
		ScoreType: st,
	}
}
