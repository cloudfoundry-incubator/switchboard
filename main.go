package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"

	"code.cloudfoundry.org/lager"

	"github.com/cloudfoundry-incubator/switchboard/api"
	"github.com/cloudfoundry-incubator/switchboard/apiaggregator"
	"github.com/cloudfoundry-incubator/switchboard/config"
	"github.com/cloudfoundry-incubator/switchboard/domain"
	apirunner "github.com/cloudfoundry-incubator/switchboard/runner/api"
	apiaggregatorrunner "github.com/cloudfoundry-incubator/switchboard/runner/apiaggregator"
	"github.com/cloudfoundry-incubator/switchboard/runner/bridge"
	"github.com/cloudfoundry-incubator/switchboard/runner/health"
	"github.com/cloudfoundry-incubator/switchboard/runner/monitor"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/sigmon"
)

func main() {
	rootConfig, err := config.NewConfig(os.Args)

	logger := rootConfig.Logger

	err = rootConfig.Validate()
	if err != nil {
		logger.Fatal("Error validating config:", err, lager.Data{"config": rootConfig})
	}

	if _, err := os.Stat(rootConfig.StaticDir); os.IsNotExist(err) {
		logger.Fatal(fmt.Sprintf("staticDir: %s does not exist", rootConfig.StaticDir), nil)
	}

	backends := domain.NewBackends(rootConfig.Proxy.Backends, logger)

	activeNodeClusterMonitor := monitor.NewClusterMonitor(
		backends,
		rootConfig.Proxy.HealthcheckTimeout(),
		logger.Session("active-monitor"),
		true,
	)

	activeNodeBridgeRunner := bridge.NewRunner(
		rootConfig.Proxy.Port,
		rootConfig.Proxy.ShutdownDelay(),
		logger.Session("active-bridge-runner"),
	)
	clusterStateManager := api.NewClusterAPI(logger)

	activeNodeClusterMonitor.RegisterBackendSubscriber(activeNodeBridgeRunner.ActiveBackendChan)
	activeNodeClusterMonitor.RegisterBackendSubscriber(clusterStateManager.ActiveBackendChan)

	clusterStateManager.RegisterTrafficEnabledChan(activeNodeBridgeRunner.TrafficEnabledChan)
	go clusterStateManager.ListenForActiveBackend()

	apiHandler := api.NewHandler(clusterStateManager, backends, logger, rootConfig.API, rootConfig.StaticDir)
	aggregatorHandler := apiaggregator.NewHandler(logger, rootConfig.API)

	members := grouper.Members{
		{
			Name:   "active-node-bridge",
			Runner: activeNodeBridgeRunner,
		},
		{
			Name:   "api-aggregator",
			Runner: apiaggregatorrunner.NewRunner(rootConfig.API.AggregatorPort, aggregatorHandler),
		},
		{
			Name:   "api",
			Runner: apirunner.NewRunner(rootConfig.API.Port, apiHandler),
		},
		{
			Name:   "active-node-monitor",
			Runner: monitor.NewRunner(activeNodeClusterMonitor, logger),
		},
	}

	if rootConfig.HealthPort != rootConfig.API.Port {
		members = append(members, grouper.Member{
			Name:   "health",
			Runner: health.NewRunner(rootConfig.HealthPort),
		})
	}

	if rootConfig.Proxy.InactiveMysqlPort != 0 {
		inactiveNodeClusterMonitor := monitor.NewClusterMonitor(
			backends,
			rootConfig.Proxy.HealthcheckTimeout(),
			logger.Session("inactive-monitor"),
			false,
		)

		inactiveNodeBridgeRunner := bridge.NewRunner(
			rootConfig.Proxy.InactiveMysqlPort,
			0,
			logger.Session("inactive-bridge-runner"),
		)

		inactiveNodeClusterMonitor.RegisterBackendSubscriber(inactiveNodeBridgeRunner.ActiveBackendChan)
		clusterStateManager.RegisterTrafficEnabledChan(inactiveNodeBridgeRunner.TrafficEnabledChan)

		members = append(members,
			grouper.Member{
				Name:   "inactive-node-bridge",
				Runner: inactiveNodeBridgeRunner,
			},
			grouper.Member{
				Name:   "inactive-node-monitor",
				Runner: monitor.NewRunner(inactiveNodeClusterMonitor, logger),
			},
		)
	}

	group := grouper.NewOrdered(os.Interrupt, members)
	process := ifrit.Invoke(sigmon.New(group))

	logger.Info("Proxy started", lager.Data{"proxyConfig": rootConfig.Proxy})

	err = <-process.Wait()
	if err != nil {
		logger.Fatal("Switchboard exited unexpectedly", err, lager.Data{"proxyConfig": rootConfig.Proxy})
	}
}
