package frontend

import (
	"bytes"
	"image/png"
	"net/url"
	"testing"
	"time"

	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/mxschmitt/playwright-go"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultBrowserWaitTime = time.Second / 2

func TestRegistrationFlow(T *testing.T) {
	helper := setupTestHelper(T)

	helper.runForAllBrowsers(T, "registration flow", func(browser playwright.Browser) func(*testing.T) {
		return func(t *testing.T) {
			// register a user
			user := fakes.BuildFakeUserCreationInput()

			page, err := browser.NewPage()
			require.NoError(t, err, "could not create page")

			_, err = page.Goto(urlToUse + "/register")
			require.NoError(t, err, "could not navigate to registration page")

			require.NoError(t, page.Type("[data-automation=registration_username_input]", user.Username))
			require.NoError(t, page.Type("[data-automation=registration_password_input]", user.Password))
			require.NoError(t, page.Type("[data-automation=registration_password_repeat_input]", user.Password))

			time.Sleep(defaultBrowserWaitTime)

			require.NoError(t, page.Click("[data-automation=register_button]"))

			time.Sleep(defaultBrowserWaitTime)

			qrCodeElement, qrCodeElementErr := page.QuerySelector("[data-automation=two_factor_qr_code]")
			require.NoError(t, qrCodeElementErr)

			img, err := png.Decode(bytes.NewReader(getScreenshotBytes(t, qrCodeElement)))
			require.NoError(t, err)

			// prepare BinaryBitmap
			bmp, bitmapErr := gozxing.NewBinaryBitmapFromImage(img)
			require.NoError(t, bitmapErr)

			// decode image
			qrReader := qrcode.NewQRCodeReader()
			result, qrCodeDecodeErr := qrReader.Decode(bmp, nil)
			require.NoError(t, qrCodeDecodeErr)

			u, secretParseErr := url.Parse(result.String())
			require.NoError(t, secretParseErr)
			twoFactorSecret := u.Query().Get("secret")
			require.NotEmpty(t, twoFactorSecret)

			code, firstCodeGenerationErr := totp.GenerateCode(twoFactorSecret, time.Now().UTC())
			require.NoError(t, firstCodeGenerationErr)
			require.NotEmpty(t, code)

			totpInputFieldFindErr := page.Type("[data-automation=totp_secret_verification_token_input]", code)
			require.NoError(t, totpInputFieldFindErr, "unexpected error finding TOTP token input field: %v", totpInputFieldFindErr)

			require.NoError(t, page.Click("[data-automation=totp_token_submit_button]"))

			time.Sleep(defaultBrowserWaitTime)
			assert.Equal(t, urlToUse+"/login", page.URL())

			// login with the newly registered user

			code, secondCodeGenerationErr := totp.GenerateCode(twoFactorSecret, time.Now().UTC())
			require.NoError(t, secondCodeGenerationErr)
			require.NotEmpty(t, code)

			require.NoError(t, page.Type("[data-automation=login_username_input]", user.Username))
			require.NoError(t, page.Type("[data-automation=login_password_input]", user.Password))
			require.NoError(t, page.Type("[data-automation=login_totp_token_input]", code))

			require.NoError(t, page.Click("[data-automation=login_button]"))
			time.Sleep(defaultBrowserWaitTime)

			assert.Equal(t, urlToUse+"/home", page.URL())
		}
	})
}
