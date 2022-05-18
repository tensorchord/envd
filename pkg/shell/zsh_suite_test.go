package shell

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestZSH(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ZSH Suite")
}
