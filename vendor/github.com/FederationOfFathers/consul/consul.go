package consul

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

var Client *ConsulClient

type ConsulClient struct {
	*api.Client
	Tags []string
}

var services = map[string]map[string][]*api.CatalogService{}
var servicesLock sync.RWMutex

func Register(serviceName string, portNumber int) error {
	return Client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name: serviceName,
		Port: portNumber,
		Tags: Client.Tags,
	})
}

func RegisterOn(serviceName string, listenOn string) error {
	listenParts := strings.Split(listenOn, ":")
	listenPort, err := strconv.Atoi(listenParts[len(listenParts)-1])
	if err != nil {
		return err
	}
	return Register(serviceName, listenPort)
}

func Service(serviceName, tag string) (string, error) {
	servicesLock.RLock()
	defer servicesLock.RUnlock()
	if _, ok := services[serviceName]; !ok {
		return "", fmt.Errorf("Not watching for service: %s", serviceName)
	}
	servers, ok := services[serviceName][tag]
	if !ok {
		return "", fmt.Errorf("Not watching for service: %s tag: %s", serviceName)
	}
	idx := rand.Intn(len(servers))
	return fmt.Sprintf("%s:%d", servers[idx].Address, servers[idx].ServicePort), nil
}

func WatchService(serviceName, tag string, interval time.Duration) {
	servicesLock.Lock()
	if _, ok := services[serviceName]; !ok {
		services[serviceName] = map[string][]*api.CatalogService{}
	}
	if _, ok := services[serviceName][tag]; !ok {
		services[serviceName][tag] = nil
	}
	servicesLock.Unlock()
	go func(serviceName, tag string, interval time.Duration) {
		t := time.Tick(interval)
		for {
			servers, _, err := Client.Catalog().Service(serviceName, tag, nil)
			if err != nil {
				log.Println("Error watching service: %s, tag: %s: %s", serviceName, tag, err.Error())
				time.Sleep(time.Second)
				continue
			}
			servicesLock.Lock()
			services[serviceName][tag] = servers
			servicesLock.Unlock()
			<-t
		}
	}(serviceName, tag, interval)
}

func Must() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	Client = &ConsulClient{
		Client: client,
	}
}
