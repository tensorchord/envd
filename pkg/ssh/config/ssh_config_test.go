package config

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ssh config", func() {
	var sshConfigPath = "ssh_config"
	When("giving a empty ssh config", func() {
		It("Should add/remove the config successfully", func() {
			env := "test-ssh-config"
			iface := "localhost"
			port := 8888
			keyPath := "key"
			err := add(sshConfigPath, env, iface, port, keyPath)
			Expect(err).NotTo(HaveOccurred())

			err = remove(sshConfigPath, env)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
