package components

// PasskeyRegistrationScript is the JavaScript for the passkey registration flow.
// #nosec G101 -- JavaScript source for passkey flow, not credentials
const PasskeyRegistrationScript = `
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
