package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resourceCorpCloudWAFCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCorpCloudWAFCertificateCreate,
		Read:   resourceCorpCloudWAFCertificateRead,
		Update: resourceCorpCloudWAFCertificateUpdate,
		Delete: resourceCorpCloudWAFCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Friendly name to identify a CloudWAF certificate",
				Required:    true,
			},
			"certificate_body": {
				Type:             schema.TypeString,
				Description:      "Body of the certificate in PEM format",
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressEquivalentTrimSpaceDiffs,
			},
			"certificate_chain": {
				Type:             schema.TypeString,
				Description:      "Certificate chain in PEM format",
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentTrimSpaceDiffs,
			},
			"private_key": {
				Type:        schema.TypeString,
				Description: "Private key of the certificate in PEM format - must be unencrypted",
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
			},
			"common_name": {
				Type:        schema.TypeString,
				Description: "Common name of the uploaded certificate",
				Computed:    true,
			},
			"expires_at": {
				Type:        schema.TypeString,
				Description: "TimeStamp for when certificate expires in RFC3339 date time format",
				Computed:    true,
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Description: "SHA1 fingerprint of the certififcate",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: `Current status of the certificate - could be one of "unknown", "active", "pendingverification", "expired", "error"`,
				Computed:    true,
			},
			"subject_alternative_names": {
				Type:        schema.TypeSet,
				Description: "Subject alternative names from the uploaded certificate",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCorpCloudWAFCertificateCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	resp, err := sc.UploadCloudWAFCertificate(pm.Corp, sigsci.UploadCloudWAFCertificateBody{
		CloudWAFCertificateBase: sigsci.CloudWAFCertificateBase{
			Name:             d.Get("name").(string),
			CertificateBody:  strings.TrimSpace(d.Get("certificate_body").(string)),
			CertificateChain: strings.TrimSpace(d.Get("certificate_chain").(string)),
		},
		PrivateKey: strings.TrimSpace(d.Get("private_key").(string)),
	})
	if err != nil {
		return err
	}

	d.SetId(resp.ID)

	return resourceCorpCloudWAFCertificateRead(d, m)
}

func resourceCorpCloudWAFCertificateRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	cwaf, err := sc.GetCloudWAFCertificate(pm.Corp, d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(d.Id())
	err = d.Set("name", cwaf.Name)
	if err != nil {
		return err
	}
	err = d.Set("certificate_body", strings.TrimSpace(cwaf.CertificateBody))
	if err != nil {
		return err
	}
	err = d.Set("certificate_chain", strings.TrimSpace(cwaf.CertificateChain))
	if err != nil {
		return err
	}
	err = d.Set("common_name", cwaf.CommonName)
	if err != nil {
		return err
	}
	err = d.Set("expires_at", cwaf.ExpiresAt)
	if err != nil {
		return err
	}
	err = d.Set("fingerprint", cwaf.Fingerprint)
	if err != nil {
		return err
	}
	err = d.Set("status", cwaf.Status)
	if err != nil {
		return err
	}
	err = d.Set("subject_alternative_names", flattenStringArray(cwaf.SubjectAlternativeNames))
	if err != nil {
		return err
	}

	return nil
}

func resourceCorpCloudWAFCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	_, err := sc.UpdateCloudWAFCertificate(pm.Corp, d.Id(), sigsci.UpdateCloudWAFCertificateBody{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return nil
	}

	return resourceCorpCloudWAFCertificateRead(d, m)
}

func resourceCorpCloudWAFCertificateDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	err := sc.DeleteCloudWAFCertificate(pm.Corp, d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
