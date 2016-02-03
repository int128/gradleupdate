import React from 'react';

import GitHub from '../repository/GitHub.jsx';

import MenuPane from './MenuPane.jsx';
import ContentPane from './ContentPane.jsx';

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};
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
      <div className="container-fluid">
        <div className="row">
          <MenuPane
            onSignOut={this.props.onUnauthorize.bind(this)}/>
          <ContentPane
            user={this.state.user}
            repos={this.state.repos}/>
        </div>
      </div>
    );
  }
}
