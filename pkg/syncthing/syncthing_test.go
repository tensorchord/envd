// Copyright 2023 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

			s.StopLocalSyncthing()

			Expect(s.IsRunning()).To(BeFalse())
		})
	})

	Describe("Syncthing config", func() {
		It("Initializes local syncthing configuration", func() {
			s, err := syncthing.InitializeLocalSyncthing("s1")
			Expect(err).To(BeNil())

			Expect(s.Port).To(Equal(syncthing.DefaultLocalPort))
			Expect(s.Config.GUI.Address()).To(Equal(fmt.Sprintf("127.0.0.1:%s", syncthing.DefaultLocalPort)))

			dirExists, err := fileutil.DirExists(s.HomeDirectory)
			Expect(err).To(BeNil())
			Expect(dirExists).To(BeTrue())

			configFilePath := syncthing.GetConfigFilePath(s.HomeDirectory)
			fileExists, err := fileutil.FileExists(configFilePath)
			Expect(err).To(BeNil())
			Expect(fileExists).To(BeTrue())

			s.StopLocalSyncthing()
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
		homeDirectory1, err := syncthing.GetHomeDirectory()
		Expect(err).To(BeNil())

		homeDirectory2, err := syncthing.GetHomeDirectory()
		Expect(err).To(BeNil())

		initConfig1 := initConfig.Copy()
		initConfig1.GUI.RawAddress = fmt.Sprintf("127.0.0.1:%s", syncthing.DefaultLocalPort)
		s1 = &syncthing.Syncthing{
			Config:        &initConfig1,
			HomeDirectory: fmt.Sprintf("%s-1", homeDirectory1),
			ApiKey:        syncthing.DefaultApiKey,
			Name:          "s1-REST",
			Port:          syncthing.DefaultLocalPort,
		}

		initConfig2 := initConfig.Copy()
		initConfig2.GUI.RawAddress = fmt.Sprintf("127.0.0.1:%s", syncthing.DefaultRemotePort)
		s2 = &syncthing.Syncthing{
			Config:        &initConfig2,
			HomeDirectory: fmt.Sprintf("%s-2", homeDirectory2),
			ApiKey:        syncthing.DefaultApiKey,
			Name:          "s2-REST",
			Port:          syncthing.DefaultRemotePort,
		}

		s1.Client = s1.NewClient()
		s2.Client = s2.NewClient()

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
		s1.StopLocalSyncthing()

		s2.StopLocalSyncthing()
	})

	It("Connects two local devices", func() {
		err := syncthing.ConnectDevices(s1, s2)
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
		s.StopLocalSyncthing()
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

var _ = Describe("Syncthing util tests", func() {
	It("Parses port correctly", func() {
		addr := "127.0.0.1:8386"

		port := syncthing.ParsePortFromAddress(addr)
		Expect(port).To(Equal("8386"))

		addr2 := "tcp://127.0.0.1:8386"
		port2 := syncthing.ParsePortFromAddress(addr2)
		Expect(port2).To(Equal("8386"))
	})
})
