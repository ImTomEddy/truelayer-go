package truelayer

type TrueLayer struct {
	clientID     string
	clientSecret string
}

func New(clientID, clientSecret string) *TrueLayer {
	return &TrueLayer{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}
