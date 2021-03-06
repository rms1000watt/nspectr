package cmd

import (
	"github.com/rms1000watt/nspectr/proxy"

	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:     "proxy",
	Short:   "Start the proxy",
	Example: `./nspectr proxy`,
	Run:     proxyFunc,
}

var proxyCfg proxy.Config

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.Flags().StringVar(&proxyCfg.Host, "host", "", "Host to listen on")
	proxyCmd.Flags().IntVar(&proxyCfg.Port, "port", 7100, "Port to listen on")
	proxyCmd.Flags().StringVar(&proxyCfg.BackendAddr, "backend-addr", "127.0.0.1:8080", "Address of backend to proxy to")

	setFlagsFromEnv(proxyCmd)
}

func proxyFunc(cmd *cobra.Command, args []string) {
	configureLogging()

	proxy.Proxy(proxyCfg)
}
