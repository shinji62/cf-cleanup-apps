package cleanup

import cfclient "github.com/cloudfoundry-community/go-cfclient"

//go:generate counterfeiter . Cleanup

type App struct {
	Name             string
	Guid             string
	SpaceName        string
	SpaceGuid        string
	OrgName          string
	OrgGuid          string
	PackageUpdatedAt string
}

type Cleanup interface {
	StopApp(*[]App) error
	DryRun(*[]App)
	ListExpiredApps(*[]cfclient.App) *[]App
	ListExpiredAppsFromApi() (*[]App, error)
	ListExpiredAppsFromApiByOrg() (*[]App, error)
	CreateReport(expireApp *[]App) map[string]int
	SetExcludedOrgs(excludeOrg string)
	GetExcludedOrgs() map[string]int
	SetIncludedOrgs(includeOrg string)
}
