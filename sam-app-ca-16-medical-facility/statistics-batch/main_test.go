package main

import (
	"covid19/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

// 稼働全日が通常のケース
func TestAggregateSubmissions_NormalOnly(t *testing.T) {
	submissions := []models.Submission{
		{
			Answer: "NORMAL",
		},
		{
			Answer: "NORMAL",
		},
		{
			Answer: "NORMAL",
		},
	}

	valid, normal, limitted, stopped, rate := aggregateSubmissions(submissions)

	assert.Equal(t, 3, valid)
	assert.Equal(t, 3, normal)
	assert.Equal(t, 0, limitted)
	assert.Equal(t, 0, stopped)
	assert.Equal(t, 1.0, rate)
}

// 稼働全日が制限のケース
func TestAggregateSubmissions_LimittedOnly(t *testing.T) {
	submissions := []models.Submission{
		{
			Answer: "LIMITTED",
		},
		{
			Answer: "LIMITTED",
		},
		{
			Answer: "LIMITTED",
		},
	}

	valid, normal, limitted, stopped, rate := aggregateSubmissions(submissions)

	assert.Equal(t, 3, valid)
	assert.Equal(t, 0, normal)
	assert.Equal(t, 3, limitted)
	assert.Equal(t, 0, stopped)
	assert.Equal(t, 0.3, rate)
}

// 稼働全日が停止のケース
func TestAggregateSubmissions_StoppedOnly(t *testing.T) {
	submissions := []models.Submission{
		{
			Answer: "STOPPED",
		},
		{
			Answer: "STOPPED",
		},
		{
			Answer: "STOPPED",
		},
	}

	valid, normal, limitted, stopped, rate := aggregateSubmissions(submissions)

	assert.Equal(t, 3, valid)
	assert.Equal(t, 0, normal)
	assert.Equal(t, 0, limitted)
	assert.Equal(t, 3, stopped)
	assert.Equal(t, 0.0, rate)
}

// 複合ケース
// 3つが均等の場合、(1+0.3+0)/3=0.43となる
func TestAggregateSubmissions_Compisite(t *testing.T) {
	submissions := []models.Submission{
		{
			Answer: "NORMAL",
		},
		{
			Answer: "LIMITTED",
		},
		{
			Answer: "STOPPED",
		},
	}

	valid, normal, limitted, stopped, rate := aggregateSubmissions(submissions)

	assert.Equal(t, 3, valid)
	assert.Equal(t, 1, normal)
	assert.Equal(t, 1, limitted)
	assert.Equal(t, 1, stopped)
	assert.Equal(t, 0.43, rate)
}

// 未回答・設置なしは計上されない
func TestAggregateSubmissions_NoValid(t *testing.T) {
	submissions := []models.Submission{
		{
			Answer: "NULL",
		},
		{
			Answer: "NOANSWER",
		},
		{
			Answer: "HOGE",
		},
	}

	valid, normal, limitted, stopped, rate := aggregateSubmissions(submissions)

	assert.Equal(t, 0, valid)
	assert.Equal(t, 0, normal)
	assert.Equal(t, 0, limitted)
	assert.Equal(t, 0, stopped)
	assert.Equal(t, 0.0, rate)
}
