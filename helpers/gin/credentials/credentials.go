package credentials

import (
	"github.com/gin-gonic/gin"
)

type Credentials struct {
	MerchantID  int
	UserID      int
	ClientID    int
	IP          string
	Fingerprint string
	UserAgent   string
	SessionID   string
	AccessToken string
}

func NewCredentials(merchantID, userID, clientID int, ip, fingerprint, userAgent, sessionID, accessToken string) Credentials {
	return Credentials{
		MerchantID:  merchantID,
		UserID:      userID,
		ClientID:    clientID,
		IP:          ip,
		Fingerprint: fingerprint,
		UserAgent:   userAgent,
		SessionID:   sessionID,
		AccessToken: accessToken,
	}
}

func GetCredentialsFromCtx(ctx *gin.Context) Credentials {
	cred, ok := ctx.MustGet("credentials").(Credentials)
	if !ok {
		return Credentials{}
	}

	return cred
}
