import React from 'react';

import Constants from '../Constants.jsx';

import User from './User.jsx';
import Footer from './Footer.jsx';
import Projects from './Projects.jsx';

import qwest from 'qwest';
import queryString from 'query-string';

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.state = {oauth: {state: null}};
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
        } else {
          this.setState({oauth: {state: 'GotError', error: 'OAuth state parameter did not match'}});
        }
      } else if (params.error_description) {
        this.setState({oauth: {state: 'GotError', error: params.error_description}});
      }
      history.replaceState(null, null, '/');
    }
  }
  onAuthorize(url, key) {
    sessionStorage.setItem('oauthKey', key);
    location.replace(url);
  }
  onGotToken(token) {
    sessionStorage.removeItem('oauthKey');
    sessionStorage.setItem('oauthToken', token);
    this.setState({oauth: {state: 'Authorized', token: token}});
  }
  onGotError(e) {
    this.setState({oauth: {state: 'GotError', error: e}});
  }
  onUnauthorize() {
    sessionStorage.removeItem('oauthToken');
    this.setState({oauth: {state: null}});
  }
  render() {
    if (this.state.oauth.state == 'Authorized') {
      return (<Authorized token={this.state.oauth.token}
        onUnauthorize={this.onUnauthorize.bind(this)}/>);
    } else if (this.state.oauth.state == 'GotCode') {
      return (<GotCode code={this.state.oauth.code}
        onGotToken={this.onGotToken.bind(this)}
        onGotError={this.onGotError.bind(this)}/>);
    } else if (this.state.oauth.state == 'GotError') {
      return (<GotError error={this.state.oauth.error}/>);
    } else {
      return (<Unauthorized onAuthorize={this.onAuthorize.bind(this)}/>);
    }
  }
}

class Unauthorized extends React.Component {
  onClick() {
    const key = Math.random().toString(36).substring(2);
    const url = 'https://github.com/login/oauth/authorize'
      + `?client_id=${Constants.oauthClientId}`
      + `&scope=${Constants.oauthScope}`
      + `&redirect_uri=${location.origin}`
      + `&state=${key}`;
    this.props.onAuthorize(url, key);
  }
  render() {
    return (
      <div className="NotAuthorized">
        <div className="jumbotron">
          <div className="container text-center">
            <h1>Gradle Update</h1>
            <p>keeps the latest Gradle wrapper on your GitHub repositories</p>
          </div>
        </div>
        <div className="container text-center">
          <button className="btn btn-default" onClick={this.onClick.bind(this)}>
            Sign in with GitHub Account
          </button>
        </div>
        <Footer/>
      </div>
    );
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

class Authorized extends React.Component {
  onClick() {
    this.props.onUnauthorize();
  }
  render() {
    return (
      <div className="container">
        <button className="btn btn-default" onClick={this.onClick.bind(this)}>
          Sign Out
        </button>
        <h2>User</h2>
        <User token={this.props.token}/>
        <Projects token={this.props.token}/>
        <Footer/>
      </div>
    );
  }
}
