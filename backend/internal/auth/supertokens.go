package auth

import (
	"os"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// InitSuperTokens initializes SuperTokens with OAuth providers and email/password
func InitSuperTokens() error {
	apiBasePath := "/auth"
	websiteBasePath := "/auth"

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: getEnv("SUPERTOKENS_CONNECTION_URI", "http://localhost:3567"),
			APIKey:        getEnv("SUPERTOKENS_API_KEY", ""),
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "InsightIQ",
			APIDomain:       getEnv("NEXT_PUBLIC_API_DOMAIN", "http://localhost:8080"),
			WebsiteDomain:   getEnv("NEXT_PUBLIC_WEBSITE_DOMAIN", "http://localhost:3000"),
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			// Third Party (OAuth) Recipe - GitHub and Google
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						// GitHub OAuth
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "github",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
										ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
									},
								},
							},
						},
						// Google OAuth
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
										ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
									},
								},
							},
						},
					},
				},
			}),
			// Email/Password Recipe - Traditional login
			emailpassword.Init(&epmodels.TypeInput{}),
			// Session Recipe
			session.Init(nil),
		},
	})

	return err
}

// getEnv retrieves environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
