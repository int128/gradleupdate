import qwest from 'qwest';

export default class {
  constructor(token) {
    this._token = token;
    this._endpoint = 'https://api.github.com';
  }
  getUser() {
    return qwest.get(`${this._endpoint}/user`, null, {
      headers: {Authorization: `token ${this._token}`},
      cache: true,  // prevent Cache-Control for CORS
    });
  }
}
