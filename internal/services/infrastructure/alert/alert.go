package alert

import "log"

type alertStub struct{}

func New() *alertStub {
	return &alertStub{}
}

func (s *alertStub) Do(email, msg string) error {
	log.Println("Letter was sent!")

	return nil
}
