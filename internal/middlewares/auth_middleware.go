package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
)

func AuthMiddleware(auth authenticator.Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// checks if header contains token
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			_ = c.Error(apperror.NewAppError(http.StatusUnauthorized, "no authorization header"))
			// needs c.Abort() here to cancel c.Next(),
			// even you return here but you have c.Next(),
			// the chain will continue, c.Abort() will cancel that
			c.Abort()
			return
		}

		// checks if token's format is right
		arr := strings.Split(bearerToken, " ")
		if len(arr) != 2 || arr[0] != "Bearer" {
			_ = c.Error(
				apperror.NewAppError(
					http.StatusUnauthorized,
					"authorization header invalid format",
				),
			)
			c.Abort()
			return
		}

		// verify token
		token, err := auth.VerifyToken(arr[1])
		if err != nil {
			_ = c.Error(
				err,
			)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			_ = c.Error(
				apperror.NewAppError(http.StatusUnauthorized, "could not cast to MapClaims"),
			)
			c.Abort()
			return
		}

		// does not need to check "exp" manual,
		// jwt library has already done that,
		// just need to have "exp" claim and
		// the library will do it for you
		//
		// exp, err := claims.GetExpirationTime()
		// if err != nil {
		// 	_ = c.Error(
		// 		apperror.NewAppError(http.StatusUnauthorized, "err while getting exp from token"),
		// 	)
		// 	c.Abort()
		// }
		//
		// if !exp.After(time.Now()) {
		// 	_ = c.Error(
		// 		apperror.NewAppError(http.StatusUnauthorized, "exp reached"),
		// 	)
		// 	c.Abort()
		// }
		//

		userID, ok := claims[authenticator.SubClaim]
		if !ok {
			_ = c.Error(
				apperror.NewAppError(http.StatusUnauthorized, "could not find sub claim"),
			)
			c.Abort()
			return
		}
		email, ok := claims[authenticator.EmailClaim]
		if !ok {
			_ = c.Error(
				apperror.NewAppError(http.StatusUnauthorized, "could not find email claim"),
			)
			c.Abort()
			return
		}

		c.Set(authenticator.SubClaim, userID)
		c.Set(authenticator.EmailClaim, email)

		c.Next()
	}
}
