package provider

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceCorpCloudWAFInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceCorpCloudWAFInstanceCreate,
		Read:   resourceCorpCloudWAFInstanceRead,
		Update: resourceCorpCloudWAFInstanceUpdate,
		Delete: resourceCorpCloudWAFInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Friendly name to identify a CloudWAF instance.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Friendly description to identify a CloudWAF instance.",
				Required:    true,
			},
			"region": {
				Type:         schema.TypeString,
				Description:  `Region the CloudWAF Instance is being deployed to. (Supported region: "us-east-1", "us-west-1", "af-south-1", "ap-northeast-1", "ap-northeast-2", "ap-south-1", "ap-southeast-1", "ap-southeast-2", "ca-central-1", "eu-central-1", "eu-north-1", "eu-west-1", "eu-west-2", "eu-west-3", "sa-east-1", "us-east-2", "us-west-2").`,
				Required:     true,
				ValidateFunc: validateRegion,
			},
			"tls_min_version": {
				Type:        schema.TypeString,
				Description: `TLS minimum version. Versions Available: "1.0", "1.2".`,
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					if val == nil {
						return nil, nil
					}
					if existsInString(val.(string), "1.0", "1.2") {
						return nil, nil
					}
					return nil, []error{errors.New(`tlsMinVersion must be "1.0" or "1.2"`)}
				},
			},
			"use_uploaded_certificates": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"workspace_configs": {
				Type:        schema.TypeSet,
				Description: "Workspace Configs",
				Required:    true,
				MaxItems:    5,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_name": {
							Type:        schema.TypeString,
							Description: "Site name.",
							Required:    true,
						},
						"instance_location": {
							Type:        schema.TypeString,
							Description: `Set instance location to "direct" or "advanced".`,
							Required:    true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								if val == nil {
									return nil, nil
								}
								if existsInString(val.(string), "direct", "advanced") {
									return nil, nil
								}
								return nil, []error{fmt.Errorf(`received instance_location %q is invalid. should be "direct" or "advanced"`, val.(string))}
							},
						},
						"client_ip_header": {
							Type:        schema.TypeString,
							Description: `Specify the request header containing the client IP address, available when InstanceLocation is set to "advanced". Default: "X-Forwarded-For".`,
							Optional:    true,
						},
						"listener_protocols": {
							Type:        schema.TypeSet,
							Description: `Specify the protocol or protocols required. ex. ["http", "https"], ["https"].`,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"routes": {
							Type:        schema.TypeSet,
							Description: "Routes",
							Required:    true,
							MaxItems:    200,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "Route unique identifier.",
										Computed:    true,
									},
									"certificate_ids": {
										Type:        schema.TypeSet,
										Description: "List of certificate IDs in string associated with request URI or domains. IDs will be available in certificate GET request.",
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"connection_pooling": {
										Type:        schema.TypeBool,
										Description: "If enabled, this will allow open TCP connections to be reused (default: true)",
										Optional:    true,
										Default:     true,
									},
									"domains": {
										Type:        schema.TypeSet,
										Description: "List of domain or request URIs, up to 100 entries.",
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										MaxItems:    100,
									},
									"origin": {
										Type:        schema.TypeString,
										Description: "Origin server URI.",
										Required:    true,
									},
									"pass_host_header": {
										Type:        schema.TypeBool,
										Description: "Pass the client supplied host header through to the upstream (including the upstream TLS handshake for use with SNI and certificate validation). If using Heroku or Server Name Indications (SNI), this must be disabled (default: false).",
										Optional:    true,
										Default:     false,
									},
									"trust_proxy_headers": {
										Type:        schema.TypeBool,
										Description: "If true, will trust proxy headers coming into the agent. If false, will ignore and drop those headers (default: false)",
										Optional:    true,
										Default:     false,
									},
								},
							},
						},
					},
				},
			},
			"deployment": {
				Type:        schema.TypeList, // use TypeList to workaround SDK TypeMap limitation with only string value support: https://github.com/hashicorp/terraform-plugin-sdk/issues/62
				Description: "The sites primary Agent key",
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Description: "Current status of the deployment.",
							Computed:    true,
						},
						"message": {
							Type:        schema.TypeString,
							Description: "CloudWAF instance message.",
							Computed:    true,
						},
						"dns_entry": {
							Type:        schema.TypeString,
							Description: "CloudWAF instance's DNS Entry.",
							Computed:    true,
						},
						"egress_ips": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:        schema.TypeString,
										Description: "Egress IP address CloudWAF will be directing traffic to origin from.",
										Computed:    true,
									},
									"status": {
										Type:        schema.TypeString,
										Description: "EgressIP Status.",
										Computed:    true,
									},
									"updated_at": {
										Type:        schema.TypeString,
										Description: "When EgressIP was last updated on.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

const (
	cloudWAFInstanceErrGetNotFound      = "Not Found"
	cloudWAFInstanceStateDeprovisioning = "deprovisioning"
	cloudWAFInstanceStateDone           = "done"
	cloudWAFInstanceStatePending        = "pending"
)

func resourceCorpCloudWAFInstanceCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	cwaf, err := sc.CreateCloudWAFInstance(pm.Corp, sigsci.CloudWAFInstanceBody{
		Name:                    d.Get("name").(string),
		Description:             d.Get("description").(string),
		Region:                  d.Get("region").(string),
		TLSMinVersion:           d.Get("tls_min_version").(string),
		UseUploadedCertificates: d.Get("use_uploaded_certificates").(bool),
		WorkspaceConfigs:        expandCloudWAFInstanceWorkspaceConfigs(d.Get("workspace_configs").(*schema.Set)),
	})
	if err != nil {
		return err
	}

	stateChangeConf := cloudWAFInstanceStateChangeConf(d, m, cwaf.ID)
	_, err = stateChangeConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Cloud WAF instance id [%s] to be created with error [%s]", cwaf.ID, err.Error())
	}

	d.SetId(cwaf.ID)

	return resourceCorpCloudWAFInstanceRead(d, m)
}

