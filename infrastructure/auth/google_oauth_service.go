package auth

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// client configuration
type GoogleOAuthService struct {
	config *oauth2.Config
}

func NewGoogleOAuthService() *GoogleOAuthService {
	return &GoogleOAuthService{
		config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			//------define and its initial value automatically given since:
			//These URLs (Auth URL and Token URL) are the same for all Google OAuth clients.
			//INTERNALLY CONFIGURED LIKE THIS👇🏾
			//Endpoint: oauth2.Endpoint{
			//	AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			//	TokenURL: "https://oauth2.googleapis.com/token",
			//}
			Endpoint: google.Endpoint,
		},
	}
}

// create google end point start with "https://accounts.google.com/o/oauth2/auth" and other qeries(client id, state)
func (g *GoogleOAuthService) GetAuthURL(state string) string {
	//state is random str given by backen user privent csrs decieve i.e it allow client the redirected response is from google
	//CSRF (Cross-Site Request Forgery)
	return g.config.AuthCodeURL(state)
}


// Get user info from google
func (g *GoogleOAuthService) GetUserInfo(ctx context.Context, code string) (*entities.GoogleUser, error) {
	//"Code" param is code sent to client server for tocken exchange(client secret as password, code as email verification, then return token as access token to client server)
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	//create authenticated client
	client := g.config.Client(ctx, token)
	//get response from "https://www.googleapis.com/oauth2/v2/userinfo" endpoint
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get user info from Google")
	}

	//Get user and bind to user
	var user entities.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}