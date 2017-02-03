package main_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var testBinary string

func TestIntegrations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integrations Suite")
}

var _ = BeforeSuite(func() {
	var err error
	testBinary, err = gexec.Build("github.com/shinji62/cf-cleanup-apps/logging/integrations/")
	fmt.Println(testBinary)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("Logging Testing", func() {

	It("Stdout test", func() {
		session, err := gexec.Start(exec.Command(testBinary), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		session.Wait()
		Expect(session.Out.Contents()).To(ContainSubstring("Should OutPut this message to stdin"))
		Expect(session.Err.Contents()).To(ContainSubstring("Should OutPut this message to sdtout Details: ErrorMessage"))
	})

})
