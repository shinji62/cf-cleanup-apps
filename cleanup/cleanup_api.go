package cleanup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/shinji62/cf-cleanup-apps/logging"
)

type Cfrequest struct {
	State string `json:"state"`
}

type CleanupCf struct {
	cf               *cfclient.Client
	dateExpired      time.Duration
	excludeSystemOrg bool
	excludedOrgs     map[string]int
}

func NewCleanupCf(cfC *cfclient.Client, since time.Duration, excludeSystem bool) Cleanup {
	return &CleanupCf{
		cf:               cfC,
		dateExpired:      since,
		excludeSystemOrg: excludeSystem,
		excludedOrgs:     map[string]int{},
	}
}

func (cl *CleanupCf) DryRun(expiredApps *[]App) {
	logging.LogStd(fmt.Sprintf("Found [%d] apps to be stopped ", len(*expiredApps)), true)
	for orgSpace, num := range cl.CreateReport(expiredApps) {
		logging.LogStd(fmt.Sprintf("Stopping [%d] apps in [%s] ", num, strings.Split(orgSpace, "-#-")), true)

	}
}

func (cl *CleanupCf) StopApp(expiredApps *[]App) error {
	stoppedReader, _ := json.Marshal(&Cfrequest{State: "STOPPED"})

	for _, app := range *expiredApps {
		//TODO this could not be so horrible
		buf := bytes.NewBuffer(stoppedReader)
		req := cl.cf.NewRequestWithBody("PUT", fmt.Sprintf("/v2/apps/%s", app.Guid), buf)

		res, err := cl.cf.DoRequest(req)
		res.Body.Close()
		if err != nil {
			return err
		}

		if res.StatusCode != 201 {
			logging.LogError(res.Status+" when calling cf api", res.StatusCode)
		}
		logging.LogStd(fmt.Sprintf("Stopping Application %s in Org: %s Space: %s Last Updated: %s", app.Name, app.OrgName, app.SpaceName, app.PackageUpdated), true)
		//Avoid Killing CF API let's put some sleep
		time.Sleep(2 * time.Second)
	}
	return nil

}

func (cl *CleanupCf) getAppsFromApi() ([]cfclient.App, error) {
	return cl.cf.ListApps()
}

func (cl *CleanupCf) ListExpiredAppsFromApi() (*[]App, error) {
	aApi, err := cl.getAppsFromApi()
	//TODO FIX THAT!!!
	// if err != nil {
	// 	return *[]App, nil
	// }
	return cl.ListExpiredApps(&aApi), err
}

func (cl *CleanupCf) ListExpiredApps(listapps *[]cfclient.App) *[]App {
	apps := []App{}
	for _, app := range *listapps {
		if !cl.isStarted(app.State) || cl.isOptOut(app.Environment) || cl.isSystemOrg(app.SpaceData.Entity.OrgData.Entity.Name) {
			continue
		}
		if !cl.isExpired(app.PackageUpdated) {
			continue
		}

		if cl.isExcludedOrg(app.SpaceData.Entity.OrgData.Entity.Name) {
			continue
		}

		apps = append(apps, App{
			app.Name,
			app.Guid,
			app.SpaceData.Entity.Name,
			app.SpaceData.Entity.Guid,
			app.SpaceData.Entity.OrgData.Entity.Name,
			app.SpaceData.Entity.OrgData.Entity.Guid,
			app.PackageUpdated,
		})
	}
	return &apps
}

func (cl *CleanupCf) isExcludedOrg(orgName string) bool {
	_, ok := cl.excludedOrgs[strings.ToLower(orgName)]
	return ok
}
func (c *CleanupCf) isStarted(state string) bool {
	return state == "STARTED"
}

func (cl *CleanupCf) isOptOut(envVar map[string]interface{}) bool {
	if val, ok := envVar["PCF_DISABLE_CLEANUP"]; ok != false && (val == "true" || val == true) {
		return true
	}
	return false
}

func (cl *CleanupCf) SetExcludedOrgs(excludeOrg string) {
	excludedOrg := map[string]int{}
	for _, kvPair := range strings.Split(excludeOrg, ",") {
		if kvPair != "" {
			excludedOrg[strings.ToLower(strings.TrimSpace(kvPair))] = 1
		}
	}
	cl.excludedOrgs = excludedOrg
}

func (cl *CleanupCf) GetExcludedOrgs() map[string]int {
	return cl.excludedOrgs
}

func (cl *CleanupCf) isSystemOrg(orgName string) bool {
	return strings.ToLower(orgName) == "system" && cl.excludeSystemOrg
}

func (cl *CleanupCf) isExpired(updated string) bool {
	t, _ := time.Parse(time.RFC3339, updated)
	return (time.Since(t) > cl.dateExpired)
}

func (cl *CleanupCf) CreateReport(expireApp *[]App) map[string]int {
	result := map[string]int{}
	for _, app := range *expireApp {
		result[app.OrgName+"-#-"+app.SpaceName]++
	}
	return result
}
