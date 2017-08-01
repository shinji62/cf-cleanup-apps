package main

import (
	"log"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/shinji62/cf-cleanup-apps/cleanup"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiEndpoint       = kingpin.Flag("api-endpoint", "Api endpoint address. For pcfdev: https://api.pcfdev.local.io").OverrideDefaultFromEnvar("API_ENDPOINT").Required().String()
	user              = kingpin.Flag("user", "Admin user.").OverrideDefaultFromEnvar("CF_API_USER").String()
	password          = kingpin.Flag("password", "Admin password.").OverrideDefaultFromEnvar("CF_API_USER").String()
	clientID          = kingpin.Flag("client-id", "Client ID.").OverrideDefaultFromEnvar("CF_API_CLIENT_ID").String()
	clientSecret      = kingpin.Flag("client-secret", "Client secret.").OverrideDefaultFromEnvar("CF_API_CLIENT_SECRET").String()
	skipSSLValidation = kingpin.Flag("skip-ssl-validation", "Please don't").Default("false").OverrideDefaultFromEnvar("SKIP_SSL_VALIDATION").Bool()
	dryRun            = kingpin.Flag("dry-run", "Dry run ").Default("false").OverrideDefaultFromEnvar("DRY_RUN").Bool()
	durationTime      = kingpin.Flag("app-expired-since", "CloudController Polling time in hour").Default("720h").OverrideDefaultFromEnvar("APP_EXPIRED_SINCE").Duration()
	excludeSystemOrg  = kingpin.Flag("exclude-system-org", "Exclude application in System org to be stopped (Most likely PCF Core App)").Default("true").OverrideDefaultFromEnvar("EXCLUDE_SYSTEM_ORG").Bool()
	excludeOrgs       = kingpin.Flag("exclude-orgs", "Org you want to exclude from cleaning : '--exclude-orgs=myorg1,myorg2 ").Default("").OverrideDefaultFromEnvar("EXCLUDE_ORGS").String()
	includeOrgs       = kingpin.Flag("include-orgs", "Org you want to include from cleaning : '--include-orgs=myorg1,myorg2 ").Default("").OverrideDefaultFromEnvar("EXCLUDE_ORGS").String()
)

func main() {

	kingpin.Parse()

	c := cfclient.Config{
		ApiAddress:        *apiEndpoint,
		Username:          *user,
		Password:          *password,
		ClientID:          *clientID,
		ClientSecret:      *clientSecret,
		SkipSslValidation: *skipSSLValidation,
	}
	cfClient, _ := cfclient.NewClient(&c)

	if *excludeOrgs != "" && *includeOrgs != "" {
		log.Fatalln("You could not use exclude and include as the same time")

	}
	cl := cleanup.NewCleanupCf(cfClient, *durationTime, *excludeSystemOrg)
	cl.SetExcludedOrgs(*excludeOrgs)
	cl.SetIncludedOrgs(*includeOrgs)
	var listExpiredApp *[]cleanup.App
	var err error
	if *includeOrgs != "" {
		listExpiredApp, err = cl.ListExpiredAppsFromApiByOrg()
	} else {
		listExpiredApp, err = cl.ListExpiredAppsFromApi()
	}
	if err != nil {
		log.Fatalf("Error during listing of Expired app", err)
	}
	switch *dryRun {
	case true:
		cl.DryRun(listExpiredApp)
	case false:
		cl.StopApp(listExpiredApp)
	}
}
