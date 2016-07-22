import qwest from "qwest";

export default class {
  constructor(token) {
    this._token = token;
    this._endpoint = '/api';
  }

  update(fullName) {
    return qwest.post(`${this._endpoint}/${fullName}/update`, null, {
      headers: {Authorization: `token ${this._token}`}
    });
  }

  getLatestGradle() {
    return qwest.get(`${this._endpoint}/latestGradle`);
  }

  findPullRequests() {
    return qwest.get(`${this._endpoint}/pullRequests`);
  }
}
