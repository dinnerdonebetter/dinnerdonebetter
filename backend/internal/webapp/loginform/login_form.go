package loginform

import (
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/webapp/design"

	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// Props holds the login form state and error messages.
type Props struct {
	UsernameError,
	PasswordError,
	TOTPError,
	GeneralError string
	// ResetSuccessMessage is shown when the user has successfully reset their password (e.g. "Your password has been reset. Sign in with your new password.").
	ResetSuccessMessage string
}

// Config allows customizing the form copy per app.
type Config struct {
	// Title is the form heading (e.g. "Login" or "Sign In").
	Title string
	// SubmitButtonText is the submit button label (e.g. "Log In" or "Sign In").
	SubmitButtonText string
	// ForgotPasswordLink is the href for the "Forgot password?" link. If empty, the link is not shown.
	ForgotPasswordLink string
}

// DefaultConfig returns config with "Login" / "Log In" copy.
func DefaultConfig() Config {
	return Config{
		Title:            "Login",
		SubmitButtonText: "Log In",
	}
}

// SignInConfig returns config with "Sign In" copy.
func SignInConfig() Config {
	return Config{
		Title:              "Sign In",
		SubmitButtonText:   "Sign In",
		ForgotPasswordLink: "/forgot_password",
	}
}

// Form renders a login form with username, password, and TOTP fields.
// Uses HTMX for submission. Pass empty strings in props for no errors.
func Form(props *Props, cfg Config, palette *design.Palette) g.Node {
	if props == nil {
		props = &Props{}
	}
	if palette == nil {
		palette = &design.StandardPalette
	}
	if cfg.Title == "" {
		cfg.Title = "Login"
	}
	if cfg.SubmitButtonText == "" {
		cfg.SubmitButtonText = "Log In"
	}

	return ghtml.Div(
		ghtml.ID("login-container"),
		ghtml.Div(
			ghtml.Class("w-full max-w-md bg-white p-8 rounded-2xl shadow-md"),
			ghtml.H2(
				ghtml.Class(fmt.Sprintf("text-2xl font-bold mb-6 text-center %s", design.TextColor(palette.Primary))),
				g.Text(cfg.Title),
			),

			g.If(props.ResetSuccessMessage != "", ghtml.Div(
				ghtml.Class(fmt.Sprintf("mb-4 p-3 rounded-md text-sm bg-green-50 border border-green-200 %s", design.TextColor(palette.Primary))),
				g.Text(props.ResetSuccessMessage),
			)),

			ghtml.Form(
				ghtml.Class("space-y-4"),
				ghtml.Method("post"),
				g.Attr("hx-post", "/login/submit"),
				g.Attr("hx-ext", "json-enc"),
				g.Attr("hx-target", "#login-container"),
				g.Attr("hx-swap", "outerHTML"),
				g.Attr("hx-request", `{"credentials":"include"}`),

				wrapInputElement("username", props.UsernameError, usernameInput("username", "username", "", palette), palette),
				wrapInputElement("password", props.PasswordError, passwordInput("password", "password", "", palette), palette),
				g.If(cfg.ForgotPasswordLink != "", ghtml.Div(
					ghtml.Class("text-right -mt-2"),
					ghtml.A(
						ghtml.Href(cfg.ForgotPasswordLink),
						ghtml.Class(fmt.Sprintf("text-sm %s hover:underline", design.TextColor(palette.Primary))),
						g.Text("Forgot password?"),
					),
				)),
				wrapInputElement("TOTP code", props.TOTPError, totpTokenInput("totp", "totpToken", "", palette), palette),

				submitButton(cfg.SubmitButtonText),

				ghtml.Div(
					ghtml.Class("relative my-4"),
					ghtml.Div(ghtml.Class("absolute inset-0 flex items-center"),
						ghtml.Div(ghtml.Class("w-full border-t border-gray-300")),
					),
					ghtml.Div(ghtml.Class("relative flex justify-center text-sm"),
						ghtml.Span(ghtml.Class("px-2 bg-white text-gray-500"), g.Text("or")),
					),
				),

				ghtml.Div(
					ghtml.Class("mt-4"),
					ghtml.Button(
						ghtml.Type("button"),
						ghtml.ID("passkey-sign-in-btn"),
						ghtml.Class("w-full py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 flex items-center justify-center gap-2"),
						g.Text("Sign in with passkey"),
					),
				),
				ghtml.Script(ghtml.Type("text/javascript"), g.Raw(passkeyScript)),

				g.If(props.GeneralError != "", ghtml.Div(
					ghtml.Class(fmt.Sprintf("mt-2 text-sm %s", design.TextColor(palette.Warning))),
					g.Text(props.GeneralError),
				)),
			),
		),
	)
}

// passkeyScript is the JavaScript for the passkey sign-in flow.
// #nosec G101 -- JavaScript source for passkey flow, not credentials
const passkeyScript = `
(function() {
	if (!window.PublicKeyCredential) return;
	var btn = document.getElementById('passkey-sign-in-btn');
	if (!btn) return;
	function b64enc(buf) {
		var b = new Uint8Array(buf);
		var s = '';
		for (var i = 0; i < b.length; i++) s += String.fromCharCode(b[i]);
		return btoa(s).replace(/\\+/g,'-').replace(/\\//g,'_').replace(/=+$/,'');
	}
	function b64dec(s) {
		s = (s + '==='.slice((s.length + 3) % 4)).replace(/-/g,'+').replace(/_/g,'/');
		return Uint8Array.from(atob(s), function(c){return c.charCodeAt(0);});
	}
	btn.onclick = function() {
		var usernameEl = document.getElementById('username');
		var username = usernameEl ? usernameEl.value.trim() : '';
		btn.disabled = true;
		fetch('/auth/passkey/authentication/options', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({username: username})
		}).then(function(r) {
			if (!r.ok) throw new Error('Failed to get options');
			return r.json();
		}).then(function(opts) {
			var raw = atob(opts.publicKeyCredentialRequestOptions);
			var obj = JSON.parse(raw);
			var pk = obj.publicKey || obj;
			if (typeof pk.challenge === 'string') pk.challenge = b64dec(pk.challenge).buffer;
			if (pk.allowCredentials) {
				pk.allowCredentials = pk.allowCredentials.map(function(c) {
					return {id: c.id, type: c.type || 'public-key', transports: c.transports};
				});
			}
			return navigator.credentials.get({publicKey: pk}).then(function(cred) {
				return {cred: cred, opts: opts, username: username};
			});
		}).then(function(data) {
			var cred = data.cred;
			if (!cred) throw new Error('No credential');
			var r = cred.response;
			var assertion = {
				id: cred.id,
				rawId: b64enc(cred.rawId),
				type: cred.type,
				response: {
					clientDataJSON: b64enc(r.clientDataJSON),
					authenticatorData: b64enc(r.authenticatorData),
					signature: b64enc(r.signature),
					userHandle: r.userHandle ? b64enc(r.userHandle) : null
				}
			};
			return fetch('/auth/passkey/authentication/verify', {
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify({
					assertionResponse: assertion,
					challenge: data.opts.challenge,
					username: data.username
				})
			});
		}).then(function(r) {
			if (!r.ok) throw new Error('Verification failed');
			window.location.href = '/';
		}).catch(function(err) {
			btn.disabled = false;
			alert(err.message || 'Passkey sign-in failed');
		});
	};
})();
`

func usernameInput(label, fieldName, content string, palette *design.Palette) g.Node {
	return ghtml.Input(
		ghtml.Type("text"),
		ghtml.ID(label),
		ghtml.Name(fieldName),
		ghtml.Value(content),
		ghtml.Class(inputClass(palette)),
		ghtml.AutoComplete("username"),
	)
}

func passwordInput(id, fieldName, content string, palette *design.Palette) g.Node {
	return ghtml.Input(
		ghtml.Type("password"),
		ghtml.ID(id),
		ghtml.Name(fieldName),
		ghtml.Value(content),
		ghtml.Class(inputClass(palette)),
		ghtml.AutoComplete("current-password"),
	)
}

func totpTokenInput(id, fieldName, content string, palette *design.Palette) g.Node {
	return ghtml.Input(
		ghtml.Type("text"),
		ghtml.ID(id),
		ghtml.Name(fieldName),
		ghtml.Value(content),
		ghtml.Class(inputClass(palette)),
		ghtml.MaxLength("6"),
		g.Attr("inputmode", "numeric"),
		ghtml.Pattern("[0-9]{6}"),
		ghtml.AutoComplete("one-time-code"),
	)
}

func inputClass(palette *design.Palette) string {
	return fmt.Sprintf("mt-1 block w-full rounded-md border-%s shadow-sm focus:border-%s focus:ring-%s",
		palette.Background.Value, palette.Primary.Value, palette.Primary.Value)
}

func wrapInputElement(label, inputError string, input g.Node, palette *design.Palette) g.Node {
	titleLabel := strings.ToUpper(label[:1]) + strings.ToLower(label[1:])
	if label == "TOTP code" {
		titleLabel = "TOTP code"
	}

	return ghtml.Div(
		ghtml.Class("space-y-1"),
		ghtml.Label(
			ghtml.For(label),
			ghtml.Class(fmt.Sprintf("block text-sm font-medium %s", design.TextColor(palette.Primary))),
			g.Text(titleLabel),
		),
		input,
		g.If(inputError != "", ghtml.Span(
			ghtml.Class(fmt.Sprintf("text-sm %s mt-1 block", design.TextColor(palette.Warning))),
			g.Text(inputError),
		)),
	)
}

func submitButton(text string) g.Node {
	return ghtml.Button(
		ghtml.Type("submit"),
		ghtml.Class("w-full py-2 px-4 bg-blue-600 text-white font-semibold rounded-md shadow hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"),
		g.Text(text),
	)
}
