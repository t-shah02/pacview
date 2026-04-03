package utils

import (
	"log/slog"
	"os/exec"
	"strings"
)

var logger *slog.Logger = GetLogger()

func parseSpaceSeparatedPacmanList(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" || strings.EqualFold(value, "none") {
		return nil
	}
	return strings.Fields(value)
}

func parseRawQiOutput() ([]PacmanPackage, bool) {
	pacmanPath, err := exec.LookPath("pacman")
	if err != nil {
		logger.Error("could not find a command called pacman on the host machine")
		return nil, false
	}

	cmd := exec.Command(pacmanPath, "-Qi")
	stdout, err := cmd.Output()
	if err != nil {
		logger.Error("failed to run command to fetch installed pacman packages on the hosted machine")
		return nil, false
	}

	return ParseQiStdout(string(stdout)), true
}

func ParseQiStdout(stdout string) []PacmanPackage {
	var packages []PacmanPackage
	var buf PacmanPackage
	inPackage := false

	flush := func() {
		if !inPackage {
			return
		}
		packages = append(packages, buf)
		buf = PacmanPackage{}
		inPackage = false
	}

	for line := range strings.SplitSeq(stdout, "\n") {
		if strings.TrimSpace(line) == "" {
			flush()
			continue
		}

		key, value, ok := strings.Cut(strings.TrimSpace(line), ":")
		if !ok {
			continue
		}
		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.TrimSpace(value)

		switch key {
		case "name":
			flush()
			buf.Name = value
			inPackage = true
		case "version":
			buf.Version = value
		case "description":
			buf.Description = value
		case "install date":
			buf.InstalledAt = value
		case "depends on":
			buf.DependsOn = parseSpaceSeparatedPacmanList(value)
		case "required by":
			buf.RequiredBy = parseSpaceSeparatedPacmanList(value)
		}
	}
	flush()

	return packages
}

func GetInstalledPacmanPackages() []PacmanPackage {
	packages, ok := parseRawQiOutput()
	if !ok {
		return nil
	}
	return packages
}
