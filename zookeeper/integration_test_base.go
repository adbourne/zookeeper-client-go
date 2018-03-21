package zookeeper

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"gopkg.in/ory-am/dockertest.v3"
	"strconv"
	"time"
)

const (
	// Zookeeper
	zookeeperDockerRepo = "zookeeper"
	zookeeperDockerTag = "3.4.10"
)

type BaseIntegrationTestSuite struct {
	suite.Suite
	Logger LoggerFunc

	DockerTestPool *dockertest.Pool

	zookeeperResource *dockertest.Resource
}

func (suite *BaseIntegrationTestSuite) InitialiseBase() {
	suite.Logger = NewStdOutLogger()

	suite.DockerTestPool = suite.createDockerConnectionPoolOrFail()
}

type ConnectionDetails struct {
	Host string
	Port int
}

func (cd *ConnectionDetails) HostPort() string {
	return fmt.Sprintf("%s:%d", cd.Host, cd.Port)
}

func (suite *BaseIntegrationTestSuite) StartZookeeper() *ConnectionDetails {
	zookeeperOptions := &dockertest.RunOptions{
		Repository: zookeeperDockerRepo,
		Tag:        zookeeperDockerTag,
	}

	suite.Logger.Log("message", "Starting Zookeeper...")
	resource := suite.startDockerTestContainerOrFail(suite.DockerTestPool, zookeeperOptions)

	suite.zookeeperResource = resource

	host := "localhost"
	port, _ := strconv.Atoi(resource.GetPort("2181/tcp"))

	connectionDetails := &ConnectionDetails{
		Host:       host,
		Port:       port,
	}

	suite.waitForZookeeperToStartUp(connectionDetails)

	return connectionDetails
}

func (suite *BaseIntegrationTestSuite) waitForZookeeperToStartUp(connectionDetails *ConnectionDetails) {
	// Connect to zookeeper
	hostPort := connectionDetails.HostPort()

	outerErr := suite.DockerTestPool.Retry(func() (err error) {
		suite.Logger.Log(
			"message", "Attempting to connect to Zookeeper...",
			"hostPort", hostPort,
		)

		time.Sleep(3 * time.Second) // TODO: This is gross, change ASAP

		return
	})

	if outerErr != nil {
		suite.FailNow("Could not connect to docker: %s", outerErr)
	}

	suite.Logger.Log("message", "Zookeeper started successfully", "host", connectionDetails.Host)
	return
}

func (suite *BaseIntegrationTestSuite) StopZookeeper() {
	if suite.zookeeperResource != nil {
		err := suite.DockerTestPool.Purge(suite.zookeeperResource)
		if err != nil {
			suite.Fail(fmt.Sprintf("Could not stop zookeeper: %s", err))
		}
	}
}

// Dockertest

func (suite *BaseIntegrationTestSuite) createDockerConnectionPoolOrFail() (pool *dockertest.Pool) {
	suite.Logger.Log("message", "Connecting to docker...")
	pool, err := dockertest.NewPool("")
	if err != nil {
		suite.Fail(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	suite.Logger.Log("message", "Connected to docker!")
	return pool
}

func (suite *BaseIntegrationTestSuite) startDockerTestContainerOrFail(pool *dockertest.Pool, options *dockertest.RunOptions) (resource *dockertest.Resource) {
	resource, err := pool.RunWithOptions(options)
	if err != nil {
		suite.Fail(fmt.Sprintf("Could not start resource %s", err))
	}
	return
}
