import React from "react";
import Constants from "../Constants.jsx";
import qwest from "qwest";
import queryString from "query-string";
import Authorized from "./Authorized.jsx";
import Unauthorized from "./Unauthorized.jsx";

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.state = {oauth: {}};
  }

  componentDidMount() {
    const token = sessionStorage.getItem('oauthToken');
    if (token) {
      this.setState({oauth: {state: 'Authorized', token: token}});
    } else if (location.search) {
      const params = queryString.parse(location.search);
      if (params.state && params.code) {
        if (sessionStorage.getItem('oauthKey') == params.state) {
          this.setState({oauth: {state: 'GotCode', code: params.code}});
          history.replaceState(null, null, '/');
        } else {
          this.setState({oauth: {state: 'Unauthorized', error: 'OAuth state parameter did not match'}});
          history.replaceState(null, null, '/');
        }
      } else if (params.error_description) {
        this.setState({oauth: {state: 'Unauthorized', error: params.error_description}});
        history.replaceState(null, null, '/');
      } else {
        this.setState({oauth: {state: 'Unauthorized'}});
      }
    } else {
      this.setState({oauth: {state: 'Unauthorized'}});
    }
  }

  render() {
    const renderer = this[`render${this.state.oauth.state}`];
    return renderer ? renderer.apply(this) : null;
  }
  renderUnauthorized() {
    return (<Unauthorized error={this.state.oauth.error}
      onAuthorize={this.authorize.bind(this)}/>);
  }
  renderGotCode() {
    return (<GotCode code={this.state.oauth.code}
      onGotToken={this.onGotToken.bind(this)}
      onGotError={this.onGotError.bind(this)}/>);
  }
  renderAuthorized() {
    return (<Authorized token={this.state.oauth.token}
      onUnauthorize={this.unauthorize.bind(this)}/>);
  }

  authorize() {
    const key = Math.random().toString(36).substring(2);
    const url = 'https://github.com/login/oauth/authorize'
      + `?client_id=${Constants.oauthClientId}`
      + `&scope=${Constants.oauthScope}`
      + `&state=${key}`;
    sessionStorage.setItem('oauthKey', key);
    location.replace(url);
  }
  unauthorize() {
    sessionStorage.removeItem('oauthToken');
    this.setState({oauth: {state: 'Unauthorized'}});
  }
  onGotToken(token) {
    sessionStorage.removeItem('oauthKey');
    sessionStorage.setItem('oauthToken', token);
    this.setState({oauth: {state: 'Authorized', token: token}});
  }
  onGotError(e) {
    this.setState({oauth: {state: 'Unauthorized', error: e}});
  }
}

class GotCode extends React.Component {
  componentDidMount() {
    qwest.post('/api/authorize', {code: this.props.code})
      .then((xhr, response) => this.props.onGotToken(response.token))
      .catch((e) => this.props.onGotError(e));
  }
  render() {
    return (
      <div className="container">
        <div className="jumbotron">
          <p className="text-center">Authorizing...</p>
          <div className="progress">
            <div className="progress-bar progress-bar-striped active" style={{width: '50%'}}></div>
          </div>
        </div>
      </div>
    );
  }
}
