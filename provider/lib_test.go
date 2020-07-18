package provider

import (
	"github.com/signalsciences/go-sigsci"
	"testing"
	"time"
)

func TestDiffTemplateDetectionsUpdate(t *testing.T) {
	old := getDefaultDetection()
	new := getDefaultDetection()

	new.Fields[0].Value = "blah"
	add, update, del := diffTemplateDetections([]sigsci.Detection{old}, []sigsci.Detection{new})
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
	add, update, del := diffTemplateDetections([]sigsci.Detection{old}, []sigsci.Detection{new})
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

func getDefaultDetection() sigsci.Detection {
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
		Created:   time.Now(),
		CreatedBy: "lib_test.go",
	}
}
