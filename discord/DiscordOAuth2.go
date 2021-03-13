package discord

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type OAuth2 struct {
	ClientId     string
	ClientSecret string
	Scope        []string
}

type AuthURLOptions struct {
	RedirectUri string
}

type TokenRequestOptions struct {
	Code        string
	RedirectUri string
}

func (d OAuth2) GenerateAuthURL(options AuthURLOptions) string {
	requestUrl, _ := url.Parse("https://discord.com/api/oauth2/authorize")

	query := requestUrl.Query()
	query.Set("response_type", "code")
	query.Set("client_id", d.ClientId)
	query.Set("scope", strings.Join(d.Scope, " "))
	query.Set("redirect_uri", options.RedirectUri)
	requestUrl.RawQuery = query.Encode()

	return requestUrl.String()
}

func (d OAuth2) TokenRequest(options TokenRequestOptions) (*TokenRequestResult, error) {
	body := url.Values{}
	body.Set("client_id", d.ClientId)
	body.Set("client_secret", d.ClientSecret)
	body.Set("grant_type", "authorization_code")
	body.Set("code", options.Code)
	body.Set("redirect_uri", options.RedirectUri)
	body.Set("scope", strings.Join(d.Scope, " "))

	response, err := http.PostForm("https://discord.com/api/oauth2/token", body)
	if err != nil {
		return nil, err
	}

	var result TokenRequestResult
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, err
}
