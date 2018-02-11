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

func (bs *BoxService) GetFolderItems(folderId string, user box.User) (*box.ItemCollection, error) {
	if user.Role == "user" {
		_, items, err := bs.Client().FolderService().GetFolderItems(folderId, &box.UrlParams{
			Limit:    1000,
			Offset:   0,
			AsUserId: user.ID,
		})
		return items, err
	} else {
		_, items, err := bs.UserClient(user.ID).FolderService().GetFolderItems(folderId, &box.UrlParams{
			Limit:  1000,
			Offset: 0,
		})
		return items, err
	}
}

func (bs *BoxService) GetEntpUsers() (*box.Users, error) {
	_, users, err := bs.Client().UserService().GetEnterpriseUsers(&box.UrlParams{
		Limit:  1000,
		Offset: 0,
	})
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
	EntpToken  *oauth2.Token
	UserTokens []UserToken
}

type UserToken struct {
	UserId    string
	UserToken *oauth2.Token
}

func (bs *BoxService) Client() *box.Client {
	if bs.EntpToken == nil || bs.EntpToken.Expiry.Before(time.Now()) || bs.EntpToken.Expiry.Equal(time.Now()) {
		bs.GetEntpToken()
	}
	c := configSource.NewClient(bs.EntpToken)
	return c
}

func (bs *BoxService) UserClient(userId string) *box.Client {
	var userToken *oauth2.Token
	tokenIndex := -1
	tokenFound := false
	for _, ut := range bs.UserTokens {
		tokenIndex++
		if ut.UserId == userId {
			tokenFound = true
			userToken = ut.UserToken
			break
		}
	}

	if tokenFound == false || userToken == nil || userToken.Expiry.Before(time.Now()) || userToken.Expiry.Equal(time.Now()) {
		privateKeyData, _ := ioutil.ReadFile(os.Getenv("PRIVATEKEY"))
		privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)

		u1 := uuid.Must(uuid.NewV4())

		token := jwt.New(jwt.GetSigningMethod("RS256"))

		claims := make(jwt.MapClaims)
		claims["iss"] = os.Getenv("CLIENTID")
		claims["sub"] = userId
		claims["box_sub_type"] = "user"
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

		userToken = new(oauth2.Token)
		accessToken, _ := jsonparser.GetString(body, "access_token")
		expiresIn, _ := jsonparser.GetInt(body, "expires_in")
		tokenType, _ := jsonparser.GetString(body, "token_type")
		if accessToken != "" {
			userToken.AccessToken = accessToken
			userToken.Expiry = time.Now().Add((time.Second * time.Duration(expiresIn)))
			userToken.TokenType = tokenType
			fmt.Printf("\nGot User AccessToken: %v and ExpireIn: %v", accessToken, expiresIn)
		} else {
			fmt.Printf("\nGot Response: %v", string(body))
		}
	}

	if tokenFound == true {
		bs.UserTokens[tokenIndex] = UserToken{
			UserId:    userId,
			UserToken: userToken,
		}
	} else {
		bs.UserTokens = append(bs.UserTokens, UserToken{
			UserId:    userId,
			UserToken: userToken,
		})
	}

	c := configSource.NewClient(userToken)
	return c
}

func (bs *BoxService) GetEntpToken() (string, error) {
	privateKeyData, _ := ioutil.ReadFile(os.Getenv("PRIVATEKEY"))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)

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
	accessToken, _ := jsonparser.GetString(body, "access_token")
	expiresIn, _ := jsonparser.GetInt(body, "expires_in")
	tokenType, _ := jsonparser.GetString(body, "token_type")
	if accessToken != "" {
		bs.EntpToken.AccessToken = accessToken
		bs.EntpToken.Expiry = time.Now().Add((time.Second * time.Duration(expiresIn)))
		bs.EntpToken.TokenType = tokenType
		fmt.Printf("\nGot AccessToken: %v and ExpireIn: %v", accessToken, expiresIn)
	} else {
		fmt.Printf("\nGot Response: %v", string(body))
	}

	return "", nil
}
