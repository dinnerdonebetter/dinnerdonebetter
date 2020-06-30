package frontend

import (
	"bytes"
	"image/png"
	"net/url"
	"testing"
	"time"

	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
	"github.com/tebeka/selenium/log"
)

func runTestOnAllSupportedBrowsers(t *testing.T, tp testProvider) {
	firefoxCaps, chromeCaps := selenium.Capabilities{}, selenium.Capabilities{}

	firefoxCaps.AddFirefox(firefox.Capabilities{
		Log: &firefox.Log{Level: firefox.Debug},
	})

	capMap := map[string]selenium.Capabilities{
		"firefox": firefoxCaps,
		"chrome":  chromeCaps,
	}

	for bn, caps := range capMap {
		caps.AddLogging(log.Capabilities{
			log.Server:      log.Debug,
			log.Browser:     log.Debug,
			log.Client:      log.Debug,
			log.Driver:      log.Debug,
			log.Performance: log.Debug,
			log.Profiler:    log.Debug,
		})
		caps["browserName"] = bn

		wd, err := selenium.NewRemote(caps, seleniumHubAddr)
		if err != nil {
			panic(err)
		}

		t.Run(bn, tp(wd))
		assert.NoError(t, wd.Quit())
	}
}

/*
func saveScreenshotTo(t *testing.T, driver selenium.WebDriver, path string) {
	t.Helper()

	screenshotAsBytes, err := driver.Screenshot()
	require.NoError(t, err)

	im, err := png.Decode(bytes.NewReader(screenshotAsBytes))
	require.NoError(t, err)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_CREATE, 0744)
	require.NoError(t, err)

	require.NoError(t, png.Encode(f, im))
}
*/

type testProvider func(driver selenium.WebDriver) func(t *testing.T)

func TestLoginPage(T *testing.T) {
	runTestOnAllSupportedBrowsers(T, func(driver selenium.WebDriver) func(t *testing.T) {
		return func(t *testing.T) {
			// Navigate to the login page.
			reqURI := urlToUse + "/login"
			require.NoError(t, driver.Get(reqURI))

			time.Sleep(time.Second)

			require.NoError(t, driver.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
				elem, err := driver.FindElement(selenium.ByID, "loginButton")
				if err != nil {
					return false, err
				}
				return elem.IsDisplayed()
			}, 10*time.Second, time.Second))

			// fetch the button.
			loginButton, loginButtonFindErr := driver.FindElement(selenium.ByID, "loginButton")
			require.NoError(t, loginButtonFindErr)

			// check that it is visible.
			actual, isDisplayedErr := loginButton.IsDisplayed()
			assert.NoError(t, isDisplayedErr)
			assert.True(t, actual)
		}
	})
}

