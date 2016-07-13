import qwest from "qwest";

export default class {
  constructor(token) {
    this._token = token;
    this._endpoint = '/api';
  }
  update(fullName, gradleVersion) {
    return qwest.post(`${this._endpoint}/${fullName}/update`, {
      gradle_version: gradleVersion
    }, {
      headers: {Authorization: `token ${this._token}`}
    });
  }
}
