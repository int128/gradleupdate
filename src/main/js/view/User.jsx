import React from 'react';

import GitHub from '../repository/GitHub.jsx';

export default class extends React.Component {
  constructor(props) {
    super(props);
    this.state = {user: null};
  }
  componentDidMount() {
    const github = new GitHub(this.props.token);
    github.getUser().then((xhr, user) => this.setState({user: user}));
  }
  render() {
    if (this.state.user) {
      return (
        <div>
          <a href={this.state.user.html_url}>
            {this.state.user.name}
            <img src={this.state.user.avatar_url}/>
          </a>
        </div>
      );
    } else {
      return null;
    }
  }
}
