package models

// ============= DTO ================

type MatchCalculationResult struct {
	MatchPercentage float64
	MatchedProfile  Profile
}

type GetRecommendationRequest struct {
	MinAge      int     `json:"min_age"`
	MaxAge      int     `json:"max_age"`
	MinDistance float64 `json:"min_distance"`
	MaxDistance float64 `json:"max_distance"`
}
