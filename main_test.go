package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-search-data-importer/features/steps"
	"github.com/ONSdigital/dp-search-data-importer/models"
	"github.com/ONSdigital/dp-search-data-importer/schema"
	dplogs "github.com/ONSdigital/log.go/v2/log"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var componentFlag = flag.Bool("component", false, "perform component tests")

type ComponentTest struct {
	t     *testing.T
	Kafka *componenttest.KafkaFeature
}

func (f *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	kafkaScenario := f.Kafka.NewScenario()
	component := steps.NewComponent(f.t, kafkaScenario)

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		if err := component.Reset(); err != nil {
			return ctx, fmt.Errorf("unable to initialise scenario: %s", err)
		}
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		component.Close()
		kafkaScenario.Close(ctx)
		return ctx, nil
	})

	component.RegisterSteps(ctx)
	kafkaScenario.RegisterSteps(ctx)
}

func (f *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	dplogs.Namespace = "dp-search-data-importer"
	f.Kafka = componenttest.NewKafkaFeature(&componenttest.KafkaOptions{
		Encoders: []componenttest.KafkaEncoderOption{
			{
				Topic:    "search-data-import",
				Encoding: "Avro",
				Encoder:  componenttest.NewAvroEncoder[models.SearchDataImport](schema.SearchDataImportEvent),
			},
		},
	})
}

func TestComponent(t *testing.T) {
	if *componentFlag {
		status := 0

		var opts = godog.Options{
			Output:   colors.Colored(os.Stdout),
			Format:   "pretty",
			Paths:    flag.Args(),
			Strict:   true,
			TestingT: t,
		}

		f := &ComponentTest{t: t}

		status = godog.TestSuite{
			Name:                 "feature_tests",
			ScenarioInitializer:  f.InitializeScenario,
			TestSuiteInitializer: f.InitializeTestSuite,
			Options:              &opts,
		}.Run()

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
