package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/maxim-nazarenko/tf-module-update/internal/conditions"
	"github.com/maxim-nazarenko/tf-module-update/internal/module"
	"github.com/maxim-nazarenko/tf-module-update/internal/processing"
	"github.com/maxim-nazarenko/tf-module-update/internal/processing/logging"
	"github.com/maxim-nazarenko/tf-module-update/internal/strategies"
)

type AppConfig struct {
	Write      bool
	LogLevel   logging.Level
	Paths      []string
	FromSource module.Source
	ToSource   module.Source
}

// =======================================================

func main() {

	config, err := parseFlags()
	if err != nil {
		log.Fatal(err.Error())
	}

	if len(config.Paths) < 1 {
		config.Paths = []string{"."}
	}

	os.Exit(run(config))
}

func run(config *AppConfig) int {
	results := processing.NewResults(config.LogLevel)
	defer func() {
		fmt.Println(results.String())
	}()

	updateCondition, err := conditionFromSource(config.FromSource)
	if err != nil {
		results.Append(err)
		return 1
	}
	results.Append(processing.NewResultFactory().Debug("searching for module sources: " + config.FromSource.String()))
	results.Append(processing.NewResultFactory().Debug("updating source with: " + config.ToSource.String()))
	strategy := strategies.NewStrictUpdater(
		func(s module.Source) module.Source {
			return s.Merge(config.ToSource)
		}).
		WithCondition(updateCondition)

	processing.NewManager(processing.Config{
		Write:            config.Write,
		ExcludeItemsFunc: processing.DefaultExclusionFunc,
	}, strategy).
		ProcessPaths(config.Paths, results)

	if results.HasErrors() {
		return 1
	}

	return 0
}

func parseFlags() (*AppConfig, error) {
	config := AppConfig{}

	flag.BoolVar(&config.Write, "write", false, "Write files. Defaults to printing possible updates without actual writing")

	var logLevel string
	flag.StringVar(&logLevel, "log.level", "info", "One of trace, debug, info, warn, error")

	var fromURL string
	var toURL string
	flag.StringVar(&fromURL, "from.url", "", "The full URL of module to find and update")
	flag.StringVar(&toURL, "to.url", "", "The full URL to update the matching module sources to")

	var fromScheme string
	var toScheme string
	flag.StringVar(&fromScheme, "from.scheme", "", "Filter modules to update by this scheme")
	flag.StringVar(&toScheme, "to.scheme", "", "Update matching modules with this new scheme")

	var fromHost string
	var toHost string
	flag.StringVar(&fromHost, "from.host", "", "Filter modules to update by this host")
	flag.StringVar(&toHost, "to.host", "", "Update matching modules with this new host")

	var fromModule string
	var toModule string
	flag.StringVar(&fromModule, "from.module", "", "Filter modules to update by this module")
	flag.StringVar(&toModule, "to.module", "", "Update matching modules with this new module")

	var fromSubmodule string
	var toSubmodule string
	flag.StringVar(&fromSubmodule, "from.submodule", "", "Filter modules to update by this submodule")
	flag.StringVar(&toSubmodule, "to.submodule", "", "Update matching modules with this new submodule")

	var fromRevisionStr string
	var toRevisionStr string
	flag.StringVar(&fromRevisionStr, "from.revision", "", "Filter modules to update by this revision")
	flag.StringVar(&toRevisionStr, "to.revision", "", "Update matching modules with this new revision")

	flag.Parse()
	// end of flags parsing

	fromSource, err := module.ParseSource(fromURL)
	if err != nil {
		return nil, err
	}

	if fromScheme != "" {
		fromSource.Scheme = fromScheme
	}

	if fromHost != "" {
		fromSource.Host = fromHost
	}

	if fromModule != "" {
		fromSource.Module = fromModule
	}

	if fromSubmodule != "" {
		fromSource.Submodule = fromSubmodule
	}

	if fromRevisionStr != "" {
		fromSource.Revision = module.Revision(fromRevisionStr)
	}

	toSource, err := module.ParseSource(toURL)
	if err != nil {
		return nil, err
	}
	if toScheme != "" {
		toSource.Scheme = toScheme
	}

	if toHost != "" {
		toSource.Host = toHost
	}

	if toModule != "" {
		toSource.Module = toModule
	}

	if toSubmodule != "" {
		toSource.Submodule = toSubmodule
	}

	if toRevisionStr != "" {
		toSource.Revision = module.Revision(toRevisionStr)
	}

	config.FromSource = fromSource
	config.ToSource = toSource

	level, err := processing.LevelFromString(logLevel)
	if err != nil {
		return nil, errors.New("cannot parse log level: " + err.Error())
	}

	config.LogLevel = level
	config.Paths = flag.Args()

	return &config, nil
}

func conditionFromSource(source module.Source) (conditions.Condition, error) {
	activeConditions := []conditions.Condition{}
	if source.Scheme != "" {
		activeConditions = append(activeConditions, conditions.SchemeMatches(source.Scheme))
	}

	if source.Host != "" {
		activeConditions = append(activeConditions, conditions.HostMatches(source.Host))
	}

	if source.Module != "" {
		activeConditions = append(activeConditions, conditions.ModuleMatches(source.Module))
	}

	if source.Submodule != "" {
		activeConditions = append(activeConditions, conditions.SubmoduleMatches(source.Submodule))
	}

	if source.Revision != "" {
		activeConditions = append(activeConditions, conditions.RevisionMatches(source.Revision))
	}

	if len(activeConditions) == 0 {
		return conditions.False(), errors.New("no conditions provided")
	}

	return conditions.All(activeConditions...), nil
}