func TestRegistrationFlow(T *testing.T) {
	const (
		loginButtonID           = "loginButton"
		registrationButtonID    = "registrationButton"
		usernameInputID         = "usernameInput"
		passwordInputID         = "passwordInput"
		passwordRepeatInputID   = "passwordRepeatInput"
		twoFactorSecretQRCodeID = "twoFactorSecretQRCode"
		totpTokenSubmitButtonID = "totpTokenSubmitButton"
		totpTokenInputID        = "totpTokenInput"
	)

	runTestOnAllSupportedBrowsers(T, func(driver selenium.WebDriver) func(t *testing.T) {
		return func(t *testing.T) {
			// Navigate to the registration page.
			reqURI := urlToUse + "/register"
			require.NoError(t, driver.Get(reqURI))

			time.Sleep(time.Second)

			require.NoError(t, driver.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
				elem, err := wd.FindElement(selenium.ByID, registrationButtonID)
				if err != nil {
					return false, err
				}
				return elem.IsDisplayed()
			}, 10*time.Second, time.Second))

			user := fakemodels.BuildFakeUserCreationInput()

			// fetch the username field and fill it
			usernameField, usernameFieldFindErr := driver.FindElement(selenium.ByID, usernameInputID)
			require.NoError(t, usernameFieldFindErr)
			require.NoError(t, usernameField.SendKeys(user.Username))

			// fetch the password field and fill it
			passwordField, passwordFieldFindErr := driver.FindElement(selenium.ByID, passwordInputID)
			require.NoError(t, passwordFieldFindErr)
			require.NoError(t, passwordField.SendKeys(user.Password))

			// fetch the password confirm field and fill it
			passwordRepeatField, passwordRepeatFieldFindErr := driver.FindElement(selenium.ByID, passwordRepeatInputID)
			require.NoError(t, passwordRepeatFieldFindErr)
			require.NoError(t, passwordRepeatField.SendKeys(user.Password))

			// fetch the button.
			registerButton, registerButtonFindErr := driver.FindElement(selenium.ByID, registrationButtonID)
			require.NoError(t, registerButtonFindErr)
			require.NoError(t, registerButton.Click())

			time.Sleep(time.Second)

			require.NoError(t, driver.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
				qrCode, err := wd.FindElement(selenium.ByID, twoFactorSecretQRCodeID)
				if err != nil {
					return false, err
				}
				return qrCode.IsDisplayed()
			}, 10*time.Second, time.Second))

			// check that it is visible.
			qrCode, twoFactorQRCodeFindErr := driver.FindElement(selenium.ByID, twoFactorSecretQRCodeID)
			assert.NoError(t, twoFactorQRCodeFindErr)
			qrCodeIsDisplayed, qrCodeIsDisplayedErr := qrCode.IsDisplayed()
			require.NoError(t, qrCodeIsDisplayedErr)
			require.True(t, qrCodeIsDisplayed)

			qrScreenshotBytes, qrCodeScreenshotErr := qrCode.Screenshot(false)
			require.NoError(t, qrCodeScreenshotErr)

			img, err := png.Decode(bytes.NewReader(qrScreenshotBytes))
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

			// fetch the totp confirmation field and fill it
			totpTokenInputField, totpInputFieldFindErr := driver.FindElement(selenium.ByID, totpTokenInputID)
			require.NoError(t, totpInputFieldFindErr)
			require.NoError(t, totpTokenInputField.SendKeys(code))

			// fetch the button.
			totpTokenSubmitButton, err := driver.FindElement(selenium.ByID, totpTokenSubmitButtonID)
			require.NoError(t, err)
			require.NoError(t, totpTokenSubmitButton.Click())

			time.Sleep(3 * time.Second)

			expectedURL := urlToUse + "/login"
			actualURL, err := driver.CurrentURL()
			require.NoError(t, err)
			assert.Equal(t, expectedURL, actualURL, "expected %q to equal %q", actualURL, expectedURL)

			// fetch the username field and fill it
			usernameField, usernameFieldFindErr = driver.FindElement(selenium.ByID, usernameInputID)
			require.NoError(t, usernameFieldFindErr)
			require.NoError(t, usernameField.SendKeys(user.Username))

			// fetch the password field and fill it
			passwordField, passwordFieldFindErr = driver.FindElement(selenium.ByID, passwordInputID)
			require.NoError(t, passwordFieldFindErr)
			require.NoError(t, passwordField.SendKeys(user.Password))

			code, secondCodeGenerationErr := totp.GenerateCode(twoFactorSecret, time.Now().UTC())
			require.NoError(t, secondCodeGenerationErr)
			require.NotEmpty(t, code)

			// fetch the TOTP code field and fill it
			totpTokenInputField, totpTokenInputFieldErr := driver.FindElement(selenium.ByID, totpTokenInputID)
			require.NoError(t, totpTokenInputFieldErr)
			require.NoError(t, totpTokenInputField.SendKeys(code))

			// fetch the button.
			loginButton, loginButtonFindErr := driver.FindElement(selenium.ByID, loginButtonID)
			require.NoError(t, loginButtonFindErr)
			require.NoError(t, loginButton.Click())

			time.Sleep(3 * time.Second)

			expectedURL = urlToUse + "/"
			actualURL, err = driver.CurrentURL()
			require.NoError(t, err)
			assert.Equal(t, expectedURL, actualURL, "expected final URL %q to equal %q", actualURL, expectedURL)
		}
	})
}
