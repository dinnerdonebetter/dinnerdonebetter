package config

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/prixfixeco/api_server/internal/database"
)

const (
	baseConfigSSMKey            = "PRIXFIXE_BASE_CONFIG"
	databaseConnectionURLSSMKey = "PRIXFIXE_DATABASE_URL"
)

func mustGetParameter(ps *ssm.SSM, paramName string) string {
	input := &ssm.GetParameterInput{Name: aws.String(paramName)}
	rawParam, err := ps.GetParameter(input)
	if err != nil {
		panic(err)
	}

	return rawParam.Parameter.String()
}

// GetConfigFromParameterStore fetches and InstanceConfig from AWS SSM Parameter Store.
func GetConfigFromParameterStore() (*InstanceConfig, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := ssm.New(sess)

	rawPartialConfig := mustGetParameter(svc, baseConfigSSMKey)

	var cfg *InstanceConfig
	if err := json.Unmarshal([]byte(rawPartialConfig), &cfg); err != nil {
		return nil, err
	}

	// fetch supplementary data from SSM
	cfg.Database.ConnectionDetails = database.ConnectionDetails(mustGetParameter(svc, databaseConnectionURLSSMKey))

	return nil, nil
}
