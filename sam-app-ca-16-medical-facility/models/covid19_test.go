package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestConvertFacilityType(t *testing.T) {
	result := convertFacilityType("入院")
	assert.Equal(t, "HOSPITAL", result)

	result = convertFacilityType("救急")
	assert.Equal(t, "EMERGENCY", result)

	result = convertFacilityType("外来")
	assert.Equal(t, "OUTPATIENT", result)

	result = convertFacilityType("それ以外")
	assert.Equal(t, "", result)
}

func TestConvertAnsType(t *testing.T) {
	result := convertAnsType("通常")
	assert.Equal(t, "NORMAL", result)

	result = convertAnsType("制限")
	assert.Equal(t, "LIMITTED", result)

	result = convertAnsType("停止")
	assert.Equal(t, "STOPPED", result)

	result = convertAnsType("未回答")
	assert.Equal(t, "NOANSWER", result)

	result = convertAnsType("それ以外")
	assert.Equal(t, "NULL", result)
}
