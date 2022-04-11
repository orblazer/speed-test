package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"github.com/showwin/speedtest-go/speedtest"
)

var (
	namespace  = "speedtest"
	labelsName = []string{"country", "city"}
	latency    = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "latency_ms",
		Help:      "The connetion latency",
	}, labelsName)
	download = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "download_mbps",
		Help:      "The connetion download speed",
	}, labelsName)
	upload = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "upload_mbps",
		Help:      "The connetion upload speed",
	}, labelsName)
)

func recordMetrics() {
	log.Println("Running speed test...")

	user, _ := speedtest.FetchUserInfo()
	serverList, _ := speedtest.FetchServers(user)
	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)

		// Register value to metrics
		latency.WithLabelValues(s.Country, s.Name).Set(float64(s.Latency.Milliseconds()))
		download.WithLabelValues(s.Country, s.Name).Set(s.DLSpeed)
		upload.WithLabelValues(s.Country, s.Name).Set(s.ULSpeed)
	}

	log.Println("Speedtest done !")
}

func main() {
	interval := flag.String("interval", "@hourly", "cron time for run speedtest")
	flag.Parse()

	c := cron.New()
	c.AddFunc(*interval, recordMetrics)
	c.Start()

	// Handle server
	log.Println("Server started on 0.0.0:2112 !")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

	c.Stop()
}
