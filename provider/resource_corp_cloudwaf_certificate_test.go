package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccResourceCorpCloudWAFCertificateCRUD_basic(t *testing.T) {
	t.Parallel()

	resourceName := "sigsci_corp_cloudwaf_certificate.test_cloudwaf_certificate"
	certificateName := "Cloud WAF Certificate by SigSci Terraform provider test"
	certificateBody := `-----BEGIN CERTIFICATE-----
MIICzDCCAbQCCQDV2NzCr6aPbDANBgkqhkiG9w0BAQsFADAoMQwwCgYDVQQLDAN3
ZWIxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0yMjA5MDMxNzE5MTVaFw0z
MjA4MzExNzE5MTVaMCgxDDAKBgNVBAsMA3dlYjEYMBYGA1UEAwwPd3d3LmV4YW1w
bGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3xLllGg3Kl0S
W16pOwH4VXyGTyByWk3gShoSnEXqWYwoGDF18YhGFFMYLUBNfbDG/jy8MLiKY20R
BhV5hpObd4ZQq4PIlTl4ZNKy07CUPX/AufdbQzrFCQy96lXBVjo6gR10TD+F/CjC
tOkM83dxtZoSPzH86eHteos41+apjgpfvVai3vkBNuZeeoxuERkxuGsfpcK2qWTg
ZFuncrWt6Plvlu70qGEIPtiFiPfQ8Rs2mdzKJEBC8nb4nqSWxIbY9Z87yS3X3C/A
0xr0W4YxOLyCN94qr+Cc3Zl6DvjOv3LWAfv4qFXApWD9f8ynAzjojqfnXtavV6+D
1SOdvMcTVQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQDG9lhmTU9EzE/B7hAhNPj2
LHlRPj8ASnSpxBzWn6Acyu9hHJbGhkR8BnRPPjihH8kv+zRaRVYxhG2sb99qg168
rPKWMbZI6ZvCQGKNjLpwUARwPOKeZ8zF+qyzxdpM9mMyzx9SI1QXDirA0BsUbAjm
RfioqCdT54F8gFrH1+AnUX4Kf2euTS65bHRgegDiIsrAmwcRrzC8ev1SDiPMUkyC
goD1A0LHXLN1LMTs6qBXIkbCjYNRkPZBRagEu68CkwjT5H4vBIl39+Lcvo2WCBSG
j4LzHGcDief95tMLhz0f5g0geV3ytrld5NSw0g1sEYJlDe8NA/aDi4gVviOt3Z5A
-----END CERTIFICATE-----`
	privateKey := `-----BEGIN PRIVATE KEY-----
MIIEwAIBADANBgkqhkiG9w0BAQEFAASCBKowggSmAgEAAoIBAQDfEuWUaDcqXRJb
Xqk7AfhVfIZPIHJaTeBKGhKcRepZjCgYMXXxiEYUUxgtQE19sMb+PLwwuIpjbREG
FXmGk5t3hlCrg8iVOXhk0rLTsJQ9f8C591tDOsUJDL3qVcFWOjqBHXRMP4X8KMK0
6Qzzd3G1mhI/Mfzp4e16izjX5qmOCl+9VqLe+QE25l56jG4RGTG4ax+lwrapZOBk
W6dyta3o+W+W7vSoYQg+2IWI99DxGzaZ3MokQELydviepJbEhtj1nzvJLdfcL8DT
GvRbhjE4vII33iqv4JzdmXoO+M6/ctYB+/ioVcClYP1/zKcDOOiOp+de1q9Xr4PV
I528xxNVAgMBAAECggEBAKPsP/aJipg/4oBwFE2/SdyP8CZvQnjnpzzs4eYiXm7F
VqVIm1INAOpokWiXSxpk8CXdPbFTuqYLfKoK182z5Fe1xMv0wE4f+D+msTBsHtL+
cQJ3KYJCyo225kwwDi2uBlXg7hglyfCdh07nvtOeX1nCyUvVEPRRSHB3pCLLZqdv
ysCCL4Sowuebcpec3w3nCuMTg+L1nxdk25C53EjsYGqMQgq0YX/CTo+M2X2Infac
3Ig5bkQohaOz724L7mc63UaT9m36vgEUZfwndxxVUBHzxQ/tqr/O4XEKmSdFrExh
yw+YFyP43WcE0m1lc2hYHKEwTM1QF1nVY60fyPJ0zCECgYEA9ljzLd8OjXKbhK7S
HZWZyR1+Lo4HnrBsLgCJVIuGuWi6gDi8hsarYpTe8RVhgqBaudMpK+Eup76WuMBt
F9RVpzkNBE7T+p1jkmCWOIkCpvf2/IoqlZ1ao5qu1fhK5HpnsYCTmUHhJmid3da9
eT34oR2WnvJ9s6fvqC1C6URtsE0CgYEA59B+pvJxvV8klPDfIgN31bcR90srHJTJ
ST3lnVMJNND6PJ3D96YJSHV0EIotqyXlv/droqerhNJ+gNOHvrXMKTJ7udCDQKfW
6IFiWaRW1SkqUEXlv59JI6Ip3B9Rfi3w3rPhNLAi/7GjSm33C+OgsMVq9ZBTMjR6
qnS872CwsykCgYEAq22+3D8K+3ezraOSaDAA8qlpc7A2sUGIJoMNDh6CRGgS0MOq
vgdmoJWEhzQfxS0dtY6yaeyr8ON6M1sFD74dVN8opcTNUutPrT81imYdyF9qKtdj
RvZXat5rqE6+nzxnCGi3TcFAkt/ea8/RzptHd6cFd9q7itfkuJ22oGmUA0kCgYEA
zmvUO+kr6wtr0czjhLA952rLbr/ateqviq65ZmxoiEWGbq+1rzKElacxIQFKRVrL
yTMS/5X6n52o1CKIgAP2tsCjeAT6u3o5XnTIFTbHs6yiZzS2rvmx8S8Xw1GICanz
EPxwj7BAmhueYkqlcErT7lT9N4m667vbdynYi/g3oHECgYEA5dCjPFl94ZGStHF2
4+cI1NCbVTnQApFI1+Vmd4lED619KSnpk/77TaYTh2I8gsK3OZP4crvee5aL41TJ
Y20XmAq/Dvv3g7QR97ND/AghEU8nnpZo1fgHzCcyZSkaBzMFeYdsfaGDncM56I0B
65yblwYq9Vyzy3hBFY6XGaFdnZ0=
-----END PRIVATE KEY-----`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCorpCloudWAFCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "sigsci_corp_cloudwaf_certificate" "test_cloudwaf_certificate"{
					name = "%s"
					certificate_body = <<CERT
%s
CERT
					private_key = <<PRIVATEKEY
%s
PRIVATEKEY
					lifecycle {
						ignore_changes = [
							private_key
						]
					}
				}`, certificateName, certificateBody, privateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", certificateName),
					resource.TestCheckResourceAttr(resourceName, "certificate_body", certificateBody),
					resource.TestCheckResourceAttr(resourceName, "certificate_chain", ""),
					resource.TestCheckResourceAttr(resourceName, "common_name", "www.example.com"),
					resource.TestCheckResourceAttr(resourceName, "expires_at", "2032-08-31T17:19:15Z"),
					resource.TestCheckResourceAttr(resourceName, "fingerprint", "d3d246a79291ce3448f13b99d34d09066861c71a"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccImportStateCheckFunction(1),
			},
		},
	})
}

func TestAccResourceCorpCloudWAFCertificateCRUD_SAN(t *testing.T) {
	t.Parallel()

	resourceNameSAN := "sigsci_corp_cloudwaf_certificate.test_cloudwaf_certificate_san"
	certificateNameSAN := "Cloud WAF Certificate by SigSci Terraform provider test with SANs"
	certificateBodySAN := `-----BEGIN CERTIFICATE-----
MIIEozCCA4ugAwIBAgIJAIy4L1oqVlDVMA0GCSqGSIb3DQEBCwUAMIG5MQswCQYD
VQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5j
aXNjbzEdMBsGA1UECgwURXhhbXBsZSBPcmdhbml6YXRpb24xJDAiBgNVBAsMG0V4
YW1wbGUgT3JnYW5pemF0aW9uYWwgVW5pdDEUMBIGA1UEAwwLZXhhbXBsZS5jb20x
IjAgBgkqhkiG9w0BCQEWE2V4YW1wbGVAZXhhbXBsZS5jb20wHhcNMjIwOTE4MjI0
OTEwWhcNMzIwOTE1MjI0OTEwWjCBuTELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNh
bGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xHTAbBgNVBAoMFEV4YW1w
bGUgT3JnYW5pemF0aW9uMSQwIgYDVQQLDBtFeGFtcGxlIE9yZ2FuaXphdGlvbmFs
IFVuaXQxFDASBgNVBAMMC2V4YW1wbGUuY29tMSIwIAYJKoZIhvcNAQkBFhNleGFt
cGxlQGV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
xJXwW8I9MRC35ONpki/yNK1lPsPS36btinsQ6qt6cFq5+9Hy61a+pUwUESwPVXxe
PinMEq10J/Wm2AP3zOUQY5bHBwv3oGOCQTeZTIl9UimEOu96/6Q5belnQBTg0iSd
WMbz/xO8asd0lsL1wdNVQ3nDkN6HQvypuE2t88O0Sz9k9+M4nKYgGyJqdOW457Ko
HRbZxicRRF+9k8yGxD1Yv5i3AozHqxLdZ8/NOScVCW2pewwxq/iPU8YguD8gxiDk
memD7hKY7HuK+YiG4ar8jchupTt6TYJoq6msumxfcCepe/bHRMONzcwDfP8U/Dn8
C2gzk86YuV3esc8FVadqKwIDAQABo4GrMIGoMHkGA1UdEQRyMHCCD3d3dy5leGFt
cGxlLmNvbYIPZnRwLmV4YW1wbGUuY29tghBtYWlsLmV4YW1wbGUuY29tggtleGFt
cGxlLm5ldIIPd3d3LmV4YW1wbGUubmV0ggtleGFtcGxlLm9yZ4IPd3d3LmV4YW1w
bGUub3JnMAkGA1UdEwQCMAAwCwYDVR0PBAQDAgTwMBMGA1UdJQQMMAoGCCsGAQUF
BwMBMA0GCSqGSIb3DQEBCwUAA4IBAQBucpfz4j1hCKxaNzlIWr8/gn5r11tdKcg4
lXk8PXD+JOoxWTeSArVLIblaE5v3KFWN7fpGECviN0e6cbb6qrAtpekkGl7bdFEg
fwl+Qruiqw99WdIru0lu6ZDvRNw8OZHUAHSK1mFC1bNPQeLHWpVvLyfgukxmjO2E
VseStN/Ld1ZO/P9CQpeMw8/GagUeVhP7oq9t0N3r9f8iPvqPU+hgFnH5FcEjIosf
Mo4XLx5DAViHYDs+k0Zv/bnwIb6b10DvjzHo0RXYX2p74pz+MDbyVfPkiyv+rIXY
67YJjaYyteLj+UP1wVKa8qmZBg8zSbZqpcU2d8CDs6aoiVKpHqg/
-----END CERTIFICATE-----`
	privateKeySAN := `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDElfBbwj0xELfk
42mSL/I0rWU+w9Lfpu2KexDqq3pwWrn70fLrVr6lTBQRLA9VfF4+KcwSrXQn9abY
A/fM5RBjlscHC/egY4JBN5lMiX1SKYQ673r/pDlt6WdAFODSJJ1YxvP/E7xqx3SW
wvXB01VDecOQ3odC/Km4Ta3zw7RLP2T34zicpiAbImp05bjnsqgdFtnGJxFEX72T
zIbEPVi/mLcCjMerEt1nz805JxUJbal7DDGr+I9TxiC4PyDGIOSZ6YPuEpjse4r5
iIbhqvyNyG6lO3pNgmirqay6bF9wJ6l79sdEw43NzAN8/xT8OfwLaDOTzpi5Xd6x
zwVVp2orAgMBAAECggEAKiN5wjGArGPJB2c32f4tDN2eNjYDna1KfcSje6oGNM89
zpzSVV/ivcvxAT1QjCJ8kRakh9xmaapeeS9gjqsLOE25m+kUy2yJHzGrypwuIM6F
aZyr4OBy7vx5BWN0TZdLoxwCcUrpuHnIpAhmZYXHYQ9YvFT26YT/XGJKR1ZL71Tm
2fR7HiO3wHCqdkNr2p/PjORsLJzeuBwX5BHCA27jJrxclsb/fvsCy4zNRjU/VapN
rvrSs+EbrB9U5xNj+NqhIJpw8bCVZHfeAK4vxNLyXU+VgboxKM+VzPx4RCcJVPGW
rcE+KhIA1HQvWvIfIjgu0igdTnuUWamWGeQVFUgwAQKBgQD2CimkNLEw68Ik7FwD
ewEALrDmioPWzOUUFxjuhZMqR4VQIefriEBTASYyV0LRgZlQZ9N2P9Y3tz0TA7gk
8cmQ/cXHuklgiqv65bxYWbF7Z4VT5Ofnj5NgB6y6Qyo8724R7tZP52FWum+oqLK0
9lyzlk9pjTwBB/avJG9GZ1rEQQKBgQDMi0I8AxgbbZ54I/v3WvyyMkMWfojHycwC
tIHqqj0DgWo6ywCXN1Uh5hYzUc8rbD8Aj0ayaEY2tt4sfroZWL18hE+7+zF5T33g
r9ugexiSWPXJYPzsfcY9L5wMsislfebdilu6zaK+1JTZrcixStVkQodN/Mnwht8G
qzhm/0KjawKBgFsWohoH4/3Pmr3ev6YOOO2fW6DOcUbp7nmEn5dW3ogNmH51Pw6F
EANq7oA+rB8yUtdgyPoDYkSYU0Uh4F/VICHMwhdSkW3riQZHXXZ8JmggiEzp9y1U
i2RHExyWVyHeJqr++Fr2t6PLPCF960Nx3hoisN3MCwX6s8pdu9Cd3Q3BAoGBAKWa
SfHUf3wVcCvM1n8Zx2VulCpuH8fBdc7q3hRj6CoiaSNYoA3N2rsrUeYS4ixB43BM
j+x5x/8cZxyXLYy/8AoUBYoogZG57iwvtR1lDCvQoo58W8oMuqdnGGyfA5fDK1tK
XaIMQytFaY4jyUzhTYty4aEefVCjoYYAshWRrR4pAoGBAOCJRQV5jT/6Zk+vjdtu
kxvINjcwcyux6QVomkSP60o4BQI6+KNVa8aguQUYdVsP7WzJjj9ivJWOUAL2S4Q9
zNwDGTj0Iz8JYyqptGl9laPt7DXZkfmROl0OEP4SS5CYTf0r2ZCwcQ30lVEJteRQ
RUh/KHYQRK0Cvo1agUHXkXvo
-----END PRIVATE KEY-----`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCorpCloudWAFCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "sigsci_corp_cloudwaf_certificate" "test_cloudwaf_certificate_san"{
					name = "%s"
					certificate_body = <<CERT
%s
CERT
					certificate_chain = <<CHAIN
%s
CHAIN
					private_key = <<PRIVATEKEY
%s
PRIVATEKEY
					lifecycle {
						ignore_changes = [
							private_key
						]
					}
				}`, certificateNameSAN, certificateBodySAN, certificateBodySAN, privateKeySAN),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceNameSAN, "name", certificateNameSAN),
					resource.TestCheckResourceAttr(resourceNameSAN, "certificate_body", certificateBodySAN),
					resource.TestCheckResourceAttr(resourceNameSAN, "certificate_chain", certificateBodySAN),
					resource.TestCheckResourceAttr(resourceNameSAN, "common_name", "example.com"),
					resource.TestCheckResourceAttr(resourceNameSAN, "expires_at", "2032-09-15T22:49:10Z"),
					resource.TestCheckResourceAttr(resourceNameSAN, "fingerprint", "78035e30621e0bbe92b9fce269a88a753ffe1154"),
					resource.TestCheckResourceAttr(resourceNameSAN, "status", "active"),
					resource.TestCheckResourceAttr(resourceNameSAN, "subject_alternative_names.#", "7"),
					resource.TestCheckResourceAttr(resourceNameSAN, "subject_alternative_names.1066711791", "www.example.com"),
					resource.TestCheckResourceAttr(resourceNameSAN, "subject_alternative_names.1320559578", "mail.example.com"),
					resource.TestCheckResourceAttr(resourceNameSAN, "subject_alternative_names.1533628156", "www.example.net"),
					resource.TestCheckResourceAttr(resourceNameSAN, "subject_alternative_names.1975881950", "example.org"),
					resource.TestCheckResourceAttr(resourceNameSAN, "subject_alternative_names.2605964798", "www.example.org"),
					resource.TestCheckResourceAttr(resourceNameSAN, "subject_alternative_names.3053388764", "example.net"),
					resource.TestCheckResourceAttr(resourceNameSAN, "subject_alternative_names.57398617", "ftp.example.com"),
					resource.TestCheckResourceAttrSet(resourceNameSAN, "created_by"),
					resource.TestCheckResourceAttrSet(resourceNameSAN, "created"),
					resource.TestCheckResourceAttrSet(resourceNameSAN, "updated_by"),
					resource.TestCheckResourceAttrSet(resourceNameSAN, "updated_at"),
				),
			},
			{
				ResourceName:      resourceNameSAN,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccImportStateCheckFunction(1),
			},
		},
	})
}

func testAccCheckCorpCloudWAFCertificateDestroy(s *terraform.State) error {
	pm := testAccProvider.Meta().(providerMetadata)
	sc := pm.Client

	resourceType := "sigsci_corp_cloudwaf_certificate"
	for _, resource := range s.RootModule().Resources {
		if resource.Type != resourceType {
			continue
		}
		readResp, err := sc.GetCloudWAFCertificate(pm.Corp, resource.Primary.Attributes["id"])
		if err == nil {
			return fmt.Errorf("%s %#v still exists", resourceType, readResp)
		}
	}
	return nil
}
