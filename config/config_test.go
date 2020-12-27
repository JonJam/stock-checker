package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestX(t *testing.T) {
	// ARRANGE
	
	c := config.newWithViper()
	

	assert.Equal(t, 123, 123, "they should be equal")
}
