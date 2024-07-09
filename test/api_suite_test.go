//go:build test || all
// +build test all

package test

import (
	"flag"
	"github.com/liuzhaomax/go-maxms/test/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"testing"
)

var (
	env      = flag.String("env", "dev", "环境")
	junitXML = "junit.xml"
)

func TestAPI(t *testing.T) {
	common.LoadConfig(*env)
	flag.Parse()
	log.Println("运行环境为：", *env)
	RegisterFailHandler(Fail)
	suiteConfig, reporterConfig := GinkgoConfiguration()
	reporterConfig.JUnitReport = junitXML
	reporterConfig.FullTrace = true
	RunSpecs(t, "API Suite", suiteConfig, reporterConfig)
}
