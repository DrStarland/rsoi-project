package utils

type Configuration struct {
	LogFile                  string `json:"log_file"`
	Port                     uint16 `json:"port"`
	RawJWKS                  string `json:"raw-jwks"`
	IdentityProviderEndpoint string `json:"identity-provider-endpoint"`
	FlightsEndpoint          string `json:"flights-endpoint"`
	TicketsEndpoint          string `json:"tickets-endpoint"`
	PrivilegesEndpoint       string `json:"privileges-endpoint"`
}

var (
	Config Configuration
)

// TODO: returnable error
func InitConfig() {
	Config = Configuration{
		"logs/server.log",
		8080,
		`{"keys":[{"kid":"oD7q2D3-11tEFQgZXfoikjHVmjcUEPU-iNGirGadNUo","alg":"RS256","e":"AQAB","kty":"RSA","n":"ygo812YXS2SMuX9iJhKZzDFqK0tsyrxkXBbwa1IiMyRIeeznbUYNYnul5WAtf4Kbo-aJxZw10My6rpJk7-bFh-oSB64myR2Gb1rowmd4w621e1Zn4QwMmvhmMYq1LEeXKu4jh2vwZs1ylCoeHfqKgW2qUtDkeXQ2W9aLFByDv1uNDF9oY2PhwrwUdGHlCJt-e4SoPlHBPr0SibMUwr5CfodRfYNOKzPT0hqqRQT6F1FMQZuMOikZY8pw6Q-OriPfcXqeWx68VeU3bmSQ3EPMHd71UDOrzY1dafkKPoLc5qGel4ktuPrrKAn1uiaNeRjN82dLTO0QiAZ5Ly7rGGcM7Q"}]}`,

		"http://users-service:8040",

		"http://flights-service:8060",
		"http://tickets-service:8070",
		"http://privileges-service:8050",
	}
}
