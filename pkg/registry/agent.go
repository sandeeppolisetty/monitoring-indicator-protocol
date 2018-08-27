package registry

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"time"

	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

type Agent struct {
	IndicatorsDocument []byte
	RegistryURI        string
	DeploymentName     string
	ProductName        string
	IntervalTime       time.Duration
}

var registrationCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:        "registration_count",
	Help:        "counter of all registration attempts",
	ConstLabels: nil,
}, []string{"status"})

func init() {
	prometheus.MustRegister(registrationCount)
}

func (a Agent) Start() {
	a.registerIndicatorDocument()

	interval := time.NewTicker(a.IntervalTime)
	for {
		select {
		case <-interval.C:
			a.registerIndicatorDocument()
		default:
		}
	}
}

func (a Agent) registerIndicatorDocument() {
	registry := fmt.Sprintf(a.RegistryURI+"/v1/register?deployment=%s&product=%s", a.DeploymentName, a.ProductName)

	body := bytes.NewBuffer(a.IndicatorsDocument)

	resp, err := http.Post(registry, "text/plain", body)
	if err != nil {
		registrationCount.WithLabelValues("err").Inc()
		log.Printf("could not make http request: %s\n", err)
	} else {
		registrationCount.WithLabelValues(strconv.Itoa(resp.StatusCode)).Inc()
	}
}
