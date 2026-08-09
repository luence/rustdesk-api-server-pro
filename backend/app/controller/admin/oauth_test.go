package admin

import (
	"errors"
	"net/url"
	"rustdesk-api-server-pro/internal/errcode"
	"testing"
)

func TestValidateOAuthRedirectURL(t *testing.T) {
	for _, valid := range []string{"", "https://desk.example.com/admin/auth/oauth/github/callback", "http://127.0.0.1/callback"} {
		if err := validateOAuthRedirectURL(valid); err != nil {
			t.Fatalf("expected %q to be valid: %v", valid, err)
		}
	}
	for _, invalid := range []string{"javascript:alert(1)", "/relative/callback", "ftp://desk.example.com/callback"} {
		if err := validateOAuthRedirectURL(invalid); err == nil {
			t.Fatalf("expected %q to be rejected", invalid)
		}
	}
}

func TestAdminLoginCallbackTarget(t *testing.T) {
	target := adminLoginCallbackTarget("/#/user/profile")
	parsed, err := url.Parse(target)
	if err != nil {
		t.Fatal(err)
	}
	fragment, err := url.Parse(parsed.Fragment)
	if err != nil {
		t.Fatal(err)
	}
	if fragment.Path != "/login" || fragment.Query().Get("redirect") != "/#/user/profile" {
		t.Fatalf("unexpected callback target: %q", target)
	}

	rootTarget := adminLoginCallbackTarget("/#/login")
	rootParsed, _ := url.Parse(rootTarget)
	rootFragment, _ := url.Parse(rootParsed.Fragment)
	if rootFragment.Query().Get("redirect") != "/" {
		t.Fatalf("login callback must fall back to root: %q", rootTarget)
	}
}

func TestCleanOAuthValues(t *testing.T) {
	values := cleanOAuthValues([]string{" read:user ", "user:email", "read:user", ""})
	if len(values) != 2 || values[0] != "read:user" || values[1] != "user:email" {
		t.Fatalf("unexpected cleaned values: %#v", values)
	}
}

func TestOAuthCallbackErrorCode(t *testing.T) {
	tests := map[string]string{
		"NoBindableOauthAccount":             errcode.ERR2208.Code,
		"context deadline exceeded":          errcode.ERR2209.Code,
		"StateInvalidOrExpired":              errcode.ERR2210.Code,
		"provider returned invalid response": errcode.ERR2212.Code,
	}
	for message, expected := range tests {
		if actual := oauthCallbackErrorCode(errors.New(message)); actual != expected {
			t.Fatalf("oauthCallbackErrorCode(%q) = %q, want %q", message, actual, expected)
		}
	}
	errcodeTests := map[errcode.Entry]string{
		errcode.ERR2023: errcode.ERR2208.Code,
		errcode.ERR2004: errcode.ERR2210.Code,
		errcode.ERR2030: errcode.ERR2209.Code,
	}
	for entry, expected := range errcodeTests {
		err := errcode.New(entry.Code, entry.Message)
		if actual := oauthCallbackErrorCode(err); actual != expected {
			t.Fatalf("oauthCallbackErrorCode(%q) = %q, want %q", err.Error(), actual, expected)
		}
	}
}

func TestMaskOAuthSecret(t *testing.T) {
	if actual := maskOAuthSecret("github-secret-a0f09583"); actual != "********a0f09583" {
		t.Fatalf("unexpected secret hint: %q", actual)
	}
	if actual := maskOAuthSecret("short"); actual != "********" {
		t.Fatalf("short secrets must not be exposed: %q", actual)
	}
	if actual := maskOAuthSecret(""); actual != "" {
		t.Fatalf("empty secret hint should stay empty: %q", actual)
	}
}
