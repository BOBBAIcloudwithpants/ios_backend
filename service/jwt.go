package service



import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/bobbaicloudwithpants/ios_backend/models"
	"github.com/spf13/viper"
	"time"
)

var secret = []byte(viper.GetString("secret"))

type Claims struct {
	UserId   int	`json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}
func GenerateToken(user_id int, username string, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		user_id,
		username,
		password,
		jwt.StandardClaims {
			ExpiresAt : expireTime.Unix(),
			Issuer : "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(secret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

// 从上下文中获取当前用户信息
func GetUserFromContext(c *gin.Context) models.User {
	claims, _ := c.MustGet("Claims").(*Claims)
	return models.User{UserId: claims.UserId, Username: claims.Username, Password: claims.Password}
}



