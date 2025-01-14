package banner

import (
	"math"
)

type WarpStats struct {
	Config
	CurrentPulls int  // Current number of pulls without 5*
	Lost5050     bool // Whether the last 5★ pull lost the 50/50
}

// NewWarpStats creates a new WarpStats instance for the specified banner type
func NewWarpStats(bannerType Type, currentPulls int, lost5050 bool) WarpStats {
	config := GetConfig(bannerType)
	config.GuaranteedRateUp = lost5050 // Set guaranteed if lost previous 50/50
	return WarpStats{
		Config:       config,
		CurrentPulls: currentPulls,
		Lost5050:     lost5050,
	}
}

// CalculateRateUpProbability calculates chance of getting the rate-up character
func (w *WarpStats) calculateRateUpProbability(probability float64) float64 {
	if w.GuaranteedRateUp {
		return probability // 100% rate-up if guaranteed
	}
	return probability * w.RateUpChance // 50% chance for rate-up otherwise
}

// calculateBaseProbability calculates the probability without pity
func (w *WarpStats) calculateBaseProbability(pulls int) float64 {
	totalRate := w.BaseRate5StarChar + w.BaseRate5StarLC
	notGetting := 1.0 - totalRate
	notGettingInNPulls := math.Pow(notGetting, float64(pulls))
	return 1.0 - notGettingInNPulls
}

// calculateWithPity calculates probability including pity system
func (w *WarpStats) calculateWithPity(pulls int) float64 {
	if pulls >= w.HardPity {
		return 1.0
	}

	probability := 0.0
	remainingPulls := pulls

	if remainingPulls > w.SoftPityStart {
		softPityPulls := remainingPulls - w.SoftPityStart
		remainingPulls = w.SoftPityStart

		baseRate := w.BaseRate5StarChar + w.BaseRate5StarLC // Total 5★ rate
		for i := 0; i < softPityPulls; i++ {
			probability += (1 - probability) * math.Min(1.0, baseRate+float64(i)*0.06)
		}
	}

	baseProbability := w.calculateBaseProbability(remainingPulls)
	probability += (1 - probability) * baseProbability

	return probability
}

// CalculateWarpProbability calculates the probability of getting a 5* based on current pity
func CalculateWarpProbability(bannerType Type, currentPulls, plannedPulls int, lost5050 bool) (float64, float64) {
	stats := NewWarpStats(bannerType, currentPulls, lost5050)
	totalPulls := currentPulls + plannedPulls

	// Calculate base 5★ probability
	baseProbability := stats.calculateWithPity(totalPulls)

	if bannerType == Limited {
		rateUpProb := stats.calculateRateUpProbability(baseProbability)
		return baseProbability, rateUpProb
	}

	// For Standard banner, character probability is exactly half of total
	return baseProbability, baseProbability * 0.5
}
