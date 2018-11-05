package mindx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertDpInfluxPoint(t *testing.T) {
	dp := RandomizeDP()
	_, err := ConvertDpInfluxPoint(&dp)
	assert.NoError(t, err)
}

func TestCreateNewBatchPoints(t *testing.T) {
	_, err := CreateNewBatchPoints("test:port")
	assert.NoError(t, err)
}
