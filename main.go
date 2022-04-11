package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"github.com/showwin/speedtest-go/speedtest"
)

var (
	labelsName = []string{"country", "city"}
	latency    = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "speedtest_latency_ms",
		Help: "The connetion latency",
	}, labelsName)
	download = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "speedtest_download_mbps",
		Help: "The connetion download speed",
	}, labelsName)
	upload = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "speedtest_upload_mbps",
		Help: "The connetion upload speed",
	}, labelsName)
)

func recordMetrics() {
	fmt.Println("Running speed test...")

	user, _ := speedtest.FetchUserInfo()
	serverList, _ := speedtest.FetchServers(user)
	targets, _ := serverList.FindServer([]int{})

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)

		// Register value to metrics
		latency.WithLabelValues(s.Country, s.Name).Observe(float64(s.Latency.Milliseconds()))
		download.WithLabelValues(s.Country, s.Name).Observe(s.DLSpeed)
		upload.WithLabelValues(s.Country, s.Name).Observe(s.ULSpeed)
	}

	fmt.Println("Speedtest done !")
}

func main() {
	interval := flag.String("interval", "@hourly", "cron time for run speedtest")
	flag.Parse()

	c := cron.New()
	c.AddFunc(*interval, recordMetrics)
	c.Start()

	// Handle server
	fmt.Println("Server started on 0.0.0:2112 !")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

	c.Stop()
}
