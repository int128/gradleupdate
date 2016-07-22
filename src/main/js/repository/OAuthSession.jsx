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
  }
}
