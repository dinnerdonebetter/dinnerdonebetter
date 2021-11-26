package config

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	BaseConfigSSMKey           = "PRIXFIXE_BASE_CONFIG"
	WriteQueueNameSSMKey       = "PRIXFIXE_WRITES_QUEUE"
	UpdateQueueNameSSMKey      = "PRIXFIXE_UPDATES_QUEUE"
	ArchiveQueueNameSSMKey     = "PRIXFIXE_ARCHIVES_QUEUE"
	DataChangesTopicNameSSMKey = "PRIXFIXE_DATA_CHANGES_QUEUE"
)

// GetConfigFromParameterStore fetches and InstanceConfig from AWS SSM Parameter Store.
func GetConfigFromParameterStore(ctx context.Context) (*InstanceConfig, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := ssm.New(sess)

	input := &ssm.GetParameterInput{Name: aws.String(BaseConfigSSMKey)}
	rawPartialCfgParam, err := svc.GetParameterWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	rawPartialConfig := rawPartialCfgParam.Parameter.String()

	var cfg *InstanceConfig
	if err = json.Unmarshal([]byte(rawPartialConfig), &cfg); err != nil {
		return nil, err
	}

	// fetch database username and password to supplement

	return nil, nil
}
