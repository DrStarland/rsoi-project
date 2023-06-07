package utils

type DBConfiguration struct {
	Type string `json:"type"`
	Name string `json:"name"`

	User     string `json:"user"`
	Password string `json:"password"`

	Host string `json:"host"`
	Port string `json:"port"`
}

type Configuration struct {
	DBConfig DBConfiguration

	Port                     uint16 `json:"port"`
	RawJWKS                  string `json:"raw-jwks"`
	IdentityProviderEndpoint string `json:"identity-provider-endpoint"`
}

var (
	Config Configuration
)

func InitConfig() {
	Config = Configuration{
		DBConfig: DBConfiguration{
			"postgres",
			"flights",
			"program",
			"test",
			"postgres-service",
			"5432",
		},
		Port: 8050,

		RawJWKS: `{"keys":[{"kid":"oD7q2D3-11tEFQgZXfoikjHVmjcUEPU-iNGirGadNUo","alg":"RS256","e":"AQAB","kty":"RSA","n":"ygo812YXS2SMuX9iJhKZzDFqK0tsyrxkXBbwa1IiMyRIeeznbUYNYnul5WAtf4Kbo-aJxZw10My6rpJk7-bFh-oSB64myR2Gb1rowmd4w621e1Zn4QwMmvhmMYq1LEeXKu4jh2vwZs1ylCoeHfqKgW2qUtDkeXQ2W9aLFByDv1uNDF9oY2PhwrwUdGHlCJt-e4SoPlHBPr0SibMUwr5CfodRfYNOKzPT0hqqRQT6F1FMQZuMOikZY8pw6Q-OriPfcXqeWx68VeU3bmSQ3EPMHd71UDOrzY1dafkKPoLc5qGel4ktuPrrKAn1uiaNeRjN82dLTO0QiAZ5Ly7rGGcM7Q"}]}`,

		IdentityProviderEndpoint: "http://users-service:8040",
	}
}
