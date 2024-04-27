package domain

import (
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
	"time"
)

type readinessTestSuite struct {
	suite.Suite
}

type readinessTestData struct {
	from, to, now string
	pct           float64
}

func TestReadiness(t *testing.T) {
	suite.Run(t, new(readinessTestSuite))
}

func mktime(data string) time.Time {
	r, _ := time.Parse(time.DateOnly, data)
	return r
}

func (s *readinessTestSuite) dataProvider() []readinessTestData {
	return []readinessTestData{
		{"2024-01-01", "2024-01-11", "2024-01-02", 0.1},
		{"2024-01-01", "2024-01-11", "2024-01-03", 0.2},
		{"2024-01-01", "2024-01-11", "2024-01-10", 0.9},
		{"2024-01-01", "2024-01-11", "2024-01-11", 1.0},
	}
}

func (s *readinessTestSuite) TestToEarly() {
	rut := readiness(mktime("2024-01-01"), mktime("2024-01-11"), mktime("2023-12-01"))
	s.Assert().InDelta(0.0, rut, 0.0)
}

func (s *readinessTestSuite) TestIncorrect() {
	rut := readiness(mktime("2024-01-21"), mktime("2024-01-11"), mktime("2023-12-01"))
	s.Assert().InDelta(0.0, rut, 0.0)
}

func (s *readinessTestSuite) TestToLate() {
	rut := readiness(mktime("2024-01-01"), mktime("2024-01-11"), mktime("2024-12-01"))
	s.Assert().InDelta(1.0, rut, 0.0)
}

func (s *readinessTestSuite) TestByProvider() {
	for i, data := range s.dataProvider() {
		s.Run(strconv.Itoa(i), func() {
			rut := readiness(mktime(data.from), mktime(data.to), mktime(data.now))
			s.Assert().InDelta(data.pct, rut, 0.0001)
		})
	}
}
