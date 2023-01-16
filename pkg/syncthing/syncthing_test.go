package syncthing_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
	})

	Describe("Syncthing", func() {
		It("Starts and stops syncthing", func() {
			s, err := syncthing.InitializeLocalSyncthing("s1")
			Expect(err).To(BeNil())

			Expect(s.IsRunning()).To(BeTrue())

			err = s.StopLocalSyncthing()
			Expect(err).To(BeNil())

			Expect(s.IsRunning()).To(BeFalse())
		})
	})

	Describe("Syncthing config", func() {
		AfterEach(func() {
		})

		It("Initializes local syncthing configuration", func() {
			s, err := syncthing.InitializeLocalSyncthing("s1")
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
		homeDirectory1 := syncthing.GetHomeDirectory("s1")
		homeDirectory2 := syncthing.GetHomeDirectory("s2")

		initConfig1 := initConfig.Copy()
		initConfig1.GUI.RawAddress = fmt.Sprintf("0.0.0.0:%s", syncthing.DefaultLocalPort)
		s1 = &syncthing.Syncthing{
			Config:        &initConfig1,
			HomeDirectory: fmt.Sprintf("%s-1", homeDirectory1),
			Client:        syncthing.NewApiClient(),
			ApiKey:        syncthing.DefaultApiKey,
			Port:          syncthing.DefaultLocalPort,
		}

		initConfig2 := initConfig.Copy()
		initConfig2.GUI.RawAddress = fmt.Sprintf("0.0.0.0:%s", syncthing.DefaultRemotePort)
		s2 = &syncthing.Syncthing{
			Config:        &initConfig2,
			HomeDirectory: fmt.Sprintf("%s-2", homeDirectory2),
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
		s, err = syncthing.InitializeLocalSyncthing("s1")
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