func cloudWAFInstanceStateChangeConf(d *schema.ResourceData, m interface{}, cwafInstanceID string) *resource.StateChangeConf {
	pm := m.(providerMetadata)
	sc := pm.Client

	return &resource.StateChangeConf{
		Pending: []string{cloudWAFInstanceStatePending},
		Target:  []string{cloudWAFInstanceStateDone},
		Refresh: func() (interface{}, string, error) {
			cwaf, err := sc.GetCloudWAFInstance(pm.Corp, cwafInstanceID)
			if err != nil {
				return 0, "", err
			}
			return cwaf, cwaf.Deployment.Status, nil
		},
		Timeout:                   d.Timeout(schema.TimeoutCreate) - time.Minute,
		Delay:                     30 * time.Second,
		MinTimeout:                5 * time.Second,
		PollInterval:              30 * time.Second,
		ContinuousTargetOccurence: 1,
	}
}

func resourceCorpCloudWAFInstanceRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	cwaf, err := sc.GetCloudWAFInstance(pm.Corp, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(d.Id())
	err = d.Set("name", cwaf.Name)
	if err != nil {
		return err
	}
	err = d.Set("description", cwaf.Description)
	if err != nil {
		return err
	}
	err = d.Set("region", cwaf.Region)
	if err != nil {
		return err
	}
	err = d.Set("tls_min_version", cwaf.TLSMinVersion)
	if err != nil {
		return err
	}
	err = d.Set("use_uploaded_certificates", cwaf.UseUploadedCertificates)
	if err != nil {
		return err
	}
	err = d.Set("workspace_configs", flattenCloudWAFInstanceWorkspaceConfigs(cwaf.WorkspaceConfigs))
	if err != nil {
		return err
	}
	err = d.Set("deployment", flattenCloudWAFInstanceDeployment(cwaf.Deployment))
	if err != nil {
		return err
	}

	return nil
}

func resourceCorpCloudWAFInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	err := sc.UpdateCloudWAFInstance(pm.Corp, d.Id(), sigsci.CloudWAFInstanceBody{
		Name:                    d.Get("name").(string),
		Description:             d.Get("description").(string),
		Region:                  d.Get("region").(string),
		TLSMinVersion:           d.Get("tls_min_version").(string),
		UseUploadedCertificates: d.Get("use_uploaded_certificates").(bool),
		WorkspaceConfigs:        expandCloudWAFInstanceWorkspaceConfigs(d.Get("workspace_configs").(*schema.Set)),
	})
	if err != nil {
		return nil
	}

	stateChangeConf := cloudWAFInstanceStateChangeConf(d, m, d.Id())
	_, err = stateChangeConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Cloud WAF instance id [%s] to be updated with error [%s]", d.Id(), err.Error())
	}

	return resourceCorpCloudWAFInstanceRead(d, m)
}

func resourceCorpCloudWAFInstanceDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client
	cwafInstanceID := d.Id()

	err := sc.DeleteCloudWAFInstance(pm.Corp, cwafInstanceID)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutDelete)-time.Minute, func() *resource.RetryError {
		cwaf, err := sc.GetCloudWAFInstance(pm.Corp, cwafInstanceID)
		if err != nil && err.Error() == cloudWAFInstanceErrGetNotFound {
			d.SetId("")
			return nil
		} else if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error describing instance [%s]", err))
		} else if cwaf.Deployment.Status == cloudWAFInstanceStateDeprovisioning {
			return resource.RetryableError(fmt.Errorf("expected instance to be deleted but was in state [%s]", cwaf.Deployment.Status))
		} else {
			return resource.NonRetryableError(fmt.Errorf("expected instance to be deleted or in in pending state but was in state [%s]", cwaf.Deployment.Status))
		}
	})
}

func expandCloudWAFInstanceWorkspaceConfigs(routesResource *schema.Set) []sigsci.CloudWAFInstanceWorkspaceConfig {
	var configs []sigsci.CloudWAFInstanceWorkspaceConfig
	for _, genericElement := range routesResource.List() {
		castElement := genericElement.(map[string]interface{})
		configs = append(configs, sigsci.CloudWAFInstanceWorkspaceConfig{
			SiteName:          castElement["site_name"].(string),
			ClientIPHeader:    castElement["client_ip_header"].(string),
			InstanceLocation:  castElement["instance_location"].(string),
			ListenerProtocols: expandStringArray(castElement["listener_protocols"].(*schema.Set)),
			Routes:            expandCloudWAFInstanceRoutes(castElement["routes"].(*schema.Set)),
		})
	}
	return configs
}

func expandCloudWAFInstanceRoutes(routesResource *schema.Set) []sigsci.CloudWAFInstanceWorkspaceRoute {
	var routes []sigsci.CloudWAFInstanceWorkspaceRoute
	for _, genericElement := range routesResource.List() {
		castElement := genericElement.(map[string]interface{})
		routes = append(routes, sigsci.CloudWAFInstanceWorkspaceRoute{
			ID:                castElement["id"].(string),
			CertificateIDs:    expandStringArray(castElement["certificate_ids"].(*schema.Set)),
			ConnectionPooling: castElement["connection_pooling"].(bool),
			Domains:           expandStringArray(castElement["domains"].(*schema.Set)),
			Origin:            castElement["origin"].(string),
			PassHostHeader:    castElement["pass_host_header"].(bool),
			TrustProxyHeaders: castElement["trust_proxy_headers"].(bool),
		})
	}
	return routes
}

func flattenCloudWAFInstanceWorkspaceConfigs(configs []sigsci.CloudWAFInstanceWorkspaceConfig) []interface{} {
	var configsMap = make([]interface{}, len(configs))
	for i, config := range configs {
		configMap := map[string]interface{}{
			"site_name":          config.SiteName,
			"client_ip_header":   config.ClientIPHeader,
			"instance_location":  config.InstanceLocation,
			"listener_protocols": flattenStringArray(config.ListenerProtocols),
			"routes":             flattenCloudWAFInstanceRoutes(config.Routes),
		}
		configsMap[i] = configMap
	}
	return configsMap
}

func flattenCloudWAFInstanceRoutes(routes []sigsci.CloudWAFInstanceWorkspaceRoute) []interface{} {
	var routesMap = make([]interface{}, len(routes))
	for i, route := range routes {
		routeMap := map[string]interface{}{
			"id":                  route.ID,
			"certificate_ids":     flattenStringArray(route.CertificateIDs),
			"connection_pooling":  route.ConnectionPooling,
			"domains":             flattenStringArray(route.Domains),
			"origin":              route.Origin,
			"pass_host_header":    route.PassHostHeader,
			"trust_proxy_headers": route.TrustProxyHeaders,
		}
		routesMap[i] = routeMap
	}
	return routesMap
}

func flattenCloudWAFInstanceDeployment(deployment sigsci.CloudWAFInstanceDeployment) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"status":     deployment.Status,
			"message":    deployment.Message,
			"dns_entry":  deployment.DNSEntry,
			"egress_ips": flattenCloudWAFInstanceEgressIPs(deployment.EgressIPs),
		},
	}
}

func flattenCloudWAFInstanceEgressIPs(egressIPs []sigsci.CloudWAFInstanceEgressIP) []interface{} {
	var egressIPsSet = make([]interface{}, len(egressIPs))
	for i, egressIP := range egressIPs {
		egressIPsSet[i] = map[string]interface{}{
			"ip":         egressIP.IP,
			"status":     egressIP.Status,
			"updated_at": egressIP.UpdatedAt}
	}
	return egressIPsSet
}
