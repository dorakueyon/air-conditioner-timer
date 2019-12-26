package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/tenntenn/natureremo"
)

func Handler() error {
	// get acccess token of nature remo from parameter store
	sess := session.Must(session.NewSession())
	ssmc := ssm.New(sess)

	keyname := "REMO_ACCESS_TOKEN"
	withDecryption := false
	param, err := ssmc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(keyname),
		WithDecryption: aws.Bool(withDecryption),
	})
	if err != nil {
		return err
	}

	token := *param.Parameter.Value

	// operate nature remo
	cli := natureremo.NewClient(token)
	ctx := context.Background()

	deviceName := "エアコン"
	operation := "on"

	// find aircon device
	as, err := cli.ApplianceService.GetAll(ctx)
	if err != nil {
		return err
	}

	var target *natureremo.Appliance
	for _, a := range as {
		if a.Nickname == deviceName {
			target = a
			break
		}
	}

	if target == nil {
		return errors.New(fmt.Sprintf("%s not found", deviceName))
	}

	// change turn-on/off
	settings := target.AirConSettings
	if operation == "on" {
		settings.Button = ""
	} else {
		settings.Button = "power-off"
	}

	err = cli.ApplianceService.UpdateAirConSettings(ctx, target, settings)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
