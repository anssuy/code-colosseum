package elo

import "math"

const kFactor = 32

func Expected(ratingA, ratingB int32) float64 {
	return 1 / (1 + math.Pow(10, float64(ratingB-ratingA)/400))
}

func Update(winnerRating, loserRating int32) (newWinnerRating, newLoserRating int32) {
	expectedWinner := Expected(winnerRating, loserRating)
	expectedLoser := Expected(loserRating, winnerRating)

	newWinnerRating = winnerRating + int32(math.Round(kFactor*(1-expectedWinner)))
	newLoserRating = loserRating + int32(math.Round(kFactor*(0-expectedLoser)))

	return newWinnerRating, newLoserRating
}
