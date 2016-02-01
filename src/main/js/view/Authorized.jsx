import React from 'react';

import GitHub from '../repository/GitHub.jsx';

import MenuPane from './MenuPane.jsx';
import ContentPane from './ContentPane.jsx';

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.state = {user: null, repos: null};
  }
  componentDidMount() {
    const github = new GitHub(this.props.token);
    github.getUser()
      .then((xhr, user) => this.setState({user: user}));
    github.findRepositories({sort: 'updated'})
      .then((xhr, repos) => this.setState({repos: repos}));
  }
  render() {
    return (
      <div className="row">
        <div className="col-lg-3 col-md-3 col-sm-3 gu-menu-pane">
          <MenuPane onSignOut={this.props.onUnauthorize.bind(this)}
            user={this.state.user}
            repos={this.state.repos}/>
        </div>
        <div className="col-lg-9 col-md-9 col-sm-9 gu-content-pane">
          <ContentPane/>
        </div>
      </div>
    );
  }
}
