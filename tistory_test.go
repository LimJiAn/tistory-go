package tistory

import (
	"fmt"
	"testing"
)

func TestTistory_GetAuthorizationCode(t *testing.T) {
	type fields struct {
		BlogURL            string
		ClientId           string
		ClientSecret       string
		AccessToken        string
		AuthenticationURL  string
		RedirectAuthURL    string
		AuthenticationCode string
	}
	type args struct {
		id       string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid credentials",
			fields: fields{
				BlogURL:      "https://example.com",
				ClientId:     "client_id",
				ClientSecret: "client_secret",
			},
			args: args{
				id:       "username",
				password: "password",
			},
			wantErr: false,
		},
		{
			name: "empty credentials",
			fields: fields{
				BlogURL:      "https://example.com",
				ClientId:     "client_id",
				ClientSecret: "client_secret",
			},
			args: args{
				id:       "",
				password: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tst := &Tistory{
				BlogURL:            tt.fields.BlogURL,
				ClientId:           tt.fields.ClientId,
				ClientSecret:       tt.fields.ClientSecret,
				AccessToken:        tt.fields.AccessToken,
				AuthenticationURL:  tt.fields.AuthenticationURL,
				RedirectAuthURL:    tt.fields.RedirectAuthURL,
				AuthenticationCode: tt.fields.AuthenticationCode,
			}
			got, err := tst.GetAuthorizationCode(tt.args.id, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tistory.GetAuthorizationCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("Tistory.GetAuthorizationCode() = %v, want non-empty string", got)
			}
		})
	}
}

func TestTistory_GetAccessToken(t *testing.T) {
	type fields struct {
		BlogURL            string
		ClientId           string
		ClientSecret       string
		AccessToken        string
		AuthenticationURL  string
		RedirectAuthURL    string
		AuthenticationCode string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "valid credentials",
			fields: fields{
				BlogURL:            "https://example.com",
				ClientId:           "client_id",
				ClientSecret:       "client_secret",
				AuthenticationCode: "authentication_code",
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			fields: fields{
				BlogURL:            "https://example.com",
				ClientId:           "client_id",
				ClientSecret:       "client_secret",
				AuthenticationCode: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tst := &Tistory{
				BlogURL:            tt.fields.BlogURL,
				ClientId:           tt.fields.ClientId,
				ClientSecret:       tt.fields.ClientSecret,
				AccessToken:        tt.fields.AccessToken,
				AuthenticationURL:  tt.fields.AuthenticationURL,
				RedirectAuthURL:    tt.fields.RedirectAuthURL,
				AuthenticationCode: tt.fields.AuthenticationCode,
			}
			got, err := tst.GetAccessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tistory.GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("Tistory.GetAccessToken() = %v, want non-empty string", got)
			}
		})
	}
}

func TestTistory_GetBlogInfo(t *testing.T) {
	type fields struct {
		BlogURL            string
		ClientId           string
		ClientSecret       string
		AccessToken        string
		AuthenticationURL  string
		RedirectAuthURL    string
		AuthenticationCode string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "valid credentials",
			fields: fields{
				BlogURL:      "https://example.com",
				ClientId:     "client_id",
				ClientSecret: "client_secret",
				AccessToken:  "access_token",
				AuthenticationURL: fmt.Sprintf(
					"https://www.tistory.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", "client_id", "https://example.com"),
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			fields: fields{
				BlogURL:            "https://example.com",
				ClientId:           "client_id",
				ClientSecret:       "client_secret",
				AccessToken:        "",
				AuthenticationURL:  "",
				RedirectAuthURL:    "",
				AuthenticationCode: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tst := &Tistory{
				BlogURL:            tt.fields.BlogURL,
				ClientId:           tt.fields.ClientId,
				ClientSecret:       tt.fields.ClientSecret,
				AccessToken:        tt.fields.AccessToken,
				AuthenticationURL:  tt.fields.AuthenticationURL,
				RedirectAuthURL:    tt.fields.RedirectAuthURL,
				AuthenticationCode: tt.fields.AuthenticationCode,
			}
			got, err := tst.GetBlogInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tistory.GetBlogInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("Tistory.GetBlogInfo() = %v, want non-empty string", got)
			}
		})
	}
}
