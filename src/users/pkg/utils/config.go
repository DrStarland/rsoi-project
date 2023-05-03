package utils

import "os"

type oktaConfig struct {
	Endpoint     string `json:"endpoint"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ClientGroup  string `json:"client_group"`
	SSWSToken    string `json:"ssws_token"`
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
		`{"keys":[{"kty":"RSA","alg":"RS256","kid":"XVh9VRM57Bic_gSk2s4owqSKYVQYhZrd7ONvInnszyQ","use":"sig","e":"AQAB","n":"r25i2X_caK8RpM5r4Gugi0N01TGL-rR_3f7vNgkXpL0RlvgJSTWjt8o_NqreLE2b9YLktYI8o7e1Asmz-U2wGA0cepppU5g-7J7B83KyWc8a71Wzj5fSBHr3_SQx2L_sPSQ2lp27fdVjIeL-c2htV69889MGz3ut3snJiGbMNdxfyEbcL8OzjUp1BYkM69A-NBc8xwCwZNWkll5sIHxIb7D5S4m_SnyZ3VTdTTvbhz-8ao3j8ofjWfuC5ys4sTLVXVrWUPKv6yrLBqcLeMezHDaUF-Ocx62dpJTC3-ZhtctmFWhHJdsK0U0VAkHaDX6qGsnpigL0ukskt04rhkr0lQ"}]}`,
		oktaConfig{
			"https://dev-35033098.okta.com",
			"0oa7udguoeLIsv0KJ5d7",
			os.Getenv("OKTA_CLIENT_SECRET"),
			"00g98aome3no6UejV5d7",
			os.Getenv("OKTA_SSWS_TOKEN"),
		},
	}
}
