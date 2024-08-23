package services

type Services struct {
	Validator    Validater
	TokenManager TokenManager
	Hash         Hasher
	AuthRepo     AuthRepo
	Alert        Alerter
}

func New(
	v Validater,
	tm TokenManager,
	h Hasher,
	authR AuthRepo,
	a Alerter,
) *Services {
	return &Services{
		Validator:    v,
		TokenManager: tm,
		Hash:         h,
		AuthRepo:     authR,
		Alert:        a,
	}
}
