package provider

import (
	"testing"
	"time"

	"github.com/signalsciences/go-sigsci"
)

func TestDiffTemplateDetectionsUpdate(t *testing.T) {
	old := getDefaultDetection()
	new := getDefaultDetection()

	new.Fields[0].Value = "blah"
	add, update, del := diffTemplateDetections("CMDEXE", []sigsci.Detection{old}, []sigsci.Detection{new})
	if len(add) > 0 || len(del) > 0 {
		t.Fail()
	}
	if len(update) != 1 {
		t.Fail()
	}
	if update[0].Fields[0].Value != "blah" {
		t.Fail()
	}
}

func TestDiffTemplateDetectionsAddDel(t *testing.T) {
	old := getDefaultDetection()
	new := getDefaultDetection()

	new.Name = "CMDEXE"
	add, update, del := diffTemplateDetections("CMDEXE", []sigsci.Detection{old}, []sigsci.Detection{new})
	if len(add) != 1 || len(del) != 1 || len(update) != 0 {
		t.Fail()
		return
	}
	if add[0].Name != "CMDEXE" {
		t.Fail()
	}
	if del[0].Name != "LOGINATTEMPT" {
		t.Fail()
	}
}

func TestDiffTemplateAlertsUpdate(t *testing.T) {
	old := getDefaultAlert()
	new := getDefaultAlert()

	new.LongName = "New Long Name"
	add, update, del := diffTemplateAlerts([]sigsci.Alert{old}, []sigsci.Alert{new})

	if len(add) > 0 || len(del) > 0 {
		t.Fail()
	}
	if len(update) != 1 {
		t.Fail()
	}
	if update[0].LongName != "New Long Name" {
		t.Fail()
	}
}

func TestDiffTemplateAlertsAddDel(t *testing.T) {
	old := getDefaultAlert()
	new := getDefaultAlert()

	new.ID = "98765"
	new.LongName = "2"
	add, update, del := diffTemplateAlerts([]sigsci.Alert{old}, []sigsci.Alert{new})

	if len(add) != 1 || len(del) != 1 || len(update) != 0 {
		t.Fail()
		return
	}
	if add[0].LongName != "2" {
		t.Fail()
	}
	if del[0].LongName != getDefaultAlert().LongName {
		t.Fail()
	}
}

func getDefaultDetection() sigsci.Detection {
	now := time.Now()
	return sigsci.Detection{
		DetectionUpdateBody: sigsci.DetectionUpdateBody{
			ID:      "123",
			Name:    "LOGINATTEMPT",
			Enabled: true,
			Fields: []sigsci.ConfiguredDetectionField{
				{Name: "path", Value: "/auth/*"},
				{Name: "path", Value: "/login"},
			},
		},
		Created:   &now,
		CreatedBy: "lib_test.go",
	}
}

func getDefaultAlert() sigsci.Alert {
	return sigsci.Alert{
		AlertUpdateBody: sigsci.AlertUpdateBody{
			LongName:          "longname",
			Interval:          60,
			Threshold:         10,
			SkipNotifications: false,
			Enabled:           true,
			Action:            "info",
		},
		ID: "654321",
	}
}

// function used for debugging. uncomment and set breakpoint at 'return nil' to inspect terraform state
/*
func testInspect() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pm := testAccProvider.Meta().(providerMetadata)
		m := s.Modules[0].Resources
		_ = m
		_ = pm.Corp == pm.Corp

		return nil
	}
}
*/

func TestValidateConditionField(t *testing.T) {
	cases := []struct {
		in      string
		wantErr bool
	}{
		{in: "scheme"},
		{in: "queryParameter"},
		{in: "value"},
		{in: "responseHeader"},

		{in: "unknownconditionfield", wantErr: true},
	}

	for _, tt := range cases {
		got, err := validateConditionField(tt.in, "key")
		if err != nil {
			t.Fatal(err)
		}

		if got != nil && !tt.wantErr {
			t.Errorf("validateConditionField(%q) returned invalid, want valid", tt.in)
		}

		if got == nil && tt.wantErr {
			t.Errorf("validateConditionField(%q) returned valid, want invalid", tt.in)
		}
	}
}
