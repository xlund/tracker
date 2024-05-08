async function handleSubmit(e) {
  e.preventDefault();
  const formData = new FormData(e.target);
  const url = e.target.action;

  const credentialOptions = await requestCredentialOptions(formData, url);

  const clientData = await createCredential(credentialOptions);

  const passkey = await fetch("/auth/passkeys/complete", {
    method: "POST",
    body: JSON.stringify(clientData.response),
  });
}

async function createCredential(credentialOptions) {
  const { publicKey } = credentialOptions;
  if (!publicKey) return Error("Passkey not found");

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
