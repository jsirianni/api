package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jsirianni/server/model"
)

// TODO(jsirianni): Handle broken backends, return 500 level errors.

var (
	errAccountNotActive = errors.New("subscription is not active")
)

// AccountRequest represents the request payload
// expected from client requests.
type AccountRequest struct {
	Key string `json:"key"`
}

func healthHandler(c *gin.Context) {
	c.Writer.WriteHeader(200)
}

// checkSubscriptionHandler returns status code 200 if the account id
// and account key combination is a valid subscription.
func (s *Server) checkSubscriptionHandler(c *gin.Context) {
	account, ok := s.isActiveAccount(c)
	if !ok {
		return
	}

	s.logger.Sugar().Debugf("account %s is active", account.ID)

	c.Writer.WriteHeader(200)
}

func (s *Server) registerDeviceHandler(c *gin.Context) {
	c.Writer.WriteHeader(200)
}

func (s *Server) accountHandler(c *gin.Context) {
	c.Writer.WriteHeader(200)
}

func (s *Server) devicesHandler(c *gin.Context) {
	c.Writer.WriteHeader(200)
}

func (s *Server) deviceHandler(c *gin.Context) {
	c.Writer.WriteHeader(200)
}

// Returns the account and true if the account is valid and active.
// Callers should return without writing status codes or response bodies
// when false.
func (s *Server) isActiveAccount(c *gin.Context) (*model.Account, bool) {
	reqBody := AccountRequest{}
	if err := c.BindJSON(&reqBody); err != nil {
		s.logger.Sugar().Debugf("failed to parse request body as json: %v", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return nil, false
	}

	accountID, ok := c.Params.Get("account")
	if !ok || accountID == "" {
		s.logger.Debug("missing account parameter")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return nil, false
	}

	if reqBody.Key == "" {
		s.logger.Debug("missing account key in request body")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return nil, false
	}

	account, err := s.store.Account(accountID)
	if err != nil {
		s.logger.Sugar().Debugf("failed to lookup account %s: %v", accountID, err)
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return nil, false
	}

	if account.Key != reqBody.Key {
		s.logger.Sugar().Debugf("invalid key for account %s", accountID)
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return nil, false
	}

	// TODO(jsirianni): IF the account and key exist but do not have
	// an active subscription, return an error to the user explaining that
	// no valid subscription was found.
	if !account.Active {
		s.logger.Sugar().Debugf("account %s is not active: %v", accountID, errAccountNotActive)
		c.AbortWithError(http.StatusPaymentRequired, errAccountNotActive)
		return nil, false
	}

	return &account, true

}
