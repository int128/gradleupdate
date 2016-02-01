import React from 'react';

import GitHub from '../repository/GitHub.jsx';

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.state = {projects: null};
  }
  componentDidMount() {
    const github = new GitHub(this.props.token);
    github.findRepositories({sort: 'updated'})
    .then((xhr, repos) => this.setState({repos: repos}));
  }
  render() {
    if (this.state.repos) {
      return (
        <div>
          {this.state.repos.map((repo) => (
            <div>
              <h3>{repo.full_name}</h3>
              <p>{repo.description}</p>
            </div>
          ))}
        </div>
      );
    } else {
      return null;
    }
  }
}
