package cmd

import (
	"aliyun-sls-exporter/pkg/collector"
	"aliyun-sls-exporter/pkg/config"
	"aliyun-sls-exporter/pkg/handler"
	"aliyun-sls-exporter/version"
	"fmt"
	"github.com/spf13/cobra"
)

const AppName = "slsmonitor"

// NewRootCommand create root command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:           AppName,
		Short:         "Exporter for aliyun slsmonitor",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(newServeMetricsCommand())
	cmd.AddCommand(newVersionCommand())
	return cmd
}

func newServeMetricsCommand() *cobra.Command {
	o := &options{
		so: &serveOption{},
	}
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve HTTP metrics handler",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			return o.Complete()
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := config.Parse(o.so.configFile)
			if err != nil {
				return err
			}
			cms, err := collector.NewCloudMonitorCollector(AppName, cfg, o.rateLimit, logger)
			if err != nil {
				return err
			}
			h, err := handler.New(o.so.listenAddress, logger, o.rateLimit, cfg, cms)
			if err != nil {
				return err
			}
			return h.Run()
		},
	}
	o.AddFlags(cmd)
	return cmd
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version info",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(version.Version())
		},
	}
}
