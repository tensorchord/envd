package syncthing_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/facebookgo/subset"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/r3labs/diff/v3"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd/pkg/syncthing"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func TestSyncthing(t *testing.T) {
    logrus.SetLevel(logrus.DebugLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Syncthing Suite")
}

var _ = Describe("Syncthing", func() {

	BeforeEach(func() {
		os.RemoveAll(syncthing.DefaultHomeDirectory())
	})

	Describe("Syncthing", func() {
		It("Starts and stops syncthing", func() {
			s, err := syncthing.InitializeLocalSyncthing()
			Expect(err).To(BeNil())

			Expect(s.IsRunning()).To(BeTrue())

			err = s.StopLocalSyncthing()
			Expect(err).To(BeNil())

			Expect(s.IsRunning()).To(BeFalse())
		})
	})

	Describe("Syncthing config", func() {
		AfterEach(func() {
			os.RemoveAll(syncthing.DefaultHomeDirectory())
		})

		It("Initializes local syncthing configuration", func() {
			s, err := syncthing.InitializeLocalSyncthing()
			Expect(err).To(BeNil())

			Expect(s.Port).To(Equal(syncthing.DefaultLocalPort))
			Expect(s.Config.GUI.Address()).To(Equal(fmt.Sprintf("0.0.0.0:%s", syncthing.DefaultLocalPort)))

			dirExists, err := fileutil.DirExists(s.HomeDirectory)
			Expect(err).To(BeNil())
			Expect(dirExists).To(BeTrue())

			configFilePath := syncthing.GetConfigFilePath(s.HomeDirectory)
			fileExists, err := fileutil.FileExists(configFilePath)
			Expect(err).To(BeNil())
			Expect(fileExists).To(BeTrue())
		})

		It("Initializes remote syncthing configuration", func() {
			s, err := syncthing.InitializeRemoteSyncthing()
			Expect(err).To(BeNil())
			Expect(s.Port).To(Equal(syncthing.DefaultRemotePort))
			Expect(s.Config).ToNot(BeNil())

			configStr := s.Config.String()
			Expect(configStr).To(ContainSubstring(fmt.Sprintf(":%s", syncthing.DefaultRemotePort)))
			Expect(configStr).To(ContainSubstring(fmt.Sprintf("api_key:\"%s\"", s.ApiKey)))
		})

	})

	Describe("Install", func() {
		It("Installs binary in cache directory", func() {
			err := syncthing.InstallSyncthing()
			Expect(err).To(BeNil())

			Expect(syncthing.IsInstalled()).To(BeTrue())
		})
	})
})

var _ = Describe("Syncthing REST API operations", func() {
	var s1 *syncthing.Syncthing
	var s2 *syncthing.Syncthing

	BeforeEach(func() {
		var err error
		initConfig := syncthing.InitLocalConfig()
		homeDirectory := syncthing.DefaultHomeDirectory()

		initConfig1 := initConfig.Copy()
		initConfig1.GUI.RawAddress = fmt.Sprintf("0.0.0.0:%s", syncthing.DefaultLocalPort)
		s1 = &syncthing.Syncthing{
			Config:        &initConfig1,
			HomeDirectory: fmt.Sprintf("%s-1", homeDirectory),
			Client:        syncthing.NewApiClient(),
			ApiKey:        syncthing.DefaultApiKey,
			Port:          syncthing.DefaultLocalPort,
		}

		initConfig2 := initConfig.Copy()
		initConfig2.GUI.RawAddress = fmt.Sprintf("0.0.0.0:%s", syncthing.DefaultRemotePort)
		s2 = &syncthing.Syncthing{
			Config:        &initConfig2,
			HomeDirectory: fmt.Sprintf("%s-2", homeDirectory),
			Client:        syncthing.NewApiClient(),
			ApiKey:        syncthing.DefaultApiKey,
			Port:          syncthing.DefaultRemotePort,
		}

		err = s1.WriteLocalConfig()
		Expect(err).To(BeNil())

		err = s2.WriteLocalConfig()
		Expect(err).To(BeNil())

		err = s1.StartLocalSyncthing()
		Expect(err).To(BeNil())

		err = s2.StartLocalSyncthing()
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err := s1.StopLocalSyncthing()
		Expect(err).To(BeNil())

		err = s2.StopLocalSyncthing()
		Expect(err).To(BeNil())
	})

	It("Connects two local devices", func() {

        err := s1.SetDeviceAddress(syncthing.DefaultLocalDeviceAddress)
		Expect(err).To(BeNil())

        err = s2.SetDeviceAddress(syncthing.DefaultRemoteDeviceAddress)
		Expect(err).To(BeNil())

		err = syncthing.ConnectDevices(s1, s2)
		Expect(err).To(BeNil())
	})

})

