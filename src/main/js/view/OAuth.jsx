import React from "react";
import {Link} from "react-router";
import OAuthSession from "../repository/OAuthSession.jsx";
import ErrorHeader from "./ErrorHeader.jsx";
import Footer from "./Footer.jsx";

export class SignIn extends React.Component {
  componentDidMount() {
    const query = this.props.location.query;
    if (query.code && OAuthSession.validateState(query.state)) {
      this.props.history.replace({pathname: '/signin/exchange', state: query.code});
    } else if (query.code) {
      this.props.history.replace({pathname: '/signin/error', state: 'OAuth state does not match'});
    } else if (query.error_description) {
      this.props.history.replace({pathname: '/signin/error', state: query.error_description});
    } else {
      OAuthSession.authorize(this.props.location.pathname);
    }
  }
  render() {
    return (<Authorizing progress="33%"/>);
  }
}

export class SignInExchange extends React.Component {
  componentDidMount() {
    OAuthSession.exchangeCodeAndToken(this.props.location.state)
      .then((xhr, response) => {
        OAuthSession.saveToken(response.token);
        this.props.history.replace({pathname: '/'});
      })
      .catch((e) =>
        this.props.history.replace({pathname: '/signin/error', state: e}));
  }
  render() {
    return (<Authorizing progress="67%"/>);
  }
}

class Authorizing extends React.Component {
  render() {
    return (
      <div className="container">
        <section>
          <div className="text-center">
            Authorizing...
          </div>
          <div className="progress">
            <div className="progress-bar progress-bar-striped active"
                 style={{width: this.props.progress}}></div>
          </div>
        </section>
        <Footer/>
      </div>
    );
  }
}

export class SignInError extends React.Component {
  render() {
    return (
      <div className="container">
        <ErrorHeader kind="OAuth Error" message={this.props.location.state || 'Unknown'}/>
        <section className="text-center">
          <Link to="/" className="btn btn-default">Back</Link>
        </section>
        <Footer/>
      </div>
    );
  }
}

export class SignOut extends React.Component {
  componentDidMount() {
    OAuthSession.expireToken();
    this.props.history.replace({pathname: '/'});
  }
  render() {
    return (<div></div>);
  }
}
