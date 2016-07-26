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
    const params = this.props.params;
    this.github.findRepository(`${params.user}/${params.repo}`)
      .then((xhr, repository) => this.setState({ghRepository: repository}))
      .catch((e) => this.setState({error: e}));
    this.gradleUpdate.findRepository(`${params.user}/${params.repo}`)
      .then((xhr, repository) => this.setState({guRepository: repository}))
      .catch((e) => this.setState({error: e}));
  }

  render() {
    return (
      <div className="container">
        <ErrorHeader kind="API Error" message={this.state.error}/>

        <section className="well">
          Not Implemented Yet
        </section>

        <Footer/>
      </div>
    );
  }
}
