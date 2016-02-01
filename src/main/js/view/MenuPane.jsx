import React from 'react';

export default class extends React.Component {
  render() {
    return (
      <div>
        <button className="btn btn-default" onClick={this.props.onSignOut.bind(this)}>
          Sign Out
        </button>

        {this.props.user ? (
          <div>
            <a href={this.props.user.html_url}>
              <img src={this.props.user.avatar_url} width="32" height="32"/>
              {this.props.user.name}
            </a>
          </div>
        ) : null}

        {this.props.repos ? (
          <ul>
            {this.props.repos.map((repo) => (<li>{repo.full_name}</li>))}
          </ul>
        ) : null}
      </div>
    );
  }
}
