package main

import (
	"github.com/Glawary/crypt/cmd"
)

// @title Crypt API
// @description Сервис по получению инфы по криптовалютам с бирж
// @tag.name Crypt
// @tag.description Криптовалюты
// @host localhost:8050
// @BasePath /
// @schemes http
func main() {
	cmd.Run()
}
