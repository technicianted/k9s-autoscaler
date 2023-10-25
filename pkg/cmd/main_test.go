package cmd

import (
	"flag"
	"os"
	"testing"

	"k8s.io/klog/v2"
)

func TestMain(m *testing.M) {
	klog.InitFlags(nil)
	flag.CommandLine.Lookup("v").Value.Set("100")
	os.Exit(m.Run())
}
