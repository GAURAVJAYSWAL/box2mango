package box2mango

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/CrowdSurge/banner"
	"github.com/buger/jsonparser"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	banner.Print("  box2mango  ")
	color.Red("\n- by Siddhartha Mukherjee <mukherjee.siddhartha@gmail.com>")
	color.Yellow("\nUsage: box2mango\n")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token, errStr, err := getEntpToken()
	if err != nil {
		fmt.Printf("Error getting enterprise token: %v : %v", errStr, err)
	}

	fmt.Printf("Token: %v", token)
}

func getEntpToken() (string, string, error) {
	privateKeyData, err := ioutil.ReadFile("private_key")
	if err != nil {
		return "", "Getting PrivateKey", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", "Parsing PrivateKey", err
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
	data.Set("grant_type", os.Getenv("JWTGRANTTYPE"))
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

	accessToken := ""
	if value, err := jsonparser.GetString(body, "access_token"); err == nil {
		accessToken = value
	}

	fmt.Printf("Got token %v", accessToken)

	return accessToken, "", nil
}
