package oauth

import (
	"errors"
	"strings"
)

type Parser struct {
	TokenStore TokenStore
}

func NewParser(tokenStore TokenStore) *Parser {
	return &Parser{
		TokenStore: tokenStore,
	}
}

func (p *Parser) Parse(accessToken string) (accessTokenClient OauthAccessToken, err error) {
	token := strings.Split(accessToken, " ")

	if !p.validTokenTypeBearer(token[0]) {
		err = errors.New(ErrorTokenTypeMismatch)
		return
	}

	accessTokenClient, err = p.TokenStore.resolveAccessTokenByAccessToken(token[1])
	if err != nil {
		return
	}

	return
}

func (p *Parser) validTokenTypeBearer(tokenType string) bool {
	if tokenType != string(Bearer) {
		return false
	}

	return true
}
