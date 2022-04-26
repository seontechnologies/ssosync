package config

import (
	"encoding/base64"

	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// Secrets ...
type Secrets struct {
	svc      *secretsmanager.SecretsManager
	secretID string
}

// NewSecrets ...
func NewSecrets(svc *secretsmanager.SecretsManager, secretID string) *Secrets {
	return &Secrets{
		svc:      svc,
		secretID: secretID,
	}
}

// GoogleAdminEmail ...
func (s *Secrets) GoogleAdminEmail() (string, error) {
	return s.getSecret("SSOSyncGoogleAdminEmail")
}

// SCIMAccessToken ...
func (s *Secrets) SCIMAccessToken() (string, error) {
	return s.getSecret("SSOSyncSCIMAccessToken")
}

// SCIMEndpointUrl ...
func (s *Secrets) SCIMEndpointUrl() (string, error) {
	return s.getSecret("SSOSyncSCIMEndpointUrl")
}

// GoogleCredentials ...
func (s *Secrets) GoogleCredentials() (string, error) {
	return s.getSecret("SSOSyncGoogleCredentials")
}

func (s *Secrets) getSecret(secretKey string) (string, error) {
	r, err := s.svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(s.secretID),
		VersionStage: aws.String("AWSCURRENT"),
	})

	if err != nil {
		return "", err
	}

	var secretString string
	var secret map[string]string

	if r.SecretString != nil {
		secretString = *r.SecretString
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(r.SecretBinary)))
		l, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, r.SecretBinary)
		if err != nil {
			return "", err
		}
		secretString = string(decodedBinarySecretBytes[:l])
	}

	if err := json.Unmarshal([]byte(secretString), &secret); err != nil {
		return "", err
	}

	if v, ok := secret[secretKey]; ok {
		return v, nil
	} else {
		return "", err
	}

}
