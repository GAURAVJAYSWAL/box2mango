package box

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"github.com/siddhartham/box"
	"golang.org/x/oauth2"
)

func getEntpUsers(tok *oauth2.Token) (*box.Users, error) {
	var (
		configSource = box.NewConfigSource(
			&oauth2.Config{
				ClientID:     os.Getenv("CLIENTID"),
				ClientSecret: os.Getenv("CLIENTSECRET"),
				Scopes:       nil,
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://app.box.com/api/oauth2/authorize",
					TokenURL: "https://app.box.com/api/oauth2/token",
				},
				RedirectURL: "http://localhost:8080/handle",
			},
		)
		c = configSource.NewClient(tok)
	)

	_, users, err := c.UserService().GetEnterpriseUsers(0, 1000)

	return users, err

}

func getEntpToken() (*oauth2.Token, string, error) {
	privateKeyData, err := ioutil.ReadFile(os.Getenv("PRIVATEKEY"))
	if err != nil {
		return nil, "Getting PrivateKey", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, "Parsing PrivateKey", err
	}

	u1 := uuid.Must(uuid.NewV4())

	token := jwt.New(jwt.GetSigningMethod("RS256"))

	claims := make(jwt.MapClaims)
	claims["iss"] = os.Getenv("CLIENTID")
	claims["sub"] = os.Getenv("ENTERPRISEID")
	claims["box_sub_type"] = "enterprise"
	claims["aud"] = "https://api.box.com/oauth2/token"
	claims["jti"] = u1.String()
	claims["exp"] = time.Now().Add(time.Second * 60).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	token.Header["kid"] = os.Getenv("PUBLICKEYID")

	tokenString, _ := token.SignedString(privateKey)

	apiURL := "https://api.box.com"
	resource := "/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	data.Add("client_id", os.Getenv("CLIENTID"))
	data.Add("client_secret", os.Getenv("CLIENTSECRET"))
	data.Add("assertion", tokenString)

	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	body, _ := ioutil.ReadAll(resp.Body)

	entpToken := new(oauth2.Token)
	if value, err := jsonparser.GetString(body, "access_token"); err == nil {
		entpToken.AccessToken = value
	}
	// if value, err := jsonparser.GetInt(body, "expires_in"); err == nil {
	// 	entpToken.Expiry = value
	// }

	fmt.Printf("Got AccessToken: %v", entpToken)

	return entpToken, "", nil
}
