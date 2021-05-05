package main

import (
	"flag"
	"log"

	"github.com/poblish/boulevard/generation"
	"golang.org/x/tools/go/packages"
)

type packagesList []string

var packageFlags packagesList

func main() {
	flag.Var(&packageFlags, "pkg", "Packages to scan")
	flag.Parse()

	conf := packages.Config{
		Mode:  packages.NeedName | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
		Tests: false,
	}

	loadedPkgs, err := packages.Load(&conf, packageFlags...)
	if err != nil {
		log.Fatalf("Could not load packages %s", err)
	}

	generator := &generation.DashboardGenerator{}
	metrics := generator.DiscoverMetrics(loadedPkgs)

	err = generator.GenerateAlertRules("alert_rules.yaml", metrics)
	if err != nil {
		log.Fatalf("Alert rule generation failed %s", err)
	}

	err = generator.GenerateGrafanaDashboard("grafana_dashboard.json", metrics)
	if err != nil {
		log.Fatalf("Generation failed %s", err)
	}
}

func (i *packagesList) String() string {
	return "my string representation"
}

func (i *packagesList) Set(value string) error {
	*i = append(*i, value)
	return nil
}
