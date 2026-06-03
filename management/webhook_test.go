package management

import (
	"encoding/json"
	"testing"
)

func TestWebHookInput_NotifiersSerializedToJSON(t *testing.T) {
	input := WebHookRequest{
		WebHook: WebHookInput{
			Name:      "test-hook",
			Notifiers: []string{"ops@example.com", "dev@example.com"},
		},
	}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	webhook, ok := result["webhook"].(map[string]interface{})
	if !ok {
		t.Fatal("expected webhook key in JSON")
	}
	notifiers, ok := webhook["notifiers"].([]interface{})
	if !ok {
		t.Fatalf("expected notifiers to be an array, got: %v", webhook["notifiers"])
	}
	if len(notifiers) != 2 {
		t.Errorf("expected 2 notifiers, got %d", len(notifiers))
	}
	if notifiers[0].(string) != "ops@example.com" {
		t.Errorf("expected ops@example.com, got %s", notifiers[0])
	}
	if notifiers[1].(string) != "dev@example.com" {
		t.Errorf("expected dev@example.com, got %s", notifiers[1])
	}
}

func TestWebHook_NotifiersDeserializedFromJSON(t *testing.T) {
	raw := `{"webhook":{"name":"test","notifiers":["ops@example.com"]}}`
	var result WebHookResponse
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(result.WebHook.Notifiers) != 1 {
		t.Errorf("expected 1 notifier, got %d", len(result.WebHook.Notifiers))
	}
	if result.WebHook.Notifiers[0] != "ops@example.com" {
		t.Errorf("expected ops@example.com, got %s", result.WebHook.Notifiers[0])
	}
}
