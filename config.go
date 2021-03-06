package main

import (
	"io/ioutil"
	"strings"

	"github.com/drone/drone/pkg/build/log"
)

type SshConfigFileSection struct {
	Host         string
	ForwardAgent string
	User         string
	HostName     string
	Port         string
}

// parseSshConfigFileSection parses a section from the ~/.ssh/config file
func parseSshConfigFileSection(content string) *SshConfigFileSection {
	section := &SshConfigFileSection{}

	for n, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if n == 0 {
			section.Host = line
		} else if strings.HasPrefix(line, "ForwardAgent") {
			section.ForwardAgent = strings.TrimSpace(strings.TrimPrefix(line, "ForwardAgent"))
		} else if strings.HasPrefix(line, "User") {
			section.User = strings.TrimSpace(strings.TrimPrefix(line, "User"))
		} else if strings.HasPrefix(line, "HostName") {
			section.HostName = strings.TrimSpace(strings.TrimPrefix(line, "HostName"))
		} else if strings.HasPrefix(line, "Port") {
			section.Port = strings.TrimSpace(strings.TrimPrefix(line, "Port"))
		}
	}
	log.Debugf("parsed ssh config file section: %s", section.Host)
	return section
}

// parseSshConfigFile parses the ~/.ssh/config file and build a list of section
func parseSshConfigFile(path string) (map[string]*SshConfigFileSection, error) {
	log.Debugf("parsing ssh config file: %s", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	sections := make(map[string]*SshConfigFileSection)
	for _, split := range strings.Split(string(content), "Host ") {
		split = strings.TrimSpace(split)
		if split == "" {
			continue
		}

		section := parseSshConfigFileSection(split)
		sections[section.Host] = section
	}

	return sections, nil
}
