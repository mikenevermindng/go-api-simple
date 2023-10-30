package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mssola/user_agent"
	"goApiByGin/db"
	"goApiByGin/ent/token"
	"goApiByGin/ent/user"
	"goApiByGin/helpers"
	customvalidator "goApiByGin/validator"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
)

func randomSecret() []byte {
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 20
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return []byte(string(buf))
}

type SignUpInput struct {
	Username string `validate:"required,email" json:"username"`
	Password string `validate:"required,password" json:"password"`
}

func SignUp(c *gin.Context) {

	var signUpInput SignUpInput

	if err := c.ShouldBind(&signUpInput); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}
	if err := customvalidator.CustomValidator.Struct(signUpInput); err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpInput.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}

	user, err := db.
		Client().User.
		Create().
		SetPassword(string(hashedPassword)).
		SetUsername(signUpInput.Username).
		Save(context.Background())

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}

	if _, err := db.
		Client().Role.
		Create().
		SetUser(user.ID).
		Save(context.Background()); err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, helpers.Response{Code: http.StatusOK, Message: "success", Data: nil})
	return
}

type SignInInput struct {
	Username string `validate:"required,email" json:"username"`
	Password string `validate:"required,password" json:"password"`
}

type SignInClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func SignIn(c *gin.Context) {
	c.Request.UserAgent()
	var signInInput SignUpInput
	if err := c.ShouldBind(&signInInput); err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}

	account, err := db.Client().User.Query().Where(user.Username(signInInput.Username)).Only(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(signInInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	refreshExpirationTime := time.Now().AddDate(0, 1, 0)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &SignInClaims{
		Username: account.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &SignInClaims{
		Username: account.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	})

	secret := randomSecret()

	tokenString, _ := token.SignedString(secret)
	refreshTokenString, _ := refreshToken.SignedString(secret)

	ua := user_agent.New(c.Request.UserAgent())

	_, insertTokenError := db.Client().
		Token.
		Create().
		SetSecret(string(secret)).
		SetToken(tokenString).
		SetRefreshToken(refreshTokenString).
		SetUser(account.ID).
		SetDevice(ua.Model()).
		SetIP(c.ClientIP()).
		Save(context.Background())

	if insertTokenError != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: insertTokenError.Error(), Data: nil})
		return
	}

	c.JSON(http.StatusOK, helpers.Response{Code: http.StatusOK, Message: "success", Data: gin.H{
		"token":        tokenString,
		"refreshToken": refreshTokenString,
	}})
	return
}

type RefreshInput struct {
	RefreshToken string `validate:"required" json:"refreshToken"`
}

func Refresh(c *gin.Context) {
	var refreshInput RefreshInput
	if err := c.ShouldBind(&refreshInput); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: nil})
		return
	}

	tokenObj, err := db.Client().Token.Query().Where(token.RefreshToken(refreshInput.RefreshToken)).Only(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: "Invalid token", Data: nil})
		return
	}

	token, parseTokenErr := jwt.Parse(refreshInput.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tokenObj.Secret), nil
	})

	if parseTokenErr != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: parseTokenErr.Error(), Data: nil})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: "Invalid token", Data: nil})
		return
	}

	username := claims["username"]
	account, getAccountErr := db.Client().User.Query().Where(user.Username(username.(string))).Only(context.Background())

	if getAccountErr != nil {
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: getAccountErr.Error(), Data: nil})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &SignInClaims{
		Username: account.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	})

	fmt.Println([]byte(tokenObj.Secret))

	if tokenString, signErr := newToken.SignedString([]byte(tokenObj.Secret)); signErr != nil {
		fmt.Println(signErr)
		c.JSON(http.StatusBadRequest, helpers.Response{Code: http.StatusBadRequest, Message: signErr.Error(), Data: nil})
	} else {
		db.Client().Token.UpdateOne(tokenObj).SetToken(tokenString).SetUpdatedAt(time.Now()).Exec(context.Background())
		c.JSON(http.StatusOK, helpers.Response{Code: http.StatusOK, Message: "success", Data: gin.H{
			"token": tokenString,
		}})
	}

	return
}
