package repositories

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewUserRepository),
	fx.Provide(NewProfileRepository),
	fx.Provide(NewAnswerRepository),
	fx.Provide(NewMatchRepository),
	fx.Provide(NewQuestionRepository),
	fx.Provide(NewRecommendationBinRepository),
)
