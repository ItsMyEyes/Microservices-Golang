package configs

import "github.com/joho/godotenv"

func init() {
	setupEnviroment()
}

func setupEnviroment() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}
