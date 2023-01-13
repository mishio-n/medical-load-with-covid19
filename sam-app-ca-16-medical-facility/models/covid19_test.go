package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestConvertFacilityType(t *testing.T) {
	result := ConvertFacilityType("入院")
	assert.Equal(t, "HOSPITAL", result)

	result = ConvertFacilityType("救急")
	assert.Equal(t, "EMERGENCY", result)

	result = ConvertFacilityType("外来")
	assert.Equal(t, "OUTPATIENT", result)

	result = ConvertFacilityType("それ以外")
	assert.Equal(t, "", result)
}

func TestConvertAnsType(t *testing.T) {
	result := ConvertAnsType("通常")
	assert.Equal(t, "NORMAL", result)

	result = ConvertAnsType("制限")
	assert.Equal(t, "LIMITTED", result)

	result = ConvertAnsType("停止")
	assert.Equal(t, "STOPPED", result)

	result = ConvertAnsType("未回答")
	assert.Equal(t, "NOANSWER", result)

	result = ConvertAnsType("それ以外")
	assert.Equal(t, "NULL", result)
}
