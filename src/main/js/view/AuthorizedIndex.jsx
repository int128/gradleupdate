import React from "react";
import {Link} from "react-router";
import OAuthSession from "../repository/OAuthSession.jsx";
import GitHub from "../repository/GitHub.jsx";
import GradleUpdate from "../repository/GradleUpdate.jsx";
import ErrorHeader from "./ErrorHeader.jsx";
import Footer from "./Footer.jsx";

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.github = new GitHub(OAuthSession.getToken());
    this.gradleUpdate = new GradleUpdate(OAuthSession.getToken());
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
    this.gradleUpdate.updateRepository(fullName);
  }

  render() {
    return (
      <div className="container">
        <ErrorHeader kind="API Error" message={this.state.error}/>

        <div className="page-header text-center">
          <h1>Gradle Update</h1>
          <GULatestGradle
            latestGradle={this.state.latestGradle}/>
        </div>

        <section>
          <div className="row">
            <div className="col-lg-3 col-md-3 col-sm-3">
              <div className="text-center">
                <GHAuthenticatedUser
                  user={this.state.user}/>
              </div>
            </div>
            <div className="col-lg-9 col-md-9 col-sm-9">
              <GHRepositories
                repos={this.state.repos}
                updateGradleWrapper={this.updateGradleWrapper.bind(this)}
              />
            </div>
          </div>
        </section>

        <Footer/>
      </div>
    );
  }
}

class GULatestGradle extends React.Component {
  render() {
    return this.props.latestGradle ?
      (
        <div>
          Latest version is <strong>Gradle {this.props.latestGradle.version}</strong>
        </div>
      ) : null;
  }
}

class GHAuthenticatedUser extends React.Component {
  signOut(e) {
    this.props.signOut();
    e.preventDefault();
  }

  render() {
    return (
      <div>
        {this.props.user ? (
          <div>
            <img src={this.props.user.avatar_url} className="img-circle" width="128" height="128"/>
            <h2>{this.props.user.name}</h2>
            <p>@{this.props.user.login}</p>
          </div>
        ) : null}
        <div>
          <Link to="/signout" className="btn btn-default">
            Sign Out
          </Link>
        </div>
      </div>
    );
  }
}

class GHRepositories extends React.Component {
  render() {
    return (
      <ul className="list-group">
        {this.props.repos ? this.props.repos.map((repo) => (
          <li className="list-group-item">
            <div className="pull-right">
              <button
                type="button"
                className="btn btn-default btn-xs"
                onClick={this.props.updateGradleWrapper.bind(this, repo.full_name)}>
                Update
              </button>
            </div>
            <a href={`/${repo.owner.login}/${repo.name}/status`}>
              {repo.owner.login}/<strong>{repo.name}</strong>
            </a>
            <div className="clearfix"></div>
          </li>
        )) : null}
      </ul>
    );
  }
}
