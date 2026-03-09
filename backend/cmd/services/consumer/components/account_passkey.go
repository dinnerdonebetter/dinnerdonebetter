package components

import (
	g "maragu.dev/gomponents"
	ghtml "maragu.dev/gomponents/html"
)

// passkeyRegistrationScript is the JavaScript for the passkey registration flow.
// #nosec G101 -- JavaScript source for passkey flow, not credentials
const passkeyRegistrationScript = `
(function() {
	if (!window.PublicKeyCredential) return;
	var btn = document.getElementById('add-passkey-btn');
	if (!btn) return;
	function b64enc(buf) {
		var b = new Uint8Array(buf);
		var s = '';
		for (var i = 0; i < b.length; i++) s += String.fromCharCode(b[i]);
		return btoa(s).replace(/\+/g,'-').replace(/\//g,'_').replace(/=+$/,'');
	}
	function b64dec(s) {
		s = (s + '==='.slice((s.length + 3) % 4)).replace(/-/g,'+').replace(/_/g,'/');
		return Uint8Array.from(atob(s), function(c){return c.charCodeAt(0);});
	}
	btn.onclick = function() {
		btn.disabled = true;
		fetch('/auth/passkey/registration/options', {
			method: 'POST',
			headers: {'Content-Type': 'application/json'},
			body: JSON.stringify({}),
			credentials: 'include'
		}).then(function(r) {
			if (!r.ok) throw new Error('Failed to get options');
			return r.json();
		}).then(function(opts) {
			var raw = atob(opts.publicKeyCredentialCreationOptions);
			var obj = JSON.parse(raw);
			var pk = obj.publicKey || obj;
			if (typeof pk.challenge === 'string') pk.challenge = b64dec(pk.challenge).buffer;
			if (pk.user && typeof pk.user.id === 'string') pk.user.id = b64dec(pk.user.id).buffer;
			return navigator.credentials.create({publicKey: pk}).then(function(cred) {
				return {cred: cred, opts: opts};
			});
		}).then(function(data) {
			var cred = data.cred;
			if (!cred) throw new Error('No credential');
			var r = cred.response;
			var attestation = {
				id: cred.id,
				rawId: b64enc(cred.rawId),
				type: cred.type,
				response: {
					clientDataJSON: b64enc(r.clientDataJSON),
					attestationObject: b64enc(r.attestationObject)
				}
			};
			return fetch('/auth/passkey/registration/verify', {
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify({
					attestationResponse: attestation,
					challenge: data.opts.challenge
				}),
				credentials: 'include'
			});
		}).then(function(r) {
			if (!r.ok) throw new Error('Registration failed');
			window.location.reload();
		}).catch(function(err) {
			btn.disabled = false;
			alert(err.message || 'Passkey registration failed');
		});
	};
})();
`

// PasskeySection renders the "Add passkey" section for the account page.
// Styled consistently with the account link cards above it.
func (r *ComponentRenderer) PasskeySection() g.Node {
	return ghtml.Div(
		ghtml.Class("block p-4 rounded-lg border border-gray-200 bg-white"),
		ghtml.Div(
			ghtml.Class("font-medium"),
			g.Text("Passkeys"),
		),
		ghtml.Div(
			ghtml.Class("text-sm text-gray-500 mt-0.5"),
			g.Text("Add a passkey to sign in quickly without a password."),
		),
		ghtml.Button(
			ghtml.Type("button"),
			ghtml.ID("add-passkey-btn"),
			ghtml.Class("mt-3 inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"),
			g.Text("Add passkey"),
		),
		ghtml.Script(ghtml.Type("text/javascript"), g.Raw(passkeyRegistrationScript)),
	)
}
