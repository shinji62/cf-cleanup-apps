package cleanup_test

import (
	"fmt"
	"time"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/shinji62/cf-cleanup-apps/cleanup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var default_App_time_3h = cfclient.App{
	Guid:             "123",
	Name:             "test app",
	SpaceURL:         "spaceurl1",
	State:            "STARTED",
	PackageUpdatedAt: fmt.Sprintf("%s", time.Now().Add(-1*(3*time.Hour)).Format(time.RFC3339)),
}

var default_exired_app = App{
	Name:             "App_expired",
	Guid:             "uid-1",
	SpaceName:        "SpaceNameOne",
	SpaceGuid:        "suid-1",
	OrgName:          "Org-Name",
	OrgGuid:          "uid-1",
	PackageUpdatedAt: "date",
}

var _ = Describe("Cleanup", func() {
	var clean Cleanup
	var apps []cfclient.App

	BeforeEach(func() {
		clean = NewCleanupCf(&cfclient.Client{}, 2*time.Hour, false)
	})

	Context("With no Application ", func() {
		It("Should return empty List of Expired", func() {
			listapp := clean.ListExpiredApps(&[]cfclient.App{})
			Expect(*listapp).To(BeEmpty())
		})
	})

	Context("With Application Data ", func() {
		BeforeEach(func() {
			apps = []cfclient.App{}
			apps = append(apps, default_App_time_3h)

		})
		AfterEach(func() {
			apps = []cfclient.App{}
		})
		It("Should return empty List when app is expired", func() {
			listapp := clean.ListExpiredApps(&apps)
			Expect(*listapp).ToNot(BeEmpty())
			Expect(*listapp).To(HaveLen(1))
		})
		It("Should not report already stopped application", func() {
			App_time_3h_stopped := default_App_time_3h
			App_time_3h_stopped.State = "STOPPED"
			apps = append(apps, App_time_3h_stopped)
			listapp := clean.ListExpiredApps(&apps)
			Expect(*listapp).ToNot(BeEmpty())
			Expect(*listapp).To(HaveLen(1))
		})
		It("Should not report not expired app", func() {
			App_time_1h := default_App_time_3h
			App_time_1h.Name = "App_time_1h"
			App_time_1h.PackageUpdatedAt = fmt.Sprintf("%s", time.Now().Add(-1*(1*time.Hour)).Format(time.RFC3339))
			apps = []cfclient.App{}
			apps = append(apps, App_time_1h)
			listapp := clean.ListExpiredApps(&apps)
			Expect(*listapp).To(BeEmpty())
		})
		It("Should not report expired ignored app", func() {
			env := make(map[string]interface{})
			env["PCF_DISABLE_CLEANUP"] = true
			App_time_3h := default_App_time_3h
			App_time_3h.Name = "App_time_3h"
			App_time_3h.Environment = env
			apps = append(apps, App_time_3h)
			listapp := clean.ListExpiredApps(&apps)
			Expect(*listapp).ToNot(BeEmpty())
			Expect(*listapp).To(HaveLen(1))
		})
	})
	Context("With Expired App", func() {
		It("Should report right number of ignore app", func() {
			report := clean.CreateReport(&[]App{default_exired_app, default_exired_app})
			clean.DryRun(&[]App{default_exired_app, default_exired_app})
			Expect(report).ToNot(BeEmpty())
			Expect(report).To(HaveKeyWithValue("Org-Name-#-SpaceNameOne", 2))
		})
	})
	Context("With Including System apps", func() {
		It("Should report right number of ignore app", func() {
			clean = NewCleanupCf(&cfclient.Client{}, 2*time.Hour, true)
			App_time_Sytem := default_App_time_3h
			App_time_Sytem.SpaceData.Entity.OrgData.Entity.Name = "system"
			apps = append(apps, App_time_Sytem)
			listapp := clean.ListExpiredApps(&apps)
			report := clean.CreateReport(listapp)
			clean.DryRun(listapp)
			Expect(report).To(BeEmpty())
		})
	})
	Context("With not including System apps", func() {
		It("Should report right number of ignore app", func() {
			App_time_Sytem := default_App_time_3h
			App_time_Sytem.SpaceData.Entity.OrgData.Entity.Name = "system"
			apps = append(apps, App_time_Sytem)
			listapp := clean.ListExpiredApps(&apps)
			report := clean.CreateReport(listapp)
			clean.DryRun(listapp)
			Expect(report).ToNot(BeEmpty())
			Expect(report).To(HaveKeyWithValue("system-#-", 2))
		})
	})

	Context("With excluded orgs", func() {
		It("Should report all excluded orgs", func() {
			clean.SetExcludedOrgs("myorg1,myorG2")
			Expect(clean.GetExcludedOrgs()).ToNot(BeEmpty())
			Expect(clean.GetExcludedOrgs()).To(Equal(map[string]int{"myorg2": 1, "myorg1": 1}))
		})
		It("Should report even if no orgs", func() {
			Expect(clean.GetExcludedOrgs()).To(BeEmpty())
			Expect(clean.GetExcludedOrgs()).To(Equal(map[string]int{}))
		})

		It("Should report right number of ignore app", func() {
			App_time_badorg := default_App_time_3h
			App_time_badorg.SpaceData.Entity.OrgData.Entity.Name = "myorgExcluded"
			App_time_badorg1 := default_App_time_3h
			App_time_badorg1.SpaceData.Entity.OrgData.Entity.Name = "myorgExcluded1"
			apps = []cfclient.App{}
			apps = append(apps, App_time_badorg, App_time_badorg1, default_App_time_3h)
			clean.SetExcludedOrgs("myorgExcluded,myorgExcluded1")
			listapp := clean.ListExpiredApps(&apps)
			report := clean.CreateReport(listapp)
			Expect(report).ToNot(BeEmpty())
			Expect(report).ToNot(HaveKey("myorgExcluded-#-"))
			Expect(report).ToNot(HaveKey("myorgExcluded1-#-"))
		})

	})

})
