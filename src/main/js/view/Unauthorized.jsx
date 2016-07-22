import React from "react";
import GradleUpdate from "../repository/GradleUpdate.jsx";
import OAuthSession from "../repository/OAuthSession.jsx";
import GUPullRequests from "./GUPullRequests.jsx";
import Footer from "./Footer.jsx";

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.gradleUpdate = new GradleUpdate(this.props.token);
  }

  componentDidMount() {
    this.gradleUpdate.findPullRequests()
      .then((xhr, pullRequests) => this.setState({pullRequests: pullRequests}))
      .catch((e) => this.setState({error: e}));
  }

  render() {
    return (
      <div>
        {this.props.error ? (<Error message={this.props.error}/>) : null}

        <div className="container">
          <section className="text-center">
            <div className="page-header">
              <h1>Gradle Update</h1>
              <p>keeps the latest Gradle wrapper on your GitHub repositories</p>
            </div>

            <button className="btn btn-default" onClick={this.props.onAuthorize}>
              Sign in with GitHub Account
            </button>
          </section>

          <section className="text-center">
            <h2>Contributions</h2>
            <GUPullRequests pullRequests={this.state.pullRequests}/>
          </section>

          <Footer/>

          {location.hostname == 'localhost' ? (<OAuthTokenManipulator/>) : null}
        </div>
      </div>
    );
  }
}

class Error extends React.Component {
  render() {
    return (
      <div className="alert alert-warning" role="alert">
        <strong>Authorization Error:</strong> {this.props.message}
      </div>
    );
  }
}

class OAuthTokenManipulator extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};
  }
  handleChange(e) {
    this.setState({oauthToken: e.target.value});
  }
  signIn(e) {
    e.preventDefault();
    OAuthSession.saveToken(this.state.oauthToken);
    location.reload();
  }
  render() {
    return (
      <div className="well text-center">
        <form onSubmit={this.signIn.bind(this)}>
          <div className="form-group">
            <input type="text" className="form-control"
                   placeholder="OAuth Token for Development"
                   required
                   onChange={this.handleChange.bind(this)}
                   value={this.state.oauthToken}/>
          </div>
          <button className="btn btn-default" disabled={!this.state.oauthToken}>
            Sign In By OAuth Token
          </button>
        </form>
      </div>
    );
  }
}
