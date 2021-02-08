package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type stonkResp struct {
	Open float32 `json:"open"`

	Close float32 `json:"close"`
}

func getStockPrice(symbol string, apiKey string) (openClose stonkResp, theError error) {
	yesterday := time.Now().AddDate(0, 0, -1)
	yesterdayFormatted := yesterday.Format("2006-01-02")
	resp, err := http.Get("https://api.polygon.io/v1/open-close/" + symbol + "/" + yesterdayFormatted + "?apiKey=" + apiKey)

	if err != nil {
		theError = err
	}

	fmt.Println(resp.StatusCode)

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&openClose)

	return
}

func getSecret(theConfig aws.Config, secretName string, secretKey string) (key string, err string) {
	svc := secretsmanager.NewFromConfig(theConfig)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, myErr := svc.GetSecretValue(context.TODO(), input)

	if myErr != nil {
		err = myErr.Error()
		fmt.Println(myErr)
		return
	}

	jsonString := *result.SecretString
	var finalResult map[string]interface{}
	json.Unmarshal([]byte(jsonString), &finalResult)
	key = finalResult[secretKey].(string)
	return
}

func main() {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithDefaultRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Running job")
	apiKey, _ := getSecret(config, "prod/polygonApiKey", "polygonApiKey")
	client := sns.NewFromConfig(config)
	msgToSend := sns.PublishInput{
		TopicArn: aws.String("arn:aws:sns:us-east-1:308841612984:StonksUpdate"),
	}

	stock := "GME"
	resp, _ := getStockPrice(stock, apiKey)

	priceDiff := resp.Open - resp.Close

	msg := "Here's your " + stock + " stock update...\n"

	if priceDiff < 0 {
		msg = msg + "Well... Even Apes get lucky sometimes"
	} else if priceDiff < 20 {
		msg = msg + "Meh, not too bad."
	} else if priceDiff < 60 {
		msg = msg + "You should probably stick to lottery tickets"
	} else {
		msg = msg + "https://www.titlemax.com"
	}
	msgToSend.Message = aws.String(msg)
	client.Publish(context.TODO(), &msgToSend)
}
