package domain_test

import (
	"testing"

	"github.com/sappy5678/cryptocom/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestTransactionID(t *testing.T) {
	transactionID := domain.TransactionID("1234567890")
	assert.Equal(t, "1234567890", transactionID.ID())
	assert.Equal(t, "1234567890-passive", transactionID.PassiveID())
}
