package service

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/tabloy/keygate/internal/license"
	"github.com/tabloy/keygate/internal/model"
)

func TestSignTokenOfflineExpiryPolicy(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	svc := &LicenseService{signingKey: priv}

	perpetual := &model.License{
		ID: "lic-perpetual", ProductID: "prod", PlanID: "perpetual",
		Status: model.StatusActive, Plan: &model.Plan{LicenseType: "perpetual"},
	}
	raw, err := svc.signToken(perpetual, "device-1")
	if err != nil {
		t.Fatal(err)
	}
	verified, err := license.Verify(raw, pub)
	if err != nil {
		t.Fatal(err)
	}
	if verified.ExpiresAt != 0 {
		t.Fatalf("perpetual token expires_at = %d, want 0", verified.ExpiresAt)
	}

	until := time.Now().Add(30 * 24 * time.Hour).Truncate(time.Second)
	prepaid := &model.License{
		ID: "lic-month", ProductID: "prod", PlanID: "month", Status: model.StatusActive,
		ValidUntil: &until,
		Plan:       &model.Plan{LicenseType: "subscription", DurationDays: 30},
	}
	raw, err = svc.signToken(prepaid, "device-1")
	if err != nil {
		t.Fatal(err)
	}
	verified, err = license.Verify(raw, pub)
	if err != nil {
		t.Fatal(err)
	}
	if verified.ExpiresAt != until.Unix() {
		t.Fatalf("prepaid token expires_at = %d, want %d", verified.ExpiresAt, until.Unix())
	}
}
