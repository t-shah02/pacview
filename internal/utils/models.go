package utils

type PacmanPackage struct {
	Name        string
	Description string
	Version     string
	InstalledAt string
	DependsOn   []string
	RequiredBy  []string
}
