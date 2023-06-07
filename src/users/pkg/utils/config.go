package utils

type oktaConfig struct {
	Endpoint     string `json:"endpoint"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ClientGroup  string `json:"client_group"`
	OktetoToken  string `json:"ssws_token"`
}

type Configuration struct {
	LogFile string     `json:"log_file"`
	Port    uint16     `json:"port"`
	RawJWKS string     `json:"raw-jwks"`
	Okta    oktaConfig `json:"okta"`
}

var (
	Config Configuration
)

// TODO: returnable error
func InitConfig() {
	Config = Configuration{
		"logs/server.log",
		8040,
		`{"keys":[{"kid":"oD7q2D3-11tEFQgZXfoikjHVmjcUEPU-iNGirGadNUo","alg":"RS256","e":"AQAB","kty":"RSA","n":"ygo812YXS2SMuX9iJhKZzDFqK0tsyrxkXBbwa1IiMyRIeeznbUYNYnul5WAtf4Kbo-aJxZw10My6rpJk7-bFh-oSB64myR2Gb1rowmd4w621e1Zn4QwMmvhmMYq1LEeXKu4jh2vwZs1ylCoeHfqKgW2qUtDkeXQ2W9aLFByDv1uNDF9oY2PhwrwUdGHlCJt-e4SoPlHBPr0SibMUwr5CfodRfYNOKzPT0hqqRQT6F1FMQZuMOikZY8pw6Q-OriPfcXqeWx68VeU3bmSQ3EPMHd71UDOrzY1dafkKPoLc5qGel4ktuPrrKAn1uiaNeRjN82dLTO0QiAZ5Ly7rGGcM7Q"}]}`,
		oktaConfig{
			Endpoint:     "https://dev-98541142.okta.com",
			ClientId:     "0oa7v8rairOUbYAvy5d7",
			ClientSecret: "iQcihL2DDY6AXyYG3_0XuPsFdWQ9w9vk98xFKBIR",
			ClientGroup:  "",
			OktetoToken:  "YR8wM1sKHc4JQb02kpmN102ZWQN4qVdbxE6J27d5pA2DZd0V", // os.Getenv("OKTETO_TOKEN"),
		},
	}
}
