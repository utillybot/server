package discord

import (
	"encoding/json"
	"io"
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

func GetGuilds(token string) ([]PartialGuild, error) {
	response, err := makeRequest("https://discord.com/api/v8/users/@me/guilds", RequestOptions{
		Method: http.MethodGet,
		Body:   nil,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		},
	})
	if err != nil {
		return nil, err
	}

	var guilds []PartialGuild

	err = json.NewDecoder(response.Body).Decode(&guilds)
	if err != nil {
		return nil, err
	}

	return guilds, nil
}

func GetUser(token string) (*User, error) {
	response, err := makeRequest("https://discord.com/api/v8/users/@me", RequestOptions{
		Method: http.MethodGet,
		Body:   nil,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		},
	})
	if err != nil {
		return nil, err
	}

	var user User

	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type RequestOptions struct {
	Method  string
	Body    io.Reader
	Headers map[string]string
}

func makeRequest(url string, options RequestOptions) (*http.Response, error) {
	req, err := http.NewRequest(options.Method, url, options.Body)
	if err != nil {
		return nil, err
	}

	for key, val := range options.Headers {
		req.Header.Set(key, val)
	}

	return http.DefaultClient.Do(req)
}
