resource "sigsci_corp_user" "test" {
  email = "test-user@fastly.com"
  role  = "user"

  memberships = [
    {
      site_name = "test-test",
      role      = "user"
    },
  ]
}

