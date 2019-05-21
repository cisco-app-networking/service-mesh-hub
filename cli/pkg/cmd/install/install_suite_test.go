package install_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var T *testing.T

func TestInstall(t *testing.T) {
	RegisterFailHandler(Fail)
	T = t
	RunSpecs(t, "Install Suite")
}
