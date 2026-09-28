package selfupdate

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCLIAssetName(t *testing.T) {
	cases := []struct {
		goos, arch, want string
		ok               bool
	}{
		{"linux", "amd64", "aroli-linux-amd64", true},
		{"linux", "arm64", "aroli-linux-arm64", true},
		{"darwin", "arm64", "", false},
		{"linux", "386", "", false},
	}
	for _, test := range cases {
		got, err := cliAssetName(test.goos, test.arch)
		if (err == nil) != test.ok || got != test.want {
			t.Errorf("cliAssetName(%q, %q) = %q, %v", test.goos, test.arch, got, err)
		}
	}
}

func TestChecksumFor(t *testing.T) {
	checksum := strings.Repeat("a", 64)
	contents := checksum + "  aroli-linux-amd64\n" + strings.Repeat("b", 64) + " *aroli-linux-arm64\n"
	got, err := checksumFor("aroli-linux-arm64", contents)
	if err != nil || got != strings.Repeat("b", 64) {
		t.Fatalf("checksumFor returned %q, %v", got, err)
	}
	if _, err := checksumFor("missing", contents); err == nil {
		t.Fatal("expected missing checksum error")
	}
	if _, err := checksumFor("bad", "short  bad\n"); err == nil {
		t.Fatal("expected malformed checksum error")
	}
	if _, err := checksumFor("bad", strings.Repeat("g", 64)+"  bad\n"); err == nil {
		t.Fatal("expected non-hex checksum error")
	}
}

func TestLatestReleaseSelectsNewestPublishedBeta(t *testing.T) {
	client := releaseTestClient(
		githubRelease{TagName: "v0.1.0-beta.2", Prerelease: true},
		githubRelease{TagName: "v0.1.0-beta.10", Prerelease: true},
		githubRelease{TagName: "v9.0.0", Prerelease: false},
		githubRelease{TagName: "v0.2.0-beta.1", Prerelease: false},
	)
	got, err := latestRelease(client)
	if err != nil || got.TagName != "v0.1.0-beta.10" {
		t.Fatalf("latestRelease() = %#v, %v", got, err)
	}
}

func TestLatestReleaseRejectsMissingBeta(t *testing.T) {
	for _, releases := range [][]githubRelease{
		{},
		{{TagName: "v1.2.3", Prerelease: false}},
		{{TagName: "v1.2.3-rc.1", Prerelease: true}},
		{{TagName: "v1.2.3/../../other", Prerelease: true}},
	} {
		if _, err := latestRelease(releaseTestClient(releases...)); err == nil {
			t.Fatalf("latestRelease accepted non-beta releases %#v", releases)
		}
	}
}

func TestLatestReleaseFallsBackToBetaTag(t *testing.T) {
	client := multiResponseClient(
		[]githubRelease{{TagName: "v3.4.6", Prerelease: false}},
		[]struct {
			Name string `json:"name"`
		}{{Name: "v3.4.6"}, {Name: "v0.1.0-beta.1"}},
	)
	got, err := latestRelease(client)
	if err != nil || got.TagName != "v0.1.0-beta.1" {
		t.Fatalf("latestRelease() fallback = %#v, %v", got, err)
	}
}

func releaseTestClient(releases ...githubRelease) *http.Client {
	body, _ := json.Marshal(releases)
	return &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader(body)),
			Request:    request,
		}, nil
	})}
}

func multiResponseClient(releases []githubRelease, tags []struct {
	Name string `json:"name"`
}) *http.Client {
	bodies := [][]byte{}
	releaseBody, _ := json.Marshal(releases)
	tagBody, _ := json.Marshal(tags)
	bodies = append(bodies, releaseBody, tagBody)
	call := 0
	return &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		body := bodies[call]
		call++
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader(body)),
			Request:    request,
		}, nil
	})}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return f(request) }
