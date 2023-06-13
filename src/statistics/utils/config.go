package utils

type KafkaConfig struct {
	Endpoints []string `json:"endpoints"`
	Topics    []string `json:"topics"`
}

type Configuration struct {
	DB      DBConfiguration `json:"db"`
	LogFile string          `json:"log_file"`
	Port    uint16          `json:"port"`
	Kafka   KafkaConfig     `json:"kafka"`
}

var (
	Config Configuration
)

// TODO: returnable error
func InitConfig() {
	Config = Configuration{
		DB: DBConfiguration{
			"postgres",
			"statistics",
			"program",
			"test",
			"postgres",
			"5432",
		},
		LogFile: "logs/server.log",
		Port:    8030,
		Kafka:   KafkaConfig{Endpoints: []string{"testus-kafka:29092"}, Topics: []string{"statistics"}},
	}
}
