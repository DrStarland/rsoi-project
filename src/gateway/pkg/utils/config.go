package utils

type Configuration struct {
	Port          uint16 `json:"port"`
	RawJWKS       string `json:"raw-jwks"`
	UsersEndpoint string
	NotesEndpoint string
	TasksEndpoint string
	CostsEndpoint string
}

var (
	Config Configuration
)

// TODO: returnable error
func InitConfig() {
	Config = Configuration{
		Port:    8080,
		RawJWKS: `{"keys":[{"kid":"oD7q2D3-11tEFQgZXfoikjHVmjcUEPU-iNGirGadNUo","alg":"RS256","e":"AQAB","kty":"RSA","n":"ygo812YXS2SMuX9iJhKZzDFqK0tsyrxkXBbwa1IiMyRIeeznbUYNYnul5WAtf4Kbo-aJxZw10My6rpJk7-bFh-oSB64myR2Gb1rowmd4w621e1Zn4QwMmvhmMYq1LEeXKu4jh2vwZs1ylCoeHfqKgW2qUtDkeXQ2W9aLFByDv1uNDF9oY2PhwrwUdGHlCJt-e4SoPlHBPr0SibMUwr5CfodRfYNOKzPT0hqqRQT6F1FMQZuMOikZY8pw6Q-OriPfcXqeWx68VeU3bmSQ3EPMHd71UDOrzY1dafkKPoLc5qGel4ktuPrrKAn1uiaNeRjN82dLTO0QiAZ5Ly7rGGcM7Q"}]}`,

		UsersEndpoint: "http://users-service:8040",
		NotesEndpoint: "http://notes-service:8050",
		TasksEndpoint: "http://tasks-service:8060",
		CostsEndpoint: "http://costs-service:8070",
	}
}
