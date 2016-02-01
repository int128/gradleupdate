import React from 'react';

import Constants from '../Constants.jsx';

import qwest from 'qwest';
import queryString from 'query-string';

import Authorized from './Authorized.jsx';
import Unauthorized from './Unauthorized.jsx';

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
          this.setState({oauth: {state: 'GotError', error: 'OAuth state parameter did not match'}});
          history.replaceState(null, null, '/');
        }
      } else if (params.error_description) {
        this.setState({oauth: {state: 'GotError', error: params.error_description}});
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
    return (<Unauthorized onAuthorize={this.authorize.bind(this)}/>);
  }
  renderGotCode() {
    return (<GotCode code={this.state.oauth.code}
      onGotToken={this.onGotToken.bind(this)}
      onGotError={this.onGotError.bind(this)}/>);
  }
  renderGotError() {
    return (<GotError error={this.state.oauth.error}/>);
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
      + `&redirect_uri=${location.origin}`
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
    this.setState({oauth: {state: 'GotError', error: e}});
  }
}

class GotCode extends React.Component {
  componentDidMount() {
    qwest.post('/authorize', {code: this.props.code})
      .then((xhr, response) => this.props.onGotToken(response.token))
      .catch((e) => this.props.onGotError(e));
  }
  render() {
    return (
      <div className="container">
        <h2>Authorization in Progress</h2>
      </div>
    );
  }
}

class GotError extends React.Component {
  render() {
    return (
      <div className="container">
        {this.props.oauthError}
      </div>
    );
  }
}
