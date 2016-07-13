import React from "react";
import GitHub from "../repository/GitHub.jsx";
import GradleUpdate from "../repository/GradleUpdate.jsx";
import MenuPane from "./MenuPane.jsx";
import ContentPane from "./ContentPane.jsx";
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
      .then((xhr, user) => this.setState({user: user}));
    this.github.findRepositories({sort: 'updated'})
      .then((xhr, repos) => this.setState({repos: repos}));
  }
  updateGradleWrapper(fullName) {
    this.gradleUpdate.update(fullName, '2.14');  //FIXME
  }
  render() {
    return (
      <div className="container">
        <div className="row">
          <div className="col-lg-3 col-md-3 col-sm-3">
            <MenuPane
              user={this.state.user}
              onSignOut={this.props.onUnauthorize.bind(this)}/>
          </div>
          <div className="col-lg-9 col-md-9 col-sm-9">
            <ContentPane
              repos={this.state.repos}
              updateGradleWrapper={this.updateGradleWrapper.bind(this)}
            />
          </div>
        </div>
        <Footer/>
      </div>
    );
  }
}
