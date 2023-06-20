package version

import (
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	version := &version{
		SemVersion: "0.0.0",
		GitCommit:  "ef30d58",
		BuildDate:  "2024-03-25_17:55:07",
		GoVersion:  "go1.20",
		Compiler:   "gc",
		Platform:   "linux/test",
	}

	Context("testing version Print function", func() {
		It("should not error", func() {
			Expect(version.Print()).Error().ShouldNot(HaveOccurred())
		})
	})

	Context("testing newVersion function", func() {
		It("should return the expected version struct", func() {
			v := newVersion()
			Expect(v.Compiler).Should(Equal(runtime.Compiler))
			Expect(v.SemVersion).Should(Equal(semVersion))
			Expect(v.GitCommit).Should(Equal(gitCommit))
			Expect(v.BuildDate).Should(Equal(buildDate))
			Expect(v.GoVersion).Should(Equal(runtime.Version()))
			Expect(v.Platform).Should(Equal(fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)))
		})
	})
})

func TestVersion(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Version Suite")
}
