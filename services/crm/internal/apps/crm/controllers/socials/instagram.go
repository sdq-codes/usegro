package crmsocials

import (
	"github.com/gofiber/fiber/v2"
)

type CRMSocialsInstagramController struct {
}

func NewCRMSocialsInstagramController() *CRMSocialsInstagramController {
	return &CRMSocialsInstagramController{}
}

// FacebookLogin godoc
// @Summary      Initiate Instagram / Facebook OAuth login
// @Description  Redirects the user to the Facebook OAuth consent screen to connect their Instagram Business account
// @Tags         CRM Socials
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Success      307  "Redirect to Facebook OAuth"
// @Failure      401  {object}  map[string]string
// @Router       /api/v1/crm/socials/instagram [get]
func (cs *CRMSocialsInstagramController) FacebookLogin(c *fiber.Ctx) error {
	//url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=instagram_business_basic",
	//	config.GetConfig().Facebook.OauthUrl, config.GetConfig().Facebook.AppId, config.GetConfig().Facebook.RedirectUrl)
	//return c.Redirect(url, fiber.StatusTemporaryRedirect)
	return nil
}
