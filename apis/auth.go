package apis

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pc01pc013/task-management/auth"
	"github.com/pc01pc013/task-management/database/entities"
	"github.com/pc01pc013/task-management/enums/authtype"
	modelsAuth "github.com/pc01pc013/task-management/models/auth"
	"github.com/pc01pc013/task-management/utils"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuthApi struct {
	db *gorm.DB
}

func NewAuthApi(dbInstance *gorm.DB) *AuthApi {
	return &AuthApi{db: dbInstance}
}

func (api *AuthApi) Login(c *gin.Context) {
	var req modelsAuth.LoginReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		context := utils.MakeResponseResultFailed("Login Failed. Get data error.")
		c.JSON(http.StatusOK, context)
		return
	}

	var userEntities entities.User
	if result := api.db.Where("Username = ? AND auth_type = ?", req.Username, authtype.Local).First(&userEntities); result.Error != nil {
		log.Printf("Login Error: %q", result.Error)
		context := utils.MakeResponseResultFailed("Login Failed. Username or Password not correct.")
		c.JSON(http.StatusOK, context)
		return
	}

	dk, err := scrypt.Key([]byte(req.Password), []byte(os.Getenv("APPSETTING_PASSWORDSALT")), 32768, 8, 1, 32)
	if err != nil {
		log.Printf("Login Error: %q", err)
		context := utils.MakeResponseResultFailed("Login Failed. Username or Password not correct.")
		c.JSON(http.StatusOK, context)
		return
	}
	if strings.Compare(*userEntities.Password, base64.StdEncoding.EncodeToString(dk)) != 0 {
		log.Printf("Login Error: Password not equal.")
		log.Println(base64.StdEncoding.EncodeToString(dk))
		context := utils.MakeResponseResultFailed("Login Failed. Username or Password not correct.")
		c.JSON(http.StatusOK, context)
		return
	}

	ss, err := auth.NewToken(req.Username, uint(authtype.Local))
	log.Printf("%v %v", ss, err)

	context := utils.MakeResponseResultSuccess(ss)
	c.JSON(http.StatusOK, context)
}

func (api *AuthApi) LoginFromGoogleAuth(c *gin.Context) {
	var req modelsAuth.LoginFromGoogleAuthReq

	if err := c.BindJSON(&req); err != nil {
		log.Printf("BindJSON Error: %q", err)
		context := utils.MakeResponseResultFailed("Login Failed. Get data error.")
		c.JSON(http.StatusOK, context)
		return
	}

	ctx := context.Background()

	config := &oauth2.Config{
		ClientID:     os.Getenv("APPSETTING_GOOGLE_OAUTH2_CLIENTID"),
		ClientSecret: os.Getenv("APPSETTING_GOOGLE_OAUTH2_CLIENTSECRET"),
		RedirectURL:  "postmessage",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	// Exchange auth code for OAuth token.
	token, err := config.Exchange(ctx, req.AuthCode)
	if err != nil {
		log.Printf("Exchange Error: %q", err)
		context := utils.MakeResponseResultFailed("Login Failed. Exchange error.")
		c.JSON(http.StatusOK, context)
		return
	}

	client := config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		log.Printf("client.Get Error: %q", err)
		context := utils.MakeResponseResultFailed("Login Failed. Exchange error.")
		c.JSON(http.StatusOK, context)
		return
	}

	defer resp.Body.Close()

	var clientres modelsAuth.GoogleAPIUserinfoRes
	if err := json.NewDecoder(resp.Body).Decode(&clientres); err != nil {
		log.Printf("json Decoder Error: %q", err)
		context := utils.MakeResponseResultFailed("Login Failed. Exchange error.")
		c.JSON(http.StatusOK, context)
		return
	}

	if !clientres.VerifiedEmail {
		log.Printf("clientres.VerifiedEmail false.")
		context := utils.MakeResponseResultFailed("Login Failed. Not Verified Email.")
		c.JSON(http.StatusOK, context)
		return
	}

	var userinfo_GoogleEntities entities.Userinfo_Google
	if result := api.db.Preload(clause.Associations).Where("UID = ?", clientres.ID).Limit(1).Find(&userinfo_GoogleEntities); result.RowsAffected == 0 {
		userinfo_GoogleEntities = entities.Userinfo_Google{
			UID:           clientres.ID,
			Email:         clientres.Email,
			VerifiedEmail: clientres.VerifiedEmail,
			Name:          clientres.Name,
			GivenName:     &clientres.GivenName,
			FamilyName:    &clientres.FamilyName,
			Picture:       &clientres.Picture,
			Locale:        &clientres.Locale,
			User: entities.User{
				Username: clientres.Email,
				Nickname: &clientres.Name,
				AuthType: uint(authtype.GoogleOAuth),
			},
		}

		if resultCreate := api.db.Create(&userinfo_GoogleEntities); resultCreate.Error != nil {
			log.Printf("Google OAuth login auto Create User Error: %q", err)
			context := utils.MakeResponseResultFailed("Login Failed. Exchange error.")
			c.JSON(http.StatusOK, context)
			return
		}

	}

	ss, err := auth.NewToken(userinfo_GoogleEntities.User.Username, uint(authtype.GoogleOAuth))
	log.Printf("%v %v", ss, err)

	context := utils.MakeResponseResultSuccess(ss)
	c.JSON(http.StatusOK, context)
}
