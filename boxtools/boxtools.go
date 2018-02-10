package boxtools

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

func (bs *BoxService) GetEntpUsers() (*box.Users, error) {
	//Todo: make this call native to this codebase
	_, users, err := bs.Client().UserService().GetEnterpriseUsers(0, 1000)
	return users, err
}

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
)

type BoxService struct {
	EntpToken *oauth2.Token
}

func (bs *BoxService) Client() *box.Client {
	c := configSource.NewClient(bs.EntpToken)
	//Todo: remove this chk with proper one
	if _, _, err := c.UserService().Me(); err != nil {
		bs.GetEntpToken()
		c = configSource.NewClient(bs.EntpToken)
	}
	return c
}

func (bs *BoxService) GetEntpToken() (string, error) {
	privateKeyData, err := ioutil.ReadFile(os.Getenv("PRIVATEKEY"))
	if err != nil {
		return "Getting PrivateKey", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "Parsing PrivateKey", err
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

	bs.EntpToken = new(oauth2.Token)
	if value, err := jsonparser.GetString(body, "access_token"); err == nil {
		bs.EntpToken.AccessToken = value
	}

	fmt.Printf("Got AccessToken: %v", bs.EntpToken)

	return "", nil
}
