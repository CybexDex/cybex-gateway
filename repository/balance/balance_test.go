package balance

import (
	"fmt"

	m "coding.net/bobxuyang/cy-gateway-BN/models"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOne(t *testing.T) {
	db := m.GetDB()

	err := db.DB().Ping()
	assert.Nil(t, err)

	bal := new(m.Balance)
	err = db.First(bal).Error
	assert.Nil(t, err)
	assert.Equal(t, false, db.NewRecord(bal))
	fmt.Println(bal)
}
