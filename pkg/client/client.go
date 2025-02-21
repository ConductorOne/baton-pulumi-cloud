package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
)

// Client represents a Pulumi API client
type Client struct {
	baseHttpClient *uhttp.BaseHttpClient
	baseURL        *url.URL
	token          string
}

// NewClient creates a new Pulumi API client
func NewClient(token string) (*Client, error) {
	baseURL, err := url.Parse("https://api.pulumi.com")
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	httpClient, err := uhttp.NewClient(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	wrapper, err := uhttp.NewBaseHttpClientWithContext(context.Background(), httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client wrapper: %w", err)
	}

	return &Client{
		baseHttpClient: wrapper,
		baseURL:        baseURL,
		token:          token,
	}, nil
}

// User represents a Pulumi user/member
type UserInfo struct {
	Name        string `json:"name"`
	GithubLogin string `json:"githubLogin"`
	AvatarUrl   string `json:"avatarUrl"`
}

type User struct {
	Role          string   `json:"role"`
	User          UserInfo `json:"user"`
	Created       string   `json:"created"`
	KnownToPulumi bool     `json:"knownToPulumi"`
	VirtualAdmin  bool     `json:"virtualAdmin"`
}

// Team represents a Pulumi team
type Team struct {
	Kind        string     `json:"kind"`
	Name        string     `json:"name"`
	DisplayName string     `json:"displayName"`
	Description string     `json:"description"`
	Members     []UserInfo `json:"members"`
	UserRole    string     `json:"userRole"`
}

// ListUsersResponse represents the paginated response from listing users
type ListUsersResponse struct {
	Members           []User `json:"members"`
	ContinuationToken string `json:"continuationToken,omitempty"`
}

// ListTeamsResponse represents the response from listing teams
type ListTeamsResponse struct {
	Teams []Team `json:"teams"`
}

// requestOptions returns the common request options for Pulumi API requests
func (c *Client) requestOptions(body interface{}) []uhttp.RequestOption {
	options := []uhttp.RequestOption{
		uhttp.WithHeader("Authorization", fmt.Sprintf("token %s", c.token)),
		uhttp.WithHeader("Accept", "application/vnd.pulumi+8"),
		uhttp.WithHeader("Accept", "application/json"),
		uhttp.WithHeader("Content-Type", "application/json"),
	}

	if body != nil {
		options = append(options, uhttp.WithJSONBody(body))
	}

	return options
}

// buildURL creates a full URL for a given path and query parameters
func (c *Client) buildURL(path string, queryParams url.Values) (*url.URL, error) {
	reqURL, err := url.Parse(fmt.Sprintf("/api/%s", path))
	if err != nil {
		return nil, fmt.Errorf("failed to parse request URL: %w", err)
	}

	reqURL = c.baseURL.ResolveReference(reqURL)
	if queryParams != nil {
		reqURL.RawQuery = queryParams.Encode()
	}

	return reqURL, nil
}

// ListUsers returns a list of all users in the organization
func (c *Client) ListUsers(ctx context.Context, orgName string, continuationToken string) (*ListUsersResponse, error) {
	queryParams := url.Values{}
	queryParams.Set("type", "backend")
	if continuationToken != "" {
		queryParams.Set("continuationToken", continuationToken)
	}

	reqURL, err := c.buildURL(fmt.Sprintf("orgs/%s/members", orgName), queryParams)
	if err != nil {
		return nil, err
	}

	req, err := c.baseHttpClient.NewRequest(ctx, "GET", reqURL, c.requestOptions(nil)...)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response ListUsersResponse
	resp, err := c.baseHttpClient.Do(req, uhttp.WithJSONResponse(&response))
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer resp.Body.Close()

	return &response, nil
}

// ListTeams returns a list of all teams in the organization
func (c *Client) ListTeams(ctx context.Context, orgName string) ([]Team, error) {
	reqURL, err := c.buildURL(fmt.Sprintf("orgs/%s/teams", orgName), nil)
	if err != nil {
		return nil, err
	}

	req, err := c.baseHttpClient.NewRequest(ctx, "GET", reqURL, c.requestOptions(nil)...)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var response ListTeamsResponse
	resp, err := c.baseHttpClient.Do(req, uhttp.WithJSONResponse(&response))
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}
	defer resp.Body.Close()

	return response.Teams, nil
}

// GetTeam returns details about a specific team including its members
func (c *Client) GetTeam(ctx context.Context, orgName, teamName string) (*Team, error) {
	reqURL, err := c.buildURL(fmt.Sprintf("orgs/%s/teams/%s", orgName, teamName), nil)
	if err != nil {
		return nil, err
	}

	req, err := c.baseHttpClient.NewRequest(ctx, "GET", reqURL, c.requestOptions(nil)...)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var team Team
	resp, err := c.baseHttpClient.Do(req, uhttp.WithJSONResponse(&team))
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}
	defer resp.Body.Close()

	return &team, nil
}

// RemoveUser removes a user from the organization
func (c *Client) RemoveUser(ctx context.Context, orgName, username string) error {
	reqURL, err := c.buildURL(fmt.Sprintf("orgs/%s/members/%s", orgName, username), nil)
	if err != nil {
		return err
	}

	req, err := c.baseHttpClient.NewRequest(ctx, "DELETE", reqURL, c.requestOptions(nil)...)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.baseHttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to remove user: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// UpdateUserRole changes a user's role in the organization
func (c *Client) UpdateUserRole(ctx context.Context, orgName, username, role string) error {
	reqURL, err := c.buildURL(fmt.Sprintf("orgs/%s/members/%s", orgName, username), nil)
	if err != nil {
		return err
	}

	body := map[string]string{
		"role": role,
	}

	req, err := c.baseHttpClient.NewRequest(ctx, "PATCH", reqURL, c.requestOptions(body)...)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.baseHttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// UpdateTeamMembership modifies a user's membership in a team
func (c *Client) UpdateTeamMembership(ctx context.Context, orgName, teamName, username, action string) error {
	reqURL, err := c.buildURL(fmt.Sprintf("orgs/%s/teams/%s", orgName, teamName), nil)
	if err != nil {
		return err
	}

	body := map[string]string{
		"memberAction": action,
		"member":       username,
	}

	req, err := c.baseHttpClient.NewRequest(ctx, "PATCH", reqURL, c.requestOptions(body)...)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.baseHttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update team membership: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
