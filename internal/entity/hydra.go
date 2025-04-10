package entity

type IDPScope struct {
	Scope  string `json:"scope" binding:"required"`
	Domain string `json:"domain" binding:"required"`
	Group  string `json:"group" binding:"required"`
	Role   string `json:"role"`
}

type OAuth2LoginRequest struct {
	Challenge                    string   `json:"challenge"`
	RequestUrl                   string   `json:"request_url"`
	RequestedAccessTokenAudience []string `json:"requested_access_token_audience,omitempty"`
	RequestedScope               []string `json:"requested_scope,omitempty"`
	SessionId                    *string  `json:"session_id,omitempty"`
	Skip                         bool     `json:"skip"`
	Subject                      string   `json:"subject"`
	AdditionalProperties         map[string]interface{}
}

type OAuth2ConsentRequest struct {
	Challenge string `json:"challenge"`

	Context        interface{} `json:"context,omitempty"`
	LoginChallenge *string     `json:"login_challenge,omitempty"`
	LoginSessionId *string     `json:"login_session_id,omitempty"`

	RequestUrl                   *string  `json:"request_url,omitempty"`
	RequestedAccessTokenAudience []string `json:"requested_access_token_audience,omitempty"`
	RequestedScope               []string `json:"requested_scope,omitempty"`
	Skip                         *bool    `json:"skip,omitempty"`
	Subject                      *string  `json:"subject,omitempty"`
	AdditionalProperties         map[string]interface{}
}

type IntrospectedOAuth2Token struct {
	// Active is a boolean indicator of whether or not the presented token is currently active.  The specifics of a token's \"active\" state will vary depending on the implementation of the authorization server and the information it keeps about its tokens, but a \"true\" value return for the \"active\" property will generally indicate that a given token has been issued by this authorization server, has not been revoked by the resource owner, and is within its given time window of validity (e.g., after its issuance time and before its expiration time).
	Active bool `json:"active"`
	// Audience contains a list of the token's intended audiences.
	Aud []string `json:"aud,omitempty"`
	// ID is aclient identifier for the OAuth 2.0 client that requested this token.
	ClientId *string `json:"client_id,omitempty"`
	// Expires at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token will expire.
	Exp *int64 `json:"exp,omitempty"`
	// Extra is arbitrary data set by the session.
	Ext map[string]interface{} `json:"ext,omitempty"`
	// Issued at is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token was originally issued.
	Iat *int64 `json:"iat,omitempty"`
	// IssuerURL is a string representing the issuer of this token
	Iss *string `json:"iss,omitempty"`
	// NotBefore is an integer timestamp, measured in the number of seconds since January 1 1970 UTC, indicating when this token is not to be used before.
	Nbf *int64 `json:"nbf,omitempty"`
	// ObfuscatedSubject is set when the subject identifier algorithm was set to \"pairwise\" during authorization. It is the `sub` value of the ID Token that was issued.
	ObfuscatedSubject *string `json:"obfuscated_subject,omitempty"`
	// Scope is a JSON string containing a space-separated list of scopes associated with this token.
	Scope *string `json:"scope,omitempty"`
	// Subject of the token, as defined in JWT [RFC7519]. Usually a machine-readable identifier of the resource owner who authorized this token.
	Sub *string `json:"sub,omitempty"`
	// TokenType is the introspected token's type, typically `Bearer`.
	TokenType *string `json:"token_type,omitempty"`
	// TokenUse is the introspected token's use, for example `access_token` or `refresh_token`.
	TokenUse *string `json:"token_use,omitempty"`
	// Username is a human-readable identifier for the resource owner who authorized this token.
	Username             *string `json:"username,omitempty"`
	AdditionalProperties map[string]interface{}
}

type UserInfo struct {
	Sub         string   `json:"sub"`
	Email       string   `json:"email"`
	Family      string   `json:"family_name"`
	GivenName   string   `json:"given_name"`
	Name        string   `json:"name"`
	Company     string   `json:"company"`
	Department  string   `json:"department"`
	Title       string   `json:"title"`
	Nickname    string   `json:"nickname"`
	PhoneNumber string   `json:"phone_number"`
	Groups      []string `json:"groups"`
}

type OAuth2Client struct {
	// OAuth 2.0 Client ID  The ID is immutable. If no ID is provided, a UUID4 will be generated.
	ClientId *string `json:"client_id,omitempty"`
	// OAuth 2.0 Client Name  The human-readable name of the client to be presented to the end-user during authorization.
	ClientName *string `json:"client_name,omitempty"`
}
