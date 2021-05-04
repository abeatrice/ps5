package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mailgun/mailgun-go/v4"
)

type TargetResponse struct {
	Product TargetProduct `json:"product"`
}

type TargetProduct struct {
	AvilableToPromiseNetwork AvilableToPromiseNetwork `json:"available_to_promise_network"`
}

type AvilableToPromiseNetwork struct {
	ProductID                 string  `json:"product_id"`
	AvilableToPromiseQuantity float32 `json:"available_to_promise_quantity"`
}

func main() {
	checkBestBuy()
	checkTarget()
	fmt.Printf("[%s] Done Checking PS5 Stock\n", time.Now().Format("2006-01-02 15:04:05"))
}

func checkBestBuy() {
	url := "https://www.bestbuy.com/site/sony-playstation-5-console/6426149.p?skuId=6426149"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", browser.Chrome())
	res, err := client.Do(req)
	check(err)

	body, err := ioutil.ReadAll(res.Body)
	check(err)

	if strings.Contains(string(body), "Add to Cart") {
		mg := mailgun.NewMailgun("sandbox7098dd5d68634781a39da81e50dfa7de.mailgun.org", os.Getenv("MAILGUN_KEY"))
		message := mg.NewMessage("abeatrice.mail@gmail.com", "PS5 Avilable!", "PS5 stock avilable at bestbuy https://www.bestbuy.com/site/sony-playstation-5-console/6426149.p?skuId=6426149", "abeatrice.mail@gmail.com")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		_, _, err := mg.Send(ctx, message)
		check(err)

		fmt.Println("ps5 avilable at bestbuy")
	}

	res.Body.Close()
}

func checkTarget() {
	// target ps5 product id 81114595
	// junk 80208042
	res, err := http.Get("https://redsky.target.com/v3/pdp/tcin/81114595?excludes=awesome_shop,question_answer_statistics,item,taxonomy,bulk_ship,rating_and_review_reviews,rating_and_review_statistics&key=eb2551e4accc14f38cc42d32fbc2b2ea")
	check(err)

	body, err := ioutil.ReadAll(res.Body)
	check(err)

	var targetPs5 TargetResponse
	err = json.Unmarshal(body, &targetPs5)
	check(err)

	if targetPs5.Product.AvilableToPromiseNetwork.AvilableToPromiseQuantity > 0 {
		mg := mailgun.NewMailgun("sandbox7098dd5d68634781a39da81e50dfa7de.mailgun.org", os.Getenv("MAILGUN_KEY"))
		message := mg.NewMessage("abeatrice.mail@gmail.com", "PS5 Avilable!", "PS5 stock avilable at target: https://www.target.com/p/playstation-5-console/-/A-81114595", "abeatrice.mail@gmail.com")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		_, _, err := mg.Send(ctx, message)
		check(err)

		fmt.Println("ps5 avilable at target")
	}

	res.Body.Close()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
