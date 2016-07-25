import qwest from "qwest";
import Constants from "../Constants.jsx";

export default {
  getToken() {
    return localStorage.getItem('oauthToken');
  },

  saveToken(token) {
    sessionStorage.removeItem('oauthKey');
    localStorage.setItem('oauthToken', token);
  },

  expireToken() {
    localStorage.removeItem('oauthToken');
  },

  validateKey(key) {
    return sessionStorage.getItem('oauthKey') == key;
  },

  saveKey(key) {
    sessionStorage.setItem('oauthKey', key);
  },

  redirectToAuthorize() {
    const key = Math.random().toString(36).substring(2);
    const url = 'https://github.com/login/oauth/authorize'
      + `?client_id=${Constants.oauthClientId}`
      + `&scope=${Constants.oauthScope}`
      + `&state=${key}`;
    this.saveKey(key);
    location.replace(url);
  },

  exchangeCodeAndToken(code) {
    return qwest.post('/api/exchange-oauth-token', {code: code});
  }
}
