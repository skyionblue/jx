package opts

const (
	// PlaceHolderPrefix is prefix for placeholders
	PlaceHolderPrefix = "REPLACE_ME"
	// PlaceHolderAppName placeholder for app name
	PlaceHolderAppName = PlaceHolderPrefix + "_APP_NAME"
	// PlaceHolderGitProvider placeholder for git provider
	PlaceHolderGitProvider = PlaceHolderPrefix + "_GIT_PROVIDER"
	// PlaceHolderOrg placeholder for org
	PlaceHolderOrg = PlaceHolderPrefix + "_ORG"
	// PlaceHolderDockerRegistryOrg placeholder for docker registry
	PlaceHolderDockerRegistryOrg = PlaceHolderPrefix + "_DOCKER_REGISTRY_ORG"

	MinimumMavenDeployVersion = "2.8.2"

	MasterBranch         = "master"
	DefaultGitIgnoreFile = `
.project
.classpath
.idea
.cache
.DS_Store
*.im?
target
work
`
)
