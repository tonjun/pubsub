package main_test

import (
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPubsub(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pubsub Suite")
}
