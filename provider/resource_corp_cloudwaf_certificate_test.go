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
MIIDvDCCAqQCCQDj4MMBbF4gWTANBgkqhkiG9w0BAQsFADCBnzELMAkGA1UEBhMC
VVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28x
FDASBgNVBAoMC0V4YW1wbGUgT3JnMRMwEQYDVQQLDApFeGFtcGxlIE9VMRQwEgYD
VQQDDAtleGFtcGxlLmNvbTEiMCAGCSqGSIb3DQEJARYTZXhhbXBsZUBleGFtcGxl
LmNvbTAeFw0yMjA5MjQwMDAxMzlaFw0yMjEwMjQwMDAxMzlaMIGfMQswCQYDVQQG
EwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNj
bzEUMBIGA1UECgwLRXhhbXBsZSBPcmcxEzARBgNVBAsMCkV4YW1wbGUgT1UxFDAS
BgNVBAMMC2V4YW1wbGUuY29tMSIwIAYJKoZIhvcNAQkBFhNleGFtcGxlQGV4YW1w
bGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnuLmCj29AAyW
fwkErscoHZ4V20DUF3DfsOPUKd5MtXSEQqFnI1fIDkDJC4KL2poTFQoZ4TBuGjcw
lmgAbQggUF4V/UobvEQaqiWeTMUh4YW0szNVvZEVzpwcE2M71KNPx72AuoHs+sgv
FszwAGIpZw1teAhwDqMPscHm/KsK4dxnOkAD+FdMVM5oYCQmf9sPS45FdYEZHueA
554QYObrh43G5tJcte9S9fESgjWfg951ESVcFCHWEG6XQwT9hux9KplgsZJfmgaf
LUaFlnuM8dldi6H9TPL+o4PRRdz8dO/NGD3IkmxncxPt6ATpPRUfgxUi8zr1wEvc
8/oo4C1VVwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQArwwMv9SYSQne12zNEEm7k
77toN9Ya+36mQOFFvNA6Vajd2b4EvlKzbnJox5OkZE6xcE1an4yhKyYYOpqApGr5
mLbdzrUHTqY9IeGrpOuBd2LXrpKtgBR27+lxXzHXd/CWIFn9YVr5IcNaYwCsYgMr
sskUi+lDJWXmkiYoRYKxvR5Ug0NYzEyxj8ZmrGYHk502BxjeW2bFHdXqAZGqoC0O
XjJ53nNvxZHhaIGK9WDDuzem+6b5r4mQVK76BLfwJ/JB2oO3BWYL5cxb6MX/1DXK
ctwO5KlIK0rx/s8nSQBB+QosaXMDP0DQGqEHWT7CQTuT1gNW8ktvzGPQrpjK7JXJ
-----END CERTIFICATE-----`
	privateKey := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCe4uYKPb0ADJZ/
CQSuxygdnhXbQNQXcN+w49Qp3ky1dIRCoWcjV8gOQMkLgovamhMVChnhMG4aNzCW
aABtCCBQXhX9Shu8RBqqJZ5MxSHhhbSzM1W9kRXOnBwTYzvUo0/HvYC6gez6yC8W
zPAAYilnDW14CHAOow+xweb8qwrh3Gc6QAP4V0xUzmhgJCZ/2w9LjkV1gRke54Dn
nhBg5uuHjcbm0ly171L18RKCNZ+D3nURJVwUIdYQbpdDBP2G7H0qmWCxkl+aBp8t
RoWWe4zx2V2Lof1M8v6jg9FF3Px0780YPciSbGdzE+3oBOk9FR+DFSLzOvXAS9zz
+ijgLVVXAgMBAAECggEBAImggSLd15jzTmk7ppK+cEE3bjc9MHodi6XtsxmRNWD4
TJhqtqwmnWO7Omp96iawz1aqKUCmcrjClZOzAqtvHo5+8Q015FBvrak0bKqTF4YC
C0Quc1aBFiKhlrA0hN7rl2+s9pSXdm7EeAWH/1xVqwdY2jnfFTGYjT+sdijm/8Yj
hpvlcyqUC3jGrc9hQHDhqwzlVrP/dhYpQIdGwTRHpUXDqwNuXJQIzLhcE0ex+5MM
gWrgAi0M/Qwbn7CyEwdcapDjX6Bt8dVloroeODEYrsDClAXB+45BnRj5HsgYSvQJ
Sn3Xqa5sGxpwOPESOuzhX/Y4f/v2iqbVGw+D0AOcvgECgYEA0dNk+MqpylUrG4dY
uKCV+QADMIllQesw7e8hqKBAubmkfynbml8FMTr9e6GFXIS2Ujwsdg5QQM3b0JqS
qDYSk5EXnWy5zy61Zz9LlGoqmLwEI9NEPtUpSctN1VglhRmW2L9XKnEJ7dgbHl+e
AcMPBI2ownETAClHg8qKm5hMsZsCgYEAwdnQTAtM3gInT6ItxmUIWjKmDQP2tVY0
mIYpiPcoPnFe0IJCzwcWuBezdvgK3x0i2UY/HfRu+d93vMoiXt6ij+zab9ewXVHQ
PmytXelIJuQ4tA38L9HZGJW6yIzIaJT6laZiQStjTEN4lDeKNIpY4SEuqGzXaAl0
CHk+DovJ1PUCgYB82Q6kZloe1QxgRelJeeuijBpZv/brARlNCdN6NVgt6kLxkyNi
uBUr1NDMxi/G/ARL7Bf8asnftV2MwtxukDX/bf6iIfZxS3aOp3++IGmWFZFVC7j4
tfbqPLjkL52rk61I7Jjd3QKubb69FOG8ZKbD69I1V/iZSPaPeW195WIE7wKBgFN0
6deDWfGOrcv7/4cVgjYK7jBWT4WceoJb6E/eUIYpmu9b1VV6MM7K7Wm/ujZ6PcGb
G5tS2+BZ1BwETi3X3dbm2tgh3P0gNu5ZLX5r67NKuBrUlokj6DpMZCDpc3KLCSMa
gdyayGJR/fyZuLeMBF3QQl0ils5km372a8ApcJhtAoGAfDup6P3pDp2Rjq3IRJod
NFJAJERVDnSYml2D9q/uitU51nqL1JcvAoWQJByMFOXuy8F97Olflw4MW/cjWnXJ
1IBEgPEGv5zlronLACOIynBo6/Iu8ULIZhx0BYDc4DXD++PxG2Dd0mbYD8vrlrUK
wOmYr1etSwJs1p3ESLrOscM=
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
				}`, certificateName, certificateBody, privateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", certificateName),
					resource.TestCheckResourceAttr(resourceName, "certificate_body", certificateBody),
					resource.TestCheckResourceAttr(resourceName, "certificate_chain", ""),
					resource.TestCheckResourceAttr(resourceName, "common_name", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "expires_at", "2022-10-24T00:01:39Z"),
					resource.TestCheckResourceAttr(resourceName, "fingerprint", "56d4f213a9061be925445280a64865927df7f88f"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
				ImportStateCheck:        testAccImportStateCheckFunction(1),
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
				),
			},
			{
				ResourceName:            resourceNameSAN,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
				ImportStateCheck:        testAccImportStateCheckFunction(1),
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
