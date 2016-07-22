import React from "react";
import GitHub from "../repository/GitHub.jsx";
import GradleUpdate from "../repository/GradleUpdate.jsx";
import GHAuthenticatedUser from "./GHAuthenticatedUser.jsx";
import GHRepositories from "./GHRepositories.jsx";
import GULatestGradle from "./GULatestGradle.jsx";
import Footer from "./Footer.jsx";

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.github = new GitHub(this.props.token);
    this.gradleUpdate = new GradleUpdate(this.props.token);
  }

  componentDidMount() {
    this.github.getUser()
      .then((xhr, user) => this.setState({user: user}))
      .catch((e) => this.setState({error: e}));
    this.github.findRepositories({sort: 'full_name'})
      .then((xhr, repos) => this.setState({repos: repos}))
      .catch((e) => this.setState({error: e}));
    this.gradleUpdate.getLatestGradle()
      .then((xhr, latestGradle) => this.setState({latestGradle: latestGradle}))
      .catch((e) => this.setState({error: e}));
  }

  updateGradleWrapper(fullName) {
    this.gradleUpdate.update(fullName);
  }

  render() {
    return (
      <div>
        {this.state.error ? (<Error message={this.state.error}/>) : null}

        <div className="page-header text-center">
          <h1>Gradle Update</h1>
          <GULatestGradle
            latestGradle={this.state.latestGradle}/>
        </div>

        <div className="container">
          <div className="row">
            <div className="col-lg-3 col-md-3 col-sm-3">
              <div className="text-center">
                <GHAuthenticatedUser
                  user={this.state.user}
                  signOut={this.props.onUnauthorize.bind(this)}/>
              </div>
            </div>
            <div className="col-lg-9 col-md-9 col-sm-9">
              <GHRepositories
                repos={this.state.repos}
                updateGradleWrapper={this.updateGradleWrapper.bind(this)}
              />
            </div>
          </div>

          <Footer/>
        </div>
      </div>
    );
  }
}

class Error extends React.Component {
  render() {
    return (
      <div className="alert alert-warning" role="alert">
        <strong>Server Error:</strong> {this.props.message}
      </div>
    );
  }
}
