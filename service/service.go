package service

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wuhan005/gadget"
	log "unknwon.dev/clog/v2"

	"dev.apicon.cn/sdk/user"
)

var (
	ErrNilGinContext = errors.New("gin context is nil")
	ErrUserNotLogin  = errors.New("user not login")
)

const (
	ENV = "APICON_SERVICE"

	ApiconAuthHeader         = "X-Apicon-Auth"
	ApiconUserIDHeader       = "X-Apicon-User-ID"
	ApiconUserNameHeader     = "X-Apicon-User-Name"
	ApiconUserEmailHeader    = "X-Apicon-User-Email"
	ApiconUserNicknameHeader = "X-Apicon-User-Nickname"
	ApiconUserKeyHeader      = "X-Apicon-Key"
	UserIPHeader             = "X-Real-Ip"
)

type Service struct {
	id   uint
	name string

	route *gin.Engine
}

// New returns a new service instance.
func New(serviceName string, serviceID uint) *Service {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(gadget.MakeSuccessJSON(gin.H{
			// TODO
		}))
	})

	// Not found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(gadget.MakeErrJSON(40400,
			fmt.Sprintf("Route not found. Find the API documentation for %s API at https://apicon.cn/_/%d", serviceName, serviceID),
		))
	})

	return &Service{
		id:   serviceID,
		name: serviceName,

		route: r,
	}
}

// Route returns the Gin engine.
func (s *Service) Route() *gin.Engine {
	return s.route
}

func (s *Service) Run(addr ...string) {
	var address string
	if len(addr) != 0 {
		address = addr[0]
	}

	err := s.route.Run(address)
	if err != nil {
		log.Fatal("Failed to start service: %v", err)
	}
}

// IsLogin returns user login status by checking the Gin request context.
func IsLogin(c *gin.Context) bool {
	if c == nil {
		return false
	}

	isLogin := c.GetHeader(ApiconAuthHeader)
	return isLogin == "ok"
}

// GetUser returns user information from the Gin request context.
func GetUser(c *gin.Context) (*user.User, error) {
	if c == nil {
		return nil, ErrNilGinContext
	}

	if c.GetHeader(ApiconAuthHeader) != "ok" {
		return nil, ErrUserNotLogin
	}

	userIDStr := c.GetHeader(ApiconUserIDHeader)
	userName := c.GetHeader(ApiconUserNameHeader)
	userEmail := c.GetHeader(ApiconUserEmailHeader)
	userNickName := c.GetHeader(ApiconUserNicknameHeader)
	userKey := c.GetHeader(ApiconUserKeyHeader)
	userIP := c.GetHeader(UserIPHeader)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, errors.Wrapf(err, "convert userID from string to int: %q", userIDStr)
	}

	return &user.User{
		ID:       uint(userID),
		Name:     userName,
		Email:    userEmail,
		NickName: userNickName,
		Key:      userKey,
		IP:       userIP,
	}, nil
}

// GetUserIP returns user IP address from the Gin request context.
func GetUserIP(c *gin.Context) (string, error) {
	if c == nil {
		return "", ErrNilGinContext
	}
	return c.GetHeader(UserIPHeader), nil
}
