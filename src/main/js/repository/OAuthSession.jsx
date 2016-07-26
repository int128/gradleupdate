import qwest from "qwest";
import Constants from "../config/Constants.jsx";

export default {
  getToken() {
    return localStorage.getItem('oauthToken');
  },

  saveToken(token) {
    localStorage.setItem('oauthToken', token);
  },

  expireToken() {
    localStorage.removeItem('oauthToken');
  },

  authorize(redirectURI) {
    const state = Math.random().toString(36).substring(2);
    sessionStorage.setItem('oauthState', state);
    location.replace('https://github.com/login/oauth/authorize'
      + `?client_id=${Constants.oauthClientId}`
      + `&redirect_uri=${location.origin}${redirectURI}`
      + `&scope=${Constants.oauthScope}`
      + `&state=${state}`);
  },

  validateState(state) {
    return sessionStorage.getItem('oauthState') == state;
  },

  exchangeCodeAndToken(code) {
    return qwest.post('/api/exchange-oauth-token', {code: code});
  }
}
