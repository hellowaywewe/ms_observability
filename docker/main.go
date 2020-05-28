package main

import (
        "log"
        "net/http"
        "strings"

        "github.com/cloudflare/ebpf_exporter/config"
        "github.com/cloudflare/ebpf_exporter/exporter"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
        "github.com/prometheus/common/version"
        "github.com/prometheus/common/expfmt"
        dto "github.com/prometheus/client_model/go"
        kingpin "gopkg.in/alecthomas/kingpin.v2"

        yaml "gopkg.in/yaml.v2"
)

func main() {
        listenAddress := kingpin.Flag("web.listen-address", "The address to listen on for HTTP requests").Default(":9435").String()
        configFile := kingpin.Flag("config.file", "Config file path").Default("config.yaml").File()
        debug := kingpin.Flag("debug", "Enable debug").Bool()
        kingpin.Version(version.Print("ebpf_exporter"))
        kingpin.HelpFlag.Short('h')
        kingpin.Parse()

        var parser expfmt.TextParser
        var parserText = func() ([]*dto.MetricFamily, error) {
            parsed, err := parser.TextToMetricFamilies(strings.NewReader(""))
            if err != nil {
                return nil, err
            }
            var result []*dto.MetricFamily
            for _, mf := range parsed {
                result = append(result, mf)
            }
            return result, nil
        }

        config := config.Config{}

        err := yaml.NewDecoder(*configFile).Decode(&config)
        if err != nil {
            log.Fatalf("Error reading config file: %s", err)
        }

        e := exporter.New(config)
        err = e.Attach()
        if err != nil {
            log.Fatalf("Error attaching exporter: %s", err)
        }

        log.Printf("Starting with %d programs found in the config", len(config.Programs))

        reg := prometheus.NewPedanticRegistry()
        reg.MustRegister(e)

        newGatherers := prometheus.Gatherers{
            reg,
            prometheus.GathererFunc(parserText),
        }
        h := promhttp.HandlerFor(
            newGatherers,
            promhttp.HandlerOpts{
                ErrorLog: nil,
                ErrorHandling: promhttp.ContinueOnError,
            },
        )

        http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request){h.ServeHTTP(w, r) })

        if *debug {
            log.Printf("Debug enabled, exporting raw tables on /tables")
            http.HandleFunc("/tables", e.TablesHandler)
        }

        log.Printf("Listening on %s", *listenAddress)
        err = http.ListenAndServe(*listenAddress, nil)
        if err != nil {
            log.Fatalf("Error listening on %s: %s", *listenAddress, err)
        }
}
