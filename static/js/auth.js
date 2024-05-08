async function handleSubmit(e) {
  e.preventDefault();
  const formData = new FormData(e.target);
  const url = e.target.action;

  const credentialOptions = await requestCredentialOptions(formData, url);

  const clientData = await createCredential(credentialOptions);

  // Stringifying the JSON removes the data in ArrayBuffers
  // on the PublicKeyCredential object and sets it to empty objects
  const passkey = await fetch("/auth/passkeys/complete", {
    method: "POST",
    body: JSON.stringify(clientData.response),
  });
}

async function createCredential(credentialOptions) {
  const { publicKey } = credentialOptions;
  if (!publicKey) return Error("Passkey not found");

  // This successfuly promts the user to sign the credential
  const credential = await navigator.credentials.create({
    publicKey: {
      ...publicKey,
      challenge: new Uint8Array(publicKey.challenge),
      user: {
        ...publicKey.user,
        id: new Uint8Array(publicKey.user.id),
      },
    },
  });

  // Creates a PublicKeyCredential object
  //   PublicKeyCredential {
  //     id: '...',
  //     rawId: ArrayBuffer(59),
  //     response: AuthenticatorAttestationResponse {
  //         clientDataJSON: ArrayBuffer(121),
  //         attestationObject: ArrayBuffer(306),
  //     },
  //     type: 'public-key'
  // }
  console.log(credential);
  return credential;
}

async function requestCredentialOptions(formData, url) {
  try {
    const res = await fetch(url, {
      method: "POST",
      body: formData,
    });
    const credential = await res.json();
    return credential;
  } catch (e) {
    console.error(e);
  }
}

addEventListener("submit", handleSubmit);
