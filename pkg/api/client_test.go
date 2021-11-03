package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"reflect"
	"testing"

	"github.com/plumming/dx/pkg/auth"

	"github.com/stretchr/testify/assert"

	"github.com/plumming/dx/pkg/httpmock"
)

func eq(t *testing.T, got interface{}, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestGraphQL(t *testing.T) {
	http := &httpmock.Registry{}

	var authConfig auth.Config = &auth.FakeConfig{
		Hosts: map[string]*auth.HostConfig{
			"github.com": {User: "user", Token: "token"},
			"other.com":  {User: "otheruser", Token: "token2"},
		},
	}

	client := NewClient(
		authConfig,
		ReplaceTripper(http),
		AddHeader("Authorization", "token OTOKEN"),
	)

	vars := map[string]interface{}{"name": "Mona"}
	response := struct {
		Viewer struct {
			Login string
		}
	}{}

	http.StubResponse(200, bytes.NewBufferString(`{"data":{"viewer":{"login":"hubot"}}}`))
	err := client.GraphQL("github.com", "QUERY", vars, &response)
	eq(t, err, nil)
	eq(t, response.Viewer.Login, "hubot")

	req := http.Requests[0]
	reqBody, _ := ioutil.ReadAll(io.LimitReader(req.Body, 10000000))
	eq(t, string(reqBody), `{"query":"QUERY","variables":{"name":"Mona"}}`)
	eq(t, req.Header.Get("Authorization"), "token OTOKEN")
}

func TestGraphQLError(t *testing.T) {
	http := &httpmock.Registry{}

	var authConfig auth.Config = &auth.FakeConfig{
		Hosts: map[string]*auth.HostConfig{
			"github.com": {User: "user", Token: "token"},
			"other.com":  {User: "otheruser", Token: "token2"},
		},
	}

	client := NewClient(authConfig, ReplaceTripper(http))

	response := struct{}{}
	http.StubResponse(200, bytes.NewBufferString(`{"errors":[{"message":"OH NO"}]}`))
	err := client.GraphQL("github.com", "", nil, &response)
	if err == nil || err.Error() != "graphql error: 'OH NO'" {
		t.Fatalf("got %q", err.Error())
	}
}

func TestRESTGetDelete(t *testing.T) {
	http := &httpmock.Registry{}

	var authConfig auth.Config = &auth.FakeConfig{
		Hosts: map[string]*auth.HostConfig{
			"github.com": {User: "user", Token: "token"},
			"other.com":  {User: "otheruser", Token: "token2"},
		},
	}

	client := NewClient(authConfig, ReplaceTripper(http))
	http.StubResponse(204, bytes.NewBuffer([]byte{}))

	r := bytes.NewReader([]byte(`{}`))
	err := client.REST("github.com", "DELETE", "applications/CLIENTID/grant", r, nil)
	eq(t, err, nil)
}

const letterBytes = "0123456789abcdef"

func TestFilterToken(t *testing.T) {
	token := generateRandomString()
	masked := Mask(token)

	t.Logf(masked)
	assert.Equal(t, len(token), len(masked))
	assert.NotEqual(t, token, masked)
}

func TestFilterToken_WithPrefix(t *testing.T) {
	token := "token " + generateRandomString()
	masked := Mask(token)

	t.Logf(masked)
	assert.Equal(t, len(token), len(masked))
	assert.NotEqual(t, token, masked)
}

func TestFilterToken_WithSuffix(t *testing.T) {
	token := generateRandomString() + " something at the end"
	masked := Mask(token)

	t.Logf(masked)
	assert.Equal(t, len(token), len(masked))
	assert.NotEqual(t, token, masked)
}

func TestFilterToken_NotToken(t *testing.T) {
	token := "not a token"
	masked := Mask(token)

	assert.Equal(t, len(token), len(masked))
	assert.Equal(t, token, masked)
}

func TestGetV3APIForHost(t *testing.T) {
	assert.Equal(t, "https://api.github.com/", GetV3APIForHost("github.com"))
	assert.Equal(t, "https://other.com/api/v3/", GetV3APIForHost("other.com"))
}

func TestGetGraphQLAPIForHost(t *testing.T) {
	assert.Equal(t, "https://api.github.com/graphql", GetGraphQLAPIForHost("github.com"))
	assert.Equal(t, "https://other.com/api/graphql", GetGraphQLAPIForHost("other.com"))
}

func generateRandomString() string {
	b := make([]byte, 40)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))] // #nosec:G404
	}
	return string(b)
}
