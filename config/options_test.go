package config_test

import (
	"flag"

	"github.com/mason-leap-lab/go-utils/config"
	"github.com/mason-leap-lab/go-utils/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MyConfig struct {
	config.Options

	Test bool   `name:"test" description:"Option \"test\"."`
	Name string `name:"name" description:"Option \"name\"."`
}

type MyLoggerConfig struct {
	config.LoggerOptions

	Test bool `name:"test" description:"Option \"test\"."`
}

type MyExtensionConfig struct {
	Extension string `name:"extension" description:"Option \"extension\"."`
}

type MyCompositeConfig struct {
	config.SeedOptions
	Logger config.LoggerOptions
}

type MyCompositeExtension struct {
	config.SeedOptions
	MyExtensionConfig
}

func checkFlagSet(flagSet *flag.FlagSet, err error) {
	Expect(flagSet).NotTo(BeNil())
	if err == config.ErrPrintUsage {
		flagSet.PrintDefaults()
	}
}

var _ = Describe("Options", func() {
	BeforeEach(func() {
		Expect(config.LogLevel).To(Equal(logger.LOG_LEVEL_INFO))
	})

	AfterEach(func() {
		config.LogLevel = logger.LOG_LEVEL_INFO
	})

	It("should polyfill fills the uninitialized Options", func() {
		Expect(config.Polyfill(config.NewOptions(), nil)).To(BeNil())

		myConfig := &MyConfig{}
		Expect(config.Polyfill(myConfig, nil)).To(BeNil())
		Expect(myConfig.Options).NotTo(BeNil())

		compsiteConfig := &MyLoggerConfig{}
		Expect(config.Polyfill(compsiteConfig, nil)).To(BeNil())
		Expect(compsiteConfig.Options).NotTo(BeNil())

		Expect(config.Polyfill(MyConfig{}, nil)).To(Equal(config.ErrNonPointer))
	})

	It("should simple config customization works", func() {
		cfg := &MyConfig{Options: config.NewOptions()}
		flagSet, err := config.ValidateOptionsWithFlags(cfg, "-test")
		checkFlagSet(flagSet, err)

		Expect(err).To(BeNil())
		Expect(cfg.Test).To(Equal(true))
	})

	It("should embeded config customization works", func() {
		var cfg MyLoggerConfig
		flagSet, err := config.ValidateOptionsWithFlags(&cfg, "-debug")
		checkFlagSet(flagSet, err)

		Expect(err).To(BeNil())
		Expect(cfg.Test).To(Equal(false))
		Expect(cfg.Debug).To(Equal(true))
		Expect(config.LogLevel).To(Equal(logger.LOG_LEVEL_ALL))
	})

	It("should composite config has no conflict", func() {
		var cfg MyCompositeConfig
		flagSet, err := config.ValidateOptionsWithFlags(&cfg, "-debug", "-seed=123")
		checkFlagSet(flagSet, err)

		Expect(err).To(BeNil())
		Expect(cfg.Logger.Debug).To(Equal(true))
		Expect(config.LogLevel).To(Equal(logger.LOG_LEVEL_ALL))
		Expect(cfg.Seed).To(Equal(int64(123)))
	})

	It("should composite extension config works", func() {
		var cfg MyCompositeExtension
		flagSet, err := config.ValidateOptionsWithFlags(&cfg, "-seed=123", "-extension=test")
		checkFlagSet(flagSet, err)

		Expect(err).To(BeNil())
		Expect(cfg.Seed).To(Equal(int64(123)))
		Expect(cfg.Extension).To(Equal("test"))
	})

	It("should yaml works for options", func() {
		var cfg MyConfig
		flagSet, err := config.ValidateOptionsWithFlags(&cfg, "-yaml=options_test.yml")
		checkFlagSet(flagSet, err)

		Expect(err).To(BeNil())
		Expect(cfg.Test).To(Equal(true))
	})

	It("should yaml works for embeded options", func() {
		var cfg MyLoggerConfig
		flagSet, err := config.ValidateOptionsWithFlags(&cfg, "-yaml=options_test.yml")
		checkFlagSet(flagSet, err)

		Expect(err).To(BeNil())
		Expect(cfg.Test).To(Equal(true))
		Expect(cfg.Debug).To(Equal(true))
		Expect(config.LogLevel).To(Equal(logger.LOG_LEVEL_ALL))
	})

	It("should yaml works for composite options", func() {
		var cfg MyCompositeConfig
		flagSet, err := config.ValidateOptionsWithFlags(&cfg, "-yaml=options_test.yml")
		checkFlagSet(flagSet, err)

		Expect(err).To(BeNil())
		Expect(cfg.Logger.Debug).To(Equal(true))
		Expect(config.LogLevel).To(Equal(logger.LOG_LEVEL_ALL))
	})

	It("should yaml override default values", func() {
		cfg := MyConfig{Name: "test"}
		flagSet, err := config.ValidateOptionsWithFlags(&cfg, "-yaml=options_test.yml")
		checkFlagSet(flagSet, err)

		Expect(err).To(BeNil())
		Expect(cfg.Test).To(Equal(true))
		Expect(cfg.Name).To(Equal("Tianium"))
	})

	It("should parameter override yaml", func() {
		cfg := MyConfig{Name: "test"}
		flagSet, err := config.ValidateOptionsWithFlags(&cfg, "-yaml=options_test.yml", "-name=Elle")
		checkFlagSet(flagSet, err)

		Expect(err).To(BeNil())
		Expect(cfg.Test).To(Equal(true))
		Expect(cfg.Name).To(Equal("Elle"))
	})
})
