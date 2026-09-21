package ai

type Service struct {
	APIKey string
}

func NewService(apiKey string) *Service {
	return &Service{
		APIKey: apiKey,
	}
}