var _ = Describe("Syncthing REST API operations", func() {
	var s *syncthing.Syncthing

	BeforeEach(func() {
		var err error
		s, err = syncthing.InitializeLocalSyncthing()
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err := s.StopLocalSyncthing()
		Expect(err).To(BeNil())
	})

	It("Connects local syncthing to running remote syncthing", func() {
		s1, err := syncthing.InitializeRemoteSyncthing()
		Expect(err).To(BeNil())

		err = syncthing.ConnectDevices(s, s1)
		Expect(err).To(BeNil())
	})

	It("Applies syncthing configuration twice", func() {
		s.Config.GUI.Debugging = false
		Expect(s.Config.GUI.Debugging).To(Equal(false))

		s.Config.GUI.Debugging = true

		ok, err := s.Ping()
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())

		err = s.ApplyConfig()
		Expect(err).To(BeNil())

		cfg, err := s.GetConfig()
		Expect(err).To(BeNil())
		Expect(cfg.GUI.Debugging).To(Equal(true))

		s.Config.GUI.Debugging = false
		err = s.ApplyConfig()
		Expect(err).To(BeNil())

		cfg, err = s.GetConfig()
		Expect(err).To(BeNil())
		Expect(cfg.GUI.Debugging).To(Equal(false))
	})

	It("Gets the most recent event", func() {
		event, err := s.GetMostRecentEvent()
		Expect(err).To(BeNil())
		Expect(event.Id > 0).To(BeTrue())
	})

})

var _ = Describe("Util tests", func() {
	Describe("Subset", func() {
		It("Subset works", func() {
			type A struct {
				Hi  string
				Bye string
			}
			type B struct {
				One   string
				Two   string
				Three A
			}

			a := A{Hi: "hi", Bye: "bye"}
			bSuperset := B{One: "one", Two: "two", Three: a}

			bSubset := B{One: "one", Two: "two"}
			Expect(subset.Check(bSubset, bSuperset)).To(BeTrue())
		})

		It("Diff works", func() {
			type B struct {
				One string
			}
			type C struct {
				One string
			}
			type A struct {
				One   string
				Two   string
				Three string
				Four  string
				Five  B
				Six   C
			}

			var sb = B{One: "one"}
			var sc = C{One: "two"}

			var a = A{One: "one", Two: "two", Three: "three", Five: sb}
			var b = A{Two: "two", Three: "five", Four: "four", Six: sc}
			var d = A{One: "two", Two: "two", Three: "five", Four: "four", Six: sc}
			var c A

			_, err := diff.Merge(a, b, &c)
			Expect(err).To(BeNil())

			Expect(c.One).ToNot(Equal("one"))                   // No changes to field
			Expect(c.Two).ToNot(Equal("two"))                   // No changes to field
			Expect(c.Three).To(Equal("five"))                   // Changed field
			Expect(c.Four).To(Equal("four"))                    // New field
			Expect(reflect.DeepEqual(c.Five, sb)).To(BeFalse()) // Removed field, no change
			Expect(reflect.DeepEqual(c.Six, sc)).To(BeTrue())   // New field

			Expect(subset.Check(c, a)).To(BeFalse()) // Original
			Expect(subset.Check(c, b)).To(BeTrue())  // With changed fields
			Expect(subset.Check(c, d)).To(BeTrue())  // With changed but irrelevant field
		})
	})
})
