package provider

import "testing"

// A remote magpie's gateway keys: a key the remote holds to no limit says
// nothing, so its card is one "no limit" bar and is not shown. A key with a
// limit keeps its card, and a key whose limit could not be read keeps its
// card too — an error is not "no limit".
func TestKeepKeyCards(t *testing.T) {
	unlimited := SubscriptionQuota{
		Provider:   "office/~keys",
		User:       "unlimited-key",
		GatewayKey: &GatewayKeyLimit{Limited: false},
		Windows:    []QuotaWindow{{Name: "Limit", Unlimited: true}},
	}
	limited := SubscriptionQuota{
		Provider:   "office/~keys",
		User:       "limited-key",
		GatewayKey: &GatewayKeyLimit{Limited: true, Period: "day"},
		Windows:    []QuotaWindow{{Name: "Tokens today", Used: 10}},
	}
	unreadable := SubscriptionQuota{
		Provider:   "office/~keys",
		User:       "error-key",
		GatewayKey: &GatewayKeyLimit{Limited: false},
		Windows:    []QuotaWindow{{Name: "Limit", Unlimited: true}},
		Error:      "remote magpie is too old to tell this key's limit",
	}
	// a key the remote calls limited whose cap it cannot read also renders
	// as the identical "no limit" card, so it goes too
	calledLimitedNoCap := SubscriptionQuota{
		Provider:   "office/~keys",
		User:       "no-cap-key",
		GatewayKey: &GatewayKeyLimit{Limited: false},
		Windows:    []QuotaWindow{{Name: "Limit", Unlimited: true}},
	}

	got := keepKeyCards([]SubscriptionQuota{unlimited, limited, unreadable, calledLimitedNoCap})
	if len(got) != 2 {
		t.Fatalf("kept %d cards, want 2 (limited and error): %+v", len(got), got)
	}
	if got[0].User != "limited-key" {
		t.Errorf("first kept card is %q, want the limited key", got[0].User)
	}
	if got[1].User != "error-key" {
		t.Errorf("second kept card is %q, want the error key", got[1].User)
	}

	if got := keepKeyCards(nil); len(got) != 0 {
		t.Errorf("no keys gave %d cards, want none", len(got))
	}
	if got := keepKeyCards([]SubscriptionQuota{unlimited}); len(got) != 0 {
		t.Errorf("an unlimited key gave %d cards, want none", len(got))
	}
}
