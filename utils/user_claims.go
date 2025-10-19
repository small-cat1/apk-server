package utils

import (
	"ApkAdmin/global"
	systemReq "ApkAdmin/model/common/request"
	"ApkAdmin/model/project"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetClientToken(c *gin.Context) string {
	token := c.Request.Header.Get("x-token")
	return token
}

func GetClientClaims(c *gin.Context) (*systemReq.CustomClaims, error) {
	token := GetToken(c)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GVA_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetClientUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetClientUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseClaims.ID
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.BaseClaims.ID
	}
}

// GetClientUserUuid  从Gin的Context中获取从jwt解析出来的用户UUID
func GetClientUserUuid(c *gin.Context) uuid.UUID {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.UUID
	}
}

// GetClientUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetClientUserInfo(c *gin.Context) *systemReq.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse
	}
}

// GetClientUserName  从Gin的Context中获取从jwt解析出来的用户名
func GetClientUserName(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.Username
		}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.Username
	}
}

func ClientLoginToken(user project.UserLogin) (token string, claims systemReq.CustomClaims, err error) {
	j := NewJWT()
	claims = j.CreateClaims(systemReq.BaseClaims{
		UUID:     user.GetUUID(),
		ID:       user.GetUserId(),
		Username: user.GetUsername(),
	})
	token, err = j.CreateToken(claims)
	return
}
