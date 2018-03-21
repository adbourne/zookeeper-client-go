package zookeeper

import (
	"testing"
	"github.com/stretchr/testify/suite"
)

type NativeClientTestSuite struct {
	BaseIntegrationTestSuite
	connectionDetails *ConnectionDetails

	underTest *NativeClient
}

func (suite *NativeClientTestSuite) SetupSuite() {
	suite.InitialiseBase()
	suite.connectionDetails = suite.StartZookeeper()

	suite.underTest = &NativeClient{}
}

func (suite *NativeClientTestSuite) TearDownSuite() {
	suite.StopZookeeper()
}


func TestRunClientTestSuite(t *testing.T) {
	suite.Run(t, new(NativeClientTestSuite))
}
