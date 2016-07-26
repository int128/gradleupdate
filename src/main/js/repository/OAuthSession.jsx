import qwest from "qwest";
import pinkySwear from "pinkyswear";
import Constants from "../config/Constants.jsx";

export default {
  getToken() {
    return localStorage.getItem('oauthToken');
  },

  expireToken() {
    localStorage.removeItem('oauthToken');
  },

  authorize() {
    const state = Math.random().toString(36).substring(2);
    const redirectURI = `${location.origin}${location.pathname}`;
    sessionStorage.setItem('oauthState', state);
    sessionStorage.setItem('oauthRedirectURI', redirectURI);
    location.replace('https://github.com/login/oauth/authorize'
      + `?client_id=${Constants.oauthClientId}`
      + `&redirect_uri=${redirectURI}`
      + `&scope=${Constants.oauthScope}`
      + `&state=${state}`);
  },

  validateState(state) {
    return sessionStorage.getItem('oauthState') == state;
  },

  exchangeCodeAndToken(code) {
    return qwest.post('/api/exchange-oauth-token', {
      code: code,
      state: sessionStorage.getItem('oauthState'),
      redirect_uri: sessionStorage.getItem('oauthRedirectURI')
    }).then((xhr, response) => {
      const promise = pinkySwear();
      if (response.access_token) {
        localStorage.setItem('oauthToken', response.access_token);
        promise(true);
      } else if (response.error_description) {
        promise(false, response.error_description);
      } else {
        promise(false);
      }
      return promise;
    });
  }
}
