package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type EnvironmentMessage struct {
	Strategy  string
	Market    string
	Spread    string
	OrderSize string
	Exchanges []string
}

// Ok Handler
func Ok(w http.ResponseWriter, req *http.Request) {
	market := viper.GetString("MARKET")
	strategy := viper.GetString("STRATEGY")
	spread := viper.GetString("SPREAD_PERCENTAGE")
	orderSize := viper.GetString("ORDER_SIZE")
	exchangesEnv := viper.GetString("EXCHANGES")
	exchanges := strings.Split(exchangesEnv, ",")

	ok := EnvironmentMessage{
		Market:    market,
		Spread:    spread,
		Strategy:  strategy,
		OrderSize: orderSize,
		Exchanges: exchanges,
	}

	j, _ := json.Marshal(ok)

	w.Header().Add("Content-Type", "application/json")

	w.Write(j)
}

func main() {
	// Read from environment variables injected in for 12 factor
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "3000")
	port := viper.GetString("PORT")

	log.Printf("Hello, I am alive and running %s", port)

	http.HandleFunc("/", Ok)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
