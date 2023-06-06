package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceCorpCloudWAFCertificateCRUD_basic(t *testing.T) {
	t.Parallel()

	resourceName := "sigsci_corp_cloudwaf_certificate.test_cloudwaf_certificate"
	certificateName := "Cloud WAF Certificate by SigSci Terraform provider test"
	certificateBody := `-----BEGIN CERTIFICATE-----
MIIDzjCCArYCCQD6uBPuCbaDuDANBgkqhkiG9w0BAQsFADCBqDELMAkGA1UEBhMC
VVMxEzARBgNVBAgMCkNhbGlmb3JuaWExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28x
HTAbBgNVBAoMFEV4YW1wbGUgT3JnYW5pemF0aW9uMRMwEQYDVQQLDApFeGFtcGxl
IE9VMRQwEgYDVQQDDAtleGFtcGxlLmNvbTEiMCAGCSqGSIb3DQEJARYTZXhhbXBs
ZUBleGFtcGxlLmNvbTAeFw0yMjA5MjQwMDE4MjRaFw0zMjA5MjEwMDE4MjRaMIGo
MQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2Fu
IEZyYW5jaXNjbzEdMBsGA1UECgwURXhhbXBsZSBPcmdhbml6YXRpb24xEzARBgNV
BAsMCkV4YW1wbGUgT1UxFDASBgNVBAMMC2V4YW1wbGUuY29tMSIwIAYJKoZIhvcN
AQkBFhNleGFtcGxlQGV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEAscvDb2j2s9bdiAIHbqRoM2qZBxdM4atSwAJQrXVe3pbne2KLZw53
kHpVtjaugfMKBnXueR1iilYu5eXtgNfrNHgq0X0+NToL/xtSgYthp89lxBYArUVy
kiM5gy8BqpApfAwQ5MMDgGflIV/mTCcCyNK3DwuOgO7oVp0V2zdtJhgvZ8e3qkuT
3dOxC27aUFYf/P88UILoc9YWRCkw2Gww/Zr908a/mgVBJ9v+/sKP3/yk8jzrRhL5
JsGWC5Gbv1gpkyzSjKyboYePvJJo5D6Fue9XZmzry3wepG1oUcLO6QpH+lTBfTjd
xHKA4sIza1J/RDBLgUBney1nMxLN8RzU5QIDAQABMA0GCSqGSIb3DQEBCwUAA4IB
AQAPRvwDkKTKCDQj5F4ZUTE9AIEs0w99KuXiWBGz3RmYl5zwZCrVWeOI+lPfCG0v
prMgh5ydUgUOqrs8S7MAkt8GaU5lb0MSKmz1jPgEEbLBp6VYv2UbrWlBz9JIxTLw
riPHNUzKb6SXk5wuoO8w7+GsBNI8fWPDSQqSWLlNsi0r4ReLxlM5WBNC10d3q2ia
jV6r8iMpiArwbJn4WSTlFuJ6crrjgbBVCFxxwoF1sHhwGg+5idxm2AHSzvENyFW4
UVhVTn9w3UPLMkEl7nAVzydpdMb/M/GLCV787BrQL35EtiCr9MSL9Gc8vR/9PzPP
QodC+xWXbig7xKLqZgQ/PbPt
-----END CERTIFICATE-----`
	privateKey := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCxy8NvaPaz1t2I
AgdupGgzapkHF0zhq1LAAlCtdV7elud7YotnDneQelW2Nq6B8woGde55HWKKVi7l
5e2A1+s0eCrRfT41Ogv/G1KBi2Gnz2XEFgCtRXKSIzmDLwGqkCl8DBDkwwOAZ+Uh
X+ZMJwLI0rcPC46A7uhWnRXbN20mGC9nx7eqS5Pd07ELbtpQVh/8/zxQguhz1hZE
KTDYbDD9mv3Txr+aBUEn2/7+wo/f/KTyPOtGEvkmwZYLkZu/WCmTLNKMrJuhh4+8
kmjkPoW571dmbOvLfB6kbWhRws7pCkf6VMF9ON3EcoDiwjNrUn9EMEuBQGd7LWcz
Es3xHNTlAgMBAAECggEAV0K/f52Pf0JUZd1BEo+ESL/nrTBFXni8W1qHiCqTzkFY
CRmbe5ABJJq2GIEL8uF6qSMWUMEYTPbxe4n2oAbY/F6B/WEvt+XuX11kiAoFetvy
gWOfH2t3SLwbDQR0F+c7RROS8wO3Yz0amt+7YuK+nhu1FqBAZ41Z4LCmOnogitIA
cCErKpHqCJbYT99eaXTt2QpXJNI8fItXaO4p8zfKxzBKybhsyu3tEerKWnhqz+25
Xr7OieYkZM3ryrIsVWZ299wH2D+gA9O+PbY2RJ0Vf3YVf8VdyYRJ5oCAaioNUfZw
HeGGlfClZzrpX1MjiNNfli05cLqX0iE2bIO+jvZkIQKBgQDadIMHfGQs8CPxbPqi
fZRuHdPedogM80f4E5RyKhTzEqTG6x0pfjkDr57rdqJEWhM3TjsMIiJRnEIomWar
2tyckmvkxuSioiWb/+HXJN5u6AtsMgQ4WLMm9HOO0ir5uSKQd0iQrCMAaGDNdV52
6eipkWLdhYAhW31bycv9cX04+QKBgQDQWlkTovazzNElJU4h1YQCq99VS1gmYb9U
HAzg3Jmu7WF4Oln+HxZMWqwR38vCuMHiwtCmsqGAEKK2ev6W6iq7iVJLKzTXV++a
612Mr+JohbHNL0bKlgTMt/i2TnmBWOlhL7xuIwduru1pQ4mM7Vh7Hv+CrVTGm4VZ
Khzq+vbCTQKBgEnprgO0ZLiHr8GZ29tqnfP8B5l3hWTMU4duKIXQEzKDFllvZ3iI
ioXiv+RvSUvTJjlKMNRUIER4mDHgZUq0THx1Vigb23PjZNI5a5I9mTzxKhw7eA4Q
hN0jTI4AMiY4K6exlE3O0DDtIAOkOIgHcH8e/9JvvwCKUgniZzCjW3kRAoGAeEId
tgLSyFbIxNrybP7zciNIBdA2MfkrWN3T5RoPLnNfVejANrg0w592P97fmiXP6xWt
Hvpt0yBG+nKlbe/8+D+7mx12I3FjIBUH6wM9+Dxqsta90oKihJMPYBKNeUYbdnf6
F8vqJ02aRK6xvwDjmDT9H6zyCKyNXDi9djeio+UCgYEAm1g/2O42eVBzYG/WTJ2H
7jplYemfbIFVpl3Uo18UJKZl/AIzm9tw/+c884naSubwQ8TukI3bDjwWu99R26bo
HQRmLMnP+t7xp84Rn4jXReWlr9sexXHPg/Lj25MdR1t3Ow53qSh4nw/cUPr42N9o
cX4iWLb38v7KEornZfofXEw=
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
					resource.TestCheckResourceAttr(resourceName, "expires_at", "2032-09-21T00:18:24Z"),
					resource.TestCheckResourceAttr(resourceName, "fingerprint", "7aa7421b97c7e8d9d49a473d63af5b89a6e034f9"),
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
