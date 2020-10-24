package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

type Credentials struct {
	Name       string
	Key        string
	Passphrase string
	Secret     string
}

type EnvironmentMessage struct {
	Strategy  string
	Market    string
	Spread    string
	OrderSize string
	Exchanges []string
	ApiKeys   []Credentials
}

// Ok Handler
func Ok(w http.ResponseWriter, req *http.Request) {
	market := viper.GetString("MARKET")
	strategy := viper.GetString("STRATEGY")
	spread := viper.GetString("SPREAD_PERCENTAGE")
	orderSize := viper.GetString("ORDER_SIZE")
	exchangesEnv := viper.GetString("EXCHANGES")
	exchanges := strings.Split(exchangesEnv, ",")

	creds := make([]Credentials, 0)

	for _, e := range exchanges {
		getCredential := func(suffix string) string {
			return viper.GetString(strings.ToUpper(fmt.Sprintf("%s_%s", e, suffix)))
		}

		key := getCredential("key")
		passphrase := getCredential("passphrase")
		secret := getCredential("secret")

		// Redact for privacy
		if len(key) > 5 {
			key = fmt.Sprintf("%s********", key[0:5])
		}

		// Redact for privacy
		if len(passphrase) > 5 {
			passphrase = fmt.Sprintf("%s********", passphrase[0:5])
		}

		// Redact for privacy
		if len(secret) > 5 {
			secret = fmt.Sprintf("%s********", secret[0:5])
		}

		creds = append(creds, Credentials{
			Name:       e,
			Key:        key,
			Passphrase: passphrase,
			Secret:     secret,
		})
	}

	ok := EnvironmentMessage{
		Market:    market,
		Spread:    spread,
		Strategy:  strategy,
		OrderSize: orderSize,
		Exchanges: exchanges,
		ApiKeys:   creds,
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
