package v7

import (
	"context"
	"io"
	"net/http"
	"time"

	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v7action"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/constant"
	uaa "code.cloudfoundry.org/cli/api/uaa/constant"
	"code.cloudfoundry.org/cli/cf/configuration/coreconfig"
	"code.cloudfoundry.org/cli/resources"
	"code.cloudfoundry.org/cli/types"
	"code.cloudfoundry.org/cli/util/configv3"
	"github.com/SermoDigital/jose/jwt"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Actor

type Actor interface {
	ApplyOrganizationQuotaByName(quotaName string, orgGUID string) (v7action.Warnings, error)
	ApplySpaceQuotaByName(quotaName string, spaceGUID string, orgGUID string) (v7action.Warnings, error)
	AssignIsolationSegmentToSpaceByNameAndSpace(isolationSegmentName string, spaceGUID string) (v7action.Warnings, error)
	Authenticate(credentials map[string]string, origin string, grantType uaa.GrantType) error
	BindSecurityGroupToSpaces(securityGroupGUID string, spaces []resources.Space, lifecycle constant.SecurityGroupLifecycle) (v7action.Warnings, error)
	CancelDeployment(deploymentGUID string) (v7action.Warnings, error)
	CheckRoute(domainName string, hostname string, path string, port int) (bool, v7action.Warnings, error)
	ClearTarget()
	CopyPackage(sourceApp resources.Application, targetApp resources.Application) (resources.Package, v7action.Warnings, error)
	CreateAndUploadBitsPackageByApplicationNameAndSpace(appName string, spaceGUID string, bitsPath string) (resources.Package, v7action.Warnings, error)
	CreateApplicationDroplet(appGUID string) (resources.Droplet, v7action.Warnings, error)
	CreateApplicationInSpace(app resources.Application, spaceGUID string) (resources.Application, v7action.Warnings, error)
	CreateBitsPackageByApplication(appGUID string) (resources.Package, v7action.Warnings, error)
	CreateBuildpack(buildpack resources.Buildpack) (resources.Buildpack, v7action.Warnings, error)
	CreateDeploymentByApplicationAndDroplet(appGUID string, dropletGUID string) (string, v7action.Warnings, error)
	CreateDeploymentByApplicationAndRevision(appGUID string, revisionGUID string) (string, v7action.Warnings, error)
	CreateDockerPackageByApplication(appGUID string, dockerImageCredentials v7action.DockerImageCredentials) (resources.Package, v7action.Warnings, error)
	CreateDockerPackageByApplicationNameAndSpace(appName string, spaceGUID string, dockerImageCredentials v7action.DockerImageCredentials) (resources.Package, v7action.Warnings, error)
	CreateIsolationSegmentByName(isolationSegment resources.IsolationSegment) (v7action.Warnings, error)
	CreateManagedServiceInstance(managedServiceInstanceParams v7action.CreateManagedServiceInstanceParams) (chan v7action.PollJobEvent, v7action.Warnings, error)
	CreateOrgRole(roleType constant.RoleType, orgGUID string, userNameOrGUID string, userOrigin string, isClient bool) (v7action.Warnings, error)
	CreateOrganization(orgName string) (resources.Organization, v7action.Warnings, error)
	CreateOrganizationQuota(name string, limits v7action.QuotaLimits) (v7action.Warnings, error)
	CreatePrivateDomain(domainName string, orgName string) (v7action.Warnings, error)
	CreateRoute(spaceGUID, domainName, hostname, path string, port int) (resources.Route, v7action.Warnings, error)
	CreateRouteBinding(params v7action.CreateRouteBindingParams) (chan v7action.PollJobEvent, v7action.Warnings, error)
	CreateSecurityGroup(name, filePath string) (v7action.Warnings, error)
	CreateServiceAppBinding(params v7action.CreateServiceAppBindingParams) (chan v7action.PollJobEvent, v7action.Warnings, error)
	CreateServiceBroker(model resources.ServiceBroker) (v7action.Warnings, error)
	CreateServiceKey(params v7action.CreateServiceKeyParams) (chan v7action.PollJobEvent, v7action.Warnings, error)
	CreateSharedDomain(domainName string, internal bool, routerGroupName string) (v7action.Warnings, error)
	CreateSpace(spaceName, orgGUID string) (resources.Space, v7action.Warnings, error)
	CreateSpaceQuota(spaceQuotaName string, orgGuid string, limits v7action.QuotaLimits) (v7action.Warnings, error)
	CreateSpaceRole(roleType constant.RoleType, orgGUID string, spaceGUID string, userNameOrGUID string, userOrigin string, isClient bool) (v7action.Warnings, error)
	CreateUser(username string, password string, origin string) (resources.User, v7action.Warnings, error)
	CreateUserProvidedServiceInstance(instance resources.ServiceInstance) (v7action.Warnings, error)
	DeleteApplicationByNameAndSpace(name, spaceGUID string, deleteRoutes bool) (v7action.Warnings, error)
	DeleteBuildpackByNameAndStack(buildpackName string, buildpackStack string) (v7action.Warnings, error)
	DeleteDomain(domain resources.Domain) (v7action.Warnings, error)
	DeleteInstanceByApplicationNameSpaceProcessTypeAndIndex(appName string, spaceGUID string, processType string, instanceIndex int) (v7action.Warnings, error)
	DeleteOrgRole(roleType constant.RoleType, orgGUID string, userNameOrGUID string, userOrigin string, isClient bool) (v7action.Warnings, error)
	DeleteOrganization(orgName string) (v7action.Warnings, error)
	DeleteOrganizationQuota(quotaName string) (v7action.Warnings, error)
	DeleteOrphanedRoutes(spaceGUID string) (v7action.Warnings, error)
	DeleteRoute(domainName, hostname, path string, port int) (v7action.Warnings, error)
	DeleteRouteBinding(params v7action.DeleteRouteBindingParams) (chan v7action.PollJobEvent, v7action.Warnings, error)
	DeleteSecurityGroup(securityGroupName string) (v7action.Warnings, error)
	DeleteServiceAppBinding(params v7action.DeleteServiceAppBindingParams) (chan v7action.PollJobEvent, v7action.Warnings, error)
	DeleteServiceBroker(serviceBrokerGUID string) (v7action.Warnings, error)
	DeleteServiceInstance(serviceInstanceName, spaceGUID string) (chan v7action.PollJobEvent, v7action.Warnings, error)
	DeleteServiceKeyByServiceInstanceAndName(serviceInstanceName, serviceKeyName, spaceGUID string) (chan v7action.PollJobEvent, v7action.Warnings, error)
	DeleteSpaceByNameAndOrganizationName(spaceName string, orgName string) (v7action.Warnings, error)
	DeleteSpaceQuotaByName(quotaName string, orgGUID string) (v7action.Warnings, error)
	DeleteSpaceRole(roleType constant.RoleType, spaceGUID string, userNameOrGUID string, userOrigin string, isClient bool) (v7action.Warnings, error)
	DeleteUser(userGuid string) (v7action.Warnings, error)
	DeleteIsolationSegmentByName(name string) (v7action.Warnings, error)
	DeleteIsolationSegmentOrganizationByName(isolationSegmentName string, orgName string) (v7action.Warnings, error)
	DiffSpaceManifest(spaceGUID string, rawManifest []byte) (resources.ManifestDiff, v7action.Warnings, error)
	DisableFeatureFlag(flagName string) (v7action.Warnings, error)
	DisableServiceAccess(offeringName, brokerName, orgName, planName string) (v7action.SkippedPlans, v7action.Warnings, error)
	DownloadCurrentDropletByAppName(appName string, spaceGUID string) ([]byte, string, v7action.Warnings, error)
	DownloadDropletByGUIDAndAppName(dropletGUID string, appName string, spaceGUID string) ([]byte, v7action.Warnings, error)
	EnableFeatureFlag(flagName string) (v7action.Warnings, error)
	EnableServiceAccess(offeringName, brokerName, orgName, planName string) (v7action.SkippedPlans, v7action.Warnings, error)
	EntitleIsolationSegmentToOrganizationByName(isolationSegmentName string, orgName string) (v7action.Warnings, error)
	GetAppFeature(appGUID string, featureName string) (resources.ApplicationFeature, v7action.Warnings, error)
	GetAppSummariesForSpace(spaceGUID string, labels string) ([]v7action.ApplicationSummary, v7action.Warnings, error)
	GetApplicationByNameAndSpace(appName string, spaceGUID string) (resources.Application, v7action.Warnings, error)
	GetApplicationMapForRoute(route resources.Route) (map[string]resources.Application, v7action.Warnings, error)
	GetApplicationDroplets(appName string, spaceGUID string) ([]resources.Droplet, v7action.Warnings, error)
	GetApplicationLabels(appName string, spaceGUID string) (map[string]types.NullString, v7action.Warnings, error)
	GetApplicationPackages(appName string, spaceGUID string) ([]resources.Package, v7action.Warnings, error)
	GetApplicationProcessHealthChecksByNameAndSpace(appName string, spaceGUID string) ([]v7action.ProcessHealthCheck, v7action.Warnings, error)
	GetApplicationRevisionsDeployed(appGUID string) ([]resources.Revision, v7action.Warnings, error)
	GetApplicationRoutes(appGUID string) ([]resources.Route, v7action.Warnings, error)
	GetApplicationTasks(appName string, sortOrder v7action.SortOrder) ([]resources.Task, v7action.Warnings, error)
	GetApplicationsByNamesAndSpace(appNames []string, spaceGUID string) ([]resources.Application, v7action.Warnings, error)
	GetBuildpackLabels(buildpackName string, buildpackStack string) (map[string]types.NullString, v7action.Warnings, error)
	GetBuildpacks(labelSelector string) ([]resources.Buildpack, v7action.Warnings, error)
	GetCurrentUser() (configv3.User, error)
	GetDefaultDomain(orgGUID string) (resources.Domain, v7action.Warnings, error)
	GetDetailedAppSummary(appName string, spaceGUID string, withObfuscatedValues bool) (v7action.DetailedApplicationSummary, v7action.Warnings, error)
	GetDomain(domainGUID string) (resources.Domain, v7action.Warnings, error)
	GetDomainByName(domainName string) (resources.Domain, v7action.Warnings, error)
	GetDomainLabels(domainName string) (map[string]types.NullString, v7action.Warnings, error)
	GetEffectiveIsolationSegmentBySpace(spaceGUID string, orgDefaultIsolationSegmentGUID string) (resources.IsolationSegment, v7action.Warnings, error)
	GetEnvironmentVariableGroup(group constant.EnvironmentVariableGroupName) (v7action.EnvironmentVariableGroup, v7action.Warnings, error)
	GetEnvironmentVariablesByApplicationNameAndSpace(appName string, spaceGUID string) (v7action.EnvironmentVariableGroups, v7action.Warnings, error)
	GetFeatureFlagByName(featureFlagName string) (resources.FeatureFlag, v7action.Warnings, error)
	GetFeatureFlags() ([]resources.FeatureFlag, v7action.Warnings, error)
	GetGlobalRunningSecurityGroups() ([]resources.SecurityGroup, v7action.Warnings, error)
	GetGlobalStagingSecurityGroups() ([]resources.SecurityGroup, v7action.Warnings, error)
	GetIsolationSegmentsByOrganization(orgName string) ([]resources.IsolationSegment, v7action.Warnings, error)
	GetIsolationSegmentByName(isoSegmentName string) (resources.IsolationSegment, v7action.Warnings, error)
	GetIsolationSegmentSummaries() ([]v7action.IsolationSegmentSummary, v7action.Warnings, error)
	GetLatestActiveDeploymentForApp(appGUID string) (resources.Deployment, v7action.Warnings, error)
	GetLoginPrompts() (map[string]coreconfig.AuthPrompt, error)
	GetNewestReadyPackageForApplication(app resources.Application) (resources.Package, v7action.Warnings, error)
	GetOrgUsersByRoleType(orgGUID string) (map[constant.RoleType][]resources.User, v7action.Warnings, error)
	GetOrganizationByName(orgName string) (resources.Organization, v7action.Warnings, error)
	GetOrganizationDomains(string, string) ([]resources.Domain, v7action.Warnings, error)
	GetOrganizationLabels(orgName string) (map[string]types.NullString, v7action.Warnings, error)
	GetOrganizationQuotaByName(orgQuotaName string) (resources.OrganizationQuota, v7action.Warnings, error)
	GetOrganizationQuotas() ([]resources.OrganizationQuota, v7action.Warnings, error)
	GetOrganizationSpaces(orgGUID string) ([]resources.Space, v7action.Warnings, error)
	GetOrganizationSpacesWithLabelSelector(orgGUID string, labelSelector string) ([]resources.Space, v7action.Warnings, error)
	GetOrganizationSummaryByName(orgName string) (v7action.OrganizationSummary, v7action.Warnings, error)
	GetOrganizations(labelSelector string) ([]resources.Organization, v7action.Warnings, error)
	GetProcessByTypeAndApplication(processType string, appGUID string) (resources.Process, v7action.Warnings, error)
	GetRawApplicationManifestByNameAndSpace(appName string, spaceGUID string) ([]byte, v7action.Warnings, error)
	GetRecentEventsByApplicationNameAndSpace(appName string, spaceGUID string) ([]v7action.Event, v7action.Warnings, error)
	GetRecentLogsForApplicationByNameAndSpace(appName string, spaceGUID string, client sharedaction.LogCacheClient) ([]sharedaction.LogMessage, v7action.Warnings, error)
	GetRootResponse() (v7action.Info, v7action.Warnings, error)
	GetRevisionByApplicationAndVersion(appGUID string, revisionVersion int) (resources.Revision, v7action.Warnings, error)
	GetRevisionsByApplicationNameAndSpace(appName string, spaceGUID string) ([]resources.Revision, v7action.Warnings, error)
	GetRouteByAttributes(domain resources.Domain, hostname string, path string, port int) (resources.Route, v7action.Warnings, error)
	GetRouteDestinationByAppGUID(route resources.Route, appGUID string) (resources.RouteDestination, error)
	GetRouteLabels(routeName string, spaceGUID string) (map[string]types.NullString, v7action.Warnings, error)
	GetRouterGroups() ([]v7action.RouterGroup, error)
	GetRouteSummaries([]resources.Route) ([]v7action.RouteSummary, v7action.Warnings, error)
	GetRoutesByOrg(orgGUID string, labels string) ([]resources.Route, v7action.Warnings, error)
	GetRoutesBySpace(spaceGUID string, labels string) ([]resources.Route, v7action.Warnings, error)
	GetSSHEnabled(appGUID string) (ccv3.SSHEnabled, v7action.Warnings, error)
	GetSSHEnabledByAppName(appName string, spaceGUID string) (ccv3.SSHEnabled, v7action.Warnings, error)
	GetSSHPasscode() (string, error)
	GetSecureShellConfigurationByApplicationNameSpaceProcessTypeAndIndex(appName string, spaceGUID string, processType string, processIndex uint) (v7action.SSHAuthentication, v7action.Warnings, error)
	GetSecurityGroup(securityGroupName string) (resources.SecurityGroup, v7action.Warnings, error)
	GetSecurityGroupSummary(securityGroupName string) (v7action.SecurityGroupSummary, v7action.Warnings, error)
	GetSecurityGroups() ([]v7action.SecurityGroupSummary, v7action.Warnings, error)
	GetServiceAccess(offeringName, brokerName, orgName string) ([]v7action.ServicePlanAccess, v7action.Warnings, error)
	GetServiceBrokerByName(serviceBrokerName string) (resources.ServiceBroker, v7action.Warnings, error)
	GetServiceBrokerLabels(serviceBrokerName string) (map[string]types.NullString, v7action.Warnings, error)
	GetServiceBrokers() ([]resources.ServiceBroker, v7action.Warnings, error)
	GetServiceKeyByServiceInstanceAndName(serviceInstanceName, serviceKeyName, spaceGUID string) (resources.ServiceCredentialBinding, v7action.Warnings, error)
	GetServiceKeyDetailsByServiceInstanceAndName(serviceInstanceName, serviceKeyName, spaceGUID string) (resources.ServiceCredentialBindingDetails, v7action.Warnings, error)
	GetServiceInstanceByNameAndSpace(serviceInstanceName, spaceGUID string) (resources.ServiceInstance, v7action.Warnings, error)
	GetServiceInstanceDetails(serviceInstanceName, spaceGUID string, omitApps bool) (v7action.ServiceInstanceDetails, v7action.Warnings, error)
	GetServiceInstanceParameters(serviceInstanceName, spaceGUID string) (v7action.ServiceInstanceParameters, v7action.Warnings, error)
	GetServiceInstanceLabels(serviceInstanceName, spaceGUID string) (map[string]types.NullString, v7action.Warnings, error)
	GetServiceInstancesForSpace(spaceGUID string, omitApps bool) ([]v7action.ServiceInstance, v7action.Warnings, error)
	GetServiceKeysByServiceInstance(serviceInstanceName, spaceGUID string) ([]resources.ServiceCredentialBinding, v7action.Warnings, error)
	GetServiceOfferingLabels(serviceOfferingName, serviceBrokerName string) (map[string]types.NullString, v7action.Warnings, error)
	GetServicePlanLabels(servicePlanName, serviceOfferingName, serviceBrokerName string) (map[string]types.NullString, v7action.Warnings, error)
	GetServicePlanByNameOfferingAndBroker(servicePlanName, serviceOfferingName, serviceBrokerName string) (resources.ServicePlan, v7action.Warnings, error)
	GetSpaceByNameAndOrganization(spaceName string, orgGUID string) (resources.Space, v7action.Warnings, error)
	GetSpaceFeature(spaceName string, orgGUID string, feature string) (bool, v7action.Warnings, error)
	GetSpaceLabels(spaceName string, orgGUID string) (map[string]types.NullString, v7action.Warnings, error)
	GetSpaceQuotaByName(spaceQuotaName string, orgGUID string) (resources.SpaceQuota, v7action.Warnings, error)
	GetSpaceQuotasByOrgGUID(orgGUID string) ([]resources.SpaceQuota, v7action.Warnings, error)
	GetSpaceSummaryByNameAndOrganization(spaceName string, orgGUID string) (v7action.SpaceSummary, v7action.Warnings, error)
	GetSpaceUsersByRoleType(spaceGuid string) (map[constant.RoleType][]resources.User, v7action.Warnings, error)
	GetStackByName(stackName string) (resources.Stack, v7action.Warnings, error)
	GetStackLabels(stackName string) (map[string]types.NullString, v7action.Warnings, error)
	GetStacks(string) ([]resources.Stack, v7action.Warnings, error)
	GetStreamingLogsForApplicationByNameAndSpace(appName string, spaceGUID string, client sharedaction.LogCacheClient) (<-chan sharedaction.LogMessage, <-chan error, context.CancelFunc, v7action.Warnings, error)
	GetTaskBySequenceIDAndApplication(sequenceID int, appGUID string) (resources.Task, v7action.Warnings, error)
	GetUAAAPIVersion() (string, error)
	GetUnstagedNewestPackageGUID(appGuid string) (string, v7action.Warnings, error)
	GetUser(username, origin string) (resources.User, error)
	MakeCurlRequest(httpMethod string, path string, customHeaders []string, httpData string, failOnHTTPError bool) ([]byte, *http.Response, error)
	MapRoute(routeGUID string, appGUID string, destinationProtocol string) (v7action.Warnings, error)
	Marketplace(filter v7action.MarketplaceFilter) ([]v7action.ServiceOfferingWithPlans, v7action.Warnings, error)
	ParseAccessToken(accessToken string) (jwt.JWT, error)
	PollBuild(buildGUID string, appName string) (resources.Droplet, v7action.Warnings, error)
	PollPackage(pkg resources.Package) (resources.Package, v7action.Warnings, error)
	PollStart(app resources.Application, noWait bool, handleProcessStats func(string)) (v7action.Warnings, error)
	PollStartForRolling(app resources.Application, deploymentGUID string, noWait bool, handleProcessStats func(string)) (v7action.Warnings, error)
	PollUploadBuildpackJob(jobURL ccv3.JobURL) (v7action.Warnings, error)
	PrepareBuildpackBits(inputPath string, tmpDirPath string, downloader v7action.Downloader) (string, error)
	PurgeServiceInstance(serviceInstanceName, spaceGUID string) (v7action.Warnings, error)
	PurgeServiceOfferingByNameAndBroker(serviceOfferingName, serviceBrokerName string) (v7action.Warnings, error)
	RefreshAccessToken() (string, error)
	RenameApplicationByNameAndSpaceGUID(oldAppName, newAppName, spaceGUID string) (resources.Application, v7action.Warnings, error)
	RenameOrganization(oldOrgName, newOrgName string) (resources.Organization, v7action.Warnings, error)
	RenameServiceInstance(currentServiceInstanceName, spaceGUID, newServiceInstanceName string) (v7action.Warnings, error)
	RenameSpaceByNameAndOrganizationGUID(oldSpaceName, newSpaceName, orgGUID string) (resources.Space, v7action.Warnings, error)
	ResetOrganizationDefaultIsolationSegment(orgGUID string) (v7action.Warnings, error)
	ResetSpaceIsolationSegment(orgGUID string, spaceGUID string) (string, v7action.Warnings, error)
	ResourceMatch(resources []sharedaction.V3Resource) ([]sharedaction.V3Resource, v7action.Warnings, error)
	RestartApplication(appGUID string, noWait bool) (v7action.Warnings, error)
	RevokeAccessAndRefreshTokens() error
	RunTask(appGUID string, task resources.Task) (resources.Task, v7action.Warnings, error)
	ScaleProcessByApplication(appGUID string, process resources.Process) (v7action.Warnings, error)
	ScheduleTokenRefresh(func(time.Duration) <-chan time.Time, chan struct{}, chan struct{}) (<-chan error, error)
	SetApplicationDroplet(appGUID string, dropletGUID string) (v7action.Warnings, error)
	SetApplicationDropletByApplicationNameAndSpace(appName string, spaceGUID string, dropletGUID string) (v7action.Warnings, error)
	SetApplicationManifest(appGUID string, rawManifest []byte) (v7action.Warnings, error)
	SetApplicationProcessHealthCheckTypeByNameAndSpace(appName string, spaceGUID string, healthCheckType constant.HealthCheckType, httpEndpoint string, processType string, invocationTimeout int64) (resources.Application, v7action.Warnings, error)
	SetEnvironmentVariableByApplicationNameAndSpace(appName string, spaceGUID string, envPair v7action.EnvironmentVariablePair) (v7action.Warnings, error)
	SetEnvironmentVariableGroup(group constant.EnvironmentVariableGroupName, envVars resources.EnvironmentVariables) (v7action.Warnings, error)
	SetOrganizationDefaultIsolationSegment(orgGUID string, isoSegGUID string) (v7action.Warnings, error)
	SetSpaceManifest(spaceGUID string, rawManifest []byte) (v7action.Warnings, error)
	SetTarget(settings v7action.TargetSettings) (v7action.Warnings, error)
	SharePrivateDomain(domainName string, orgName string) (v7action.Warnings, error)
	ShareServiceInstanceToSpaceAndOrg(serviceInstanceName, targetedSpaceGUID, targetedOrgGUID string, sharedToDetails v7action.ServiceInstanceSharingParams) (v7action.Warnings, error)
	StageApplicationPackage(pkgGUID string) (resources.Build, v7action.Warnings, error)
	StagePackage(packageGUID, appName, spaceGUID string) (<-chan resources.Droplet, <-chan v7action.Warnings, <-chan error)
	StartApplication(appGUID string) (v7action.Warnings, error)
	StopApplication(appGUID string) (v7action.Warnings, error)
	TerminateTask(taskGUID string) (resources.Task, v7action.Warnings, error)
	UnbindSecurityGroup(securityGroupName string, orgGUID string, spaceGUID string, lifecycle constant.SecurityGroupLifecycle) (v7action.Warnings, error)
	UnmapRoute(routeGUID string, destinationGUID string) (v7action.Warnings, error)
	UnsetEnvironmentVariableByApplicationNameAndSpace(appName string, spaceGUID string, EnvironmentVariableName string) (v7action.Warnings, error)
	UnsetSpaceQuota(spaceQuotaName, spaceName, orgGUID string) (v7action.Warnings, error)
	UnsharePrivateDomain(domainName string, orgName string) (v7action.Warnings, error)
	UnshareServiceInstanceFromSpaceAndOrg(serviceInstanceName, targetedSpaceGUID, targetedOrgGUID string, unshareFromDetails v7action.ServiceInstanceSharingParams) (v7action.Warnings, error)
	UpdateAppFeature(app resources.Application, enabled bool, featureName string) (v7action.Warnings, error)
	UpdateApplication(app resources.Application) (resources.Application, v7action.Warnings, error)
	UpdateApplicationLabelsByApplicationName(string, string, map[string]types.NullString) (v7action.Warnings, error)
	UpdateBuildpackByNameAndStack(buildpackName string, buildpackStack string, buildpack resources.Buildpack) (resources.Buildpack, v7action.Warnings, error)
	UpdateBuildpackLabelsByBuildpackNameAndStack(string, string, map[string]types.NullString) (v7action.Warnings, error)
	UpdateDestination(string, string, string) (v7action.Warnings, error)
	UpdateDomainLabelsByDomainName(string, map[string]types.NullString) (v7action.Warnings, error)
	UpdateManagedServiceInstance(params v7action.UpdateManagedServiceInstanceParams) (chan v7action.PollJobEvent, v7action.Warnings, error)
	UpgradeManagedServiceInstance(serviceInstanceName, spaceGUID string) (v7action.Warnings, error)
	UpdateOrganizationLabelsByOrganizationName(string, map[string]types.NullString) (v7action.Warnings, error)
	UpdateOrganizationQuota(quotaName string, newName string, limits v7action.QuotaLimits) (v7action.Warnings, error)
	UpdateProcessByTypeAndApplication(processType string, appGUID string, updatedProcess resources.Process) (v7action.Warnings, error)
	UpdateRouteLabels(string, string, map[string]types.NullString) (v7action.Warnings, error)
	UpdateSecurityGroup(name, filePath string) (v7action.Warnings, error)
	UpdateSecurityGroupGloballyEnabled(securityGroupName string, lifecycle constant.SecurityGroupLifecycle, enabled bool) (v7action.Warnings, error)
	UpdateServiceBroker(serviceBrokerGUID string, model resources.ServiceBroker) (v7action.Warnings, error)
	UpdateServiceBrokerLabelsByServiceBrokerName(string, map[string]types.NullString) (v7action.Warnings, error)
	UpdateServiceInstanceLabels(serviceInstanceName, spaceGUID string, labels map[string]types.NullString) (v7action.Warnings, error)
	UpdateServiceOfferingLabels(serviceOfferingName string, serviceBrokerName string, labels map[string]types.NullString) (v7action.Warnings, error)
	UpdateServicePlanLabels(servicePlanName string, serviceOfferingName string, serviceBrokerName string, labels map[string]types.NullString) (v7action.Warnings, error)
	UpdateSpaceFeature(spaceName string, orgGUID string, enableds bool, feature string) (v7action.Warnings, error)
	UpdateSpaceLabelsBySpaceName(string, string, map[string]types.NullString) (v7action.Warnings, error)
	UpdateSpaceQuota(quotaName, orgGUID, newName string, limits v7action.QuotaLimits) (v7action.Warnings, error)
	UpdateStackLabelsByStackName(string, map[string]types.NullString) (v7action.Warnings, error)
	UpdateUserPassword(userGUID string, oldPassword string, newPassword string) error
	UpdateUserProvidedServiceInstance(serviceInstanceName, spaceGUID string, serviceInstanceUpdates resources.ServiceInstance) (v7action.Warnings, error)
	UploadBitsPackage(pkg resources.Package, matchedResources []sharedaction.V3Resource, newResources io.Reader, newResourcesLength int64) (resources.Package, v7action.Warnings, error)
	UploadBuildpack(guid string, pathToBuildpackBits string, progressBar v7action.SimpleProgressBar) (ccv3.JobURL, v7action.Warnings, error)
	UploadDroplet(dropletGUID string, dropletPath string, progressReader io.Reader, fileSize int64) (v7action.Warnings, error)
}