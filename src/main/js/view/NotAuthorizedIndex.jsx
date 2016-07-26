import React from "react";
import {Link} from "react-router";
import GradleUpdate from "../repository/GradleUpdate.jsx";
import ErrorHeader from "./ErrorHeader.jsx";
import Footer from "./Footer.jsx";

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.gradleUpdate = new GradleUpdate();
  }

  componentDidMount() {
    this.gradleUpdate.findPullRequests()
      .then((xhr, pullRequests) => this.setState({pullRequests: pullRequests}))
      .catch((e) => this.setState({error: e}));
  }

  render() {
    return (
      <div className="container">
        <ErrorHeader kind="API Error" message={this.state.error}/>

        <section className="text-center">
          <div className="page-header">
            <h1>Gradle Update</h1>
            <p>keeps the latest Gradle wrapper on your GitHub repositories</p>
          </div>

          <Link to="/signin" className="btn btn-primary">
            Sign in with GitHub Account
          </Link>
        </section>

        <section className="text-center">
          <h2>Contributions</h2>
          <GUPullRequests pullRequests={this.state.pullRequests}/>
        </section>

        <Footer/>
      </div>
    );
  }
}

class GUPullRequests extends React.Component {
  render() {
    return (
      <ul className="list-group">
        {this.props.pullRequests ? this.props.pullRequests.map((pullRequest) => (
          <li className="list-group-item">
            <div className="pull-right">
              <small>{pullRequest.createdAt}</small>
            </div>
            Sent <a href={pullRequest.url}>Pull Request on {pullRequest.fullName}</a> to
            update from {pullRequest.fromVersion} to {pullRequest.toVersion}
            <div class="clearfix"></div>
          </li>
        )) : null}
      </ul>
    );
  }
}
