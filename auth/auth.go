package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// token黑名单 accessToken : ExpiresAt
	tokenBlacklist = make(map[string]int64)
	// 用户更新名单 ID : IssuedAt
	userUpdatelist = make(map[int64]int64)
	// AccessTokenKey AccessToken
	AccessTokenKey = "accesstoken"
	// AuthKey AuthKey
	AuthKey = "auth"

	cleanDay = time.Now().Day()
)

//MiddleWare MiddleWare
func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader(AccessTokenKey)
		if expiresAt, ok := tokenBlacklist[accessToken]; ok {
			if time.Now().Unix() >= expiresAt {
				delete(tokenBlacklist, accessToken)
			}
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if time.Now().Day() > cleanDay {
			cleanBlackList()
			cleanDay = time.Now().Day()
		}
		auth, err := parse(accessToken)
		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		// 签发时间
		iss := auth.IssuedAt
		if update, ok := userUpdatelist[auth.ID]; ok && update > iss {
			tokenBlacklist[accessToken] = auth.ExpiresAt
			delete(userUpdatelist, auth.ID)
		}
		c.Set(AuthKey, auth)
		c.Next()
	}
}

func cleanBlackList() {
	now := time.Now().Unix()
	var keyset []string
	for key, expiresAt := range tokenBlacklist {
		if now >= expiresAt {
			keyset = append(keyset, key)
		}
	}
	for _, key := range keyset {
		delete(tokenBlacklist, key)
	}
}
