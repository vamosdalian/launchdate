package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateS3Client(t *testing.T) {
	t.Run("Create client with endpoint", func(t *testing.T) {
		secretID := "test-secret-id"
		secretKey := "test-secret-key"
		endpoint := "https://r2.cloudflarestorage.com/account-id"
		region := "auto"

		client, err := CreateS3Client(secretID, secretKey, region, endpoint)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Create client without endpoint", func(t *testing.T) {
		secretID := "test-secret-id"
		secretKey := "test-secret-key"
		endpoint := ""
		region := "us-west-1"

		client, err := CreateS3Client(secretID, secretKey, region, endpoint)
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})
}
