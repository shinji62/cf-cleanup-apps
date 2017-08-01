
# CF-CLEANUP-APPS

This nifty util stop application which have been not been **Updated** since 720 hours (by default ~1month) hours

Why this application?
There is already some application like autosleeping but I think for small usage it's too much overhead.

I think that like development environment need to be alive, meaning application have to be push often, so that's why I decided to create this simple application.

So **Updated** mean application which have been push, application source code updated.

```
./cf-cleanup-apps --api-endpoint=https://api.local.pcfdev.io --skip-ssl-validation --user=admin --password=admin

```

Of course if you don't want your application to be stopped  you can still set the environment variable `PCF_DISABLE_CLEANUP` to true to keep your application running.

```
cf set-env myapp PCF_DISABLE_CLEANUP true
```


# Options

```
./cf-cleanup-apps --help                                                                      ➜  ~/.gvm/pkgsets/go1.7.4/global/src/github.com/shinji62/cf-cleanup-apps
usage: cf-cleanup-apps --api-endpoint=API-ENDPOINT [<flags>]

Flags:
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --api-endpoint=API-ENDPOINT    Api endpoint address. For pcfdev: https://api.pcfdev.local.io
  --user=USER                    Admin user.
  --password=PASSWORD            Admin password.
  --client-id=CLIENT-ID          Client ID.
  --client-secret=CLIENT-SECRET  Client secret.
  --skip-ssl-validation          Please don't
  --dry-run                      Dry run
  --app-expired-since=720h       CloudController Polling time in hour
  --exclude-system-org           Exclude application in System org to be stopped (Most likely PCF Core App)
  --exclude-orgs=""              Org you want to exclude from cleaning : '--exclude-orgs=myorg1,myorg
  --include-orgs=""              Org you want to include from cleaning : '--include-orgs=myorg1,myorg2
```

# Authentification definition
We support 2 types of auth.

## Oauth (Preferred One)

Create Client id/ client Secret

```
uaac target <https://uaa.[your> cf system domain] --skip-ssl-validation

uaac token client get admin -s [your admin-secret]

uaac client add cf-cleanup \
 --secret [your_client_secret] \
 --authorized_grant_types client_credentials,refresh_token \
 --authorities **cloud_controller.admin
```

** Since `cf v241` you can use `cloud_controller.admin_read_only` instead of `cloud_controller.admin`

## Credentials
Create CloudFoundry user
```
cf create-user [cf-cleanup user] [cf-cleanup password]

uaac target <https://uaa.[your> cf system domain] --skip-ssl-validation

uaac token client get admin -s [your admin-secret]

uaac member add **cloud_controller.admin [cf-cleanup user]
```

** Since `cf v241` you can use `cloud_controller.admin_read_only` instead of `cloud_controller.admin`




# To test and build

```
# Setup repo
go get github.com/shinji62/cf-cleanup-apps
cd $GOPATH/src/github.com/shinji62/cf-cleanup-apps

# Test
ginkgo -r .

# Build binary
godep go build
```


# Run this task in Concourse or Other CI

```
./cf-cleanup-apps --api-endpoint=https://api.local.pcfdev.io --skip-ssl-validation --user=admin --password=admin
```




# Run as Task in Cloud Foundry (No YET TESTED)

This application should be run as task everyday of as you want but this do make sense to run as Long Running Process (LRP)

1. Download the latest release of cf-cleanup-apps.

  ```
  git clone https://github.com/shinji62/cf-cleanup-apps
   cd cf-cleanup-apps
  ```

2. Utilize the CF cli to authenticate with your PCF instance.

  ```
  cf login -a https://api.[your cf system domain] -u [your id] --skip-ssl-validation
  ```

3. Push cf-cleanup-apps.

  ```
  cf push cleanup-apps --no-start
  ```

4. Set environment variables with cf  in the [manifest.yml](./manifest.yml).




6. Run the task the app. (cf v247+)

  ```
  cf run-task cf-cleanup-apps "./cf-cleanup-apps" --name cleaning
  ```

  If you are using the offline version of the go buildpack and your app fails to stage then open up the Godeps/Godeps.json file and change the `GoVersion` to a supported one by the buildpacks and repush.
