package aisera

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/textproto"
	"strings"

	"golang.org/x/exp/slog"
)

const (
	aiseraAdminCookie = "admin:token"
	sessionHeader     = "session-header"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Permissions map[string]bool

type LoginResponse struct {
	ID             string      `json:"id"`
	Username       string      `json:"username"`
	TenantID       string      `json:"tenantId"`
	Skus           []string    `json:"skus"`
	Roles          []any       `json:"roles"`
	SessionToken   string      `json:"sessionToken"`
	PermissionsMap Permissions `json:"permissionsMap"`
}

func (l LoginRequest) JSON() []byte {
	val, _ := json.Marshal(l)
	return val
}

func Login(ctx context.Context, loginRequest LoginRequest) (Offerings, error) {
	logger.Debug("logging in", slog.String("username", loginRequest.Username))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, MustURL("aisera/login").String(), ToJSONReader(loginRequest))
	if err != nil {
		return nil, fmt.Errorf("error creatig login request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	loginResponse := LoginResponse{}
	resp, err := Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making the request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to login: invalid status code: %d", resp.StatusCode)
	}
	err = JSONReaderToVal(resp.Body, &loginResponse)
	if err != nil {
		return nil, fmt.Errorf("error reading login response in: %w", err)
	}
	authKey := getCookie(resp.Header, aiseraAdminCookie)
	if authKey == "" {
		return nil, errors.New("could not get an authentication key")
	}
	return offering{
		auhtKey: authKey,
	}, nil
}

func getCookie(headerValues http.Header, name string) string {
	for _, line := range headerValues["Set-Cookie"] {
		parts := strings.Split(textproto.TrimString(line), ";")
		if len(parts) == 1 && parts[0] == "" {
			continue
		}
		parts[0] = textproto.TrimString(parts[0])
		cookieName, value, ok := strings.Cut(parts[0], "=")
		if !ok {
			continue
		}
		if name == cookieName {
			return value
		}
	}
	return ""
}
