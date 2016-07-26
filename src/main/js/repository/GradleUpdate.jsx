import qwest from "qwest";

export default class {
  constructor(token) {
    if (token) {
      this._authorization = `token ${token}`;
    }
    this._endpoint = '/api';
  }

  updateRepository(fullName) {
    return qwest.post(`${this._endpoint}/${fullName}/update`, null, {
      headers: {Authorization: this._authorization},
      responseType: 'json'
    });
  }

  findRepository(fullName) {
    return qwest.get(`${this._endpoint}/${fullName}/status`, null, {
      responseType: 'json'
    });
  }

  getLatestGradle() {
    return qwest.get(`${this._endpoint}/latestGradle`, null, {
      responseType: 'json'
    });
  }

  findPullRequests() {
    return qwest.get(`${this._endpoint}/pullRequests`, null, {
      responseType: 'json'
    });
  }
}
