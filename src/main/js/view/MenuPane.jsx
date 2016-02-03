import React from 'react';

export default class extends React.Component {
  onClickRepo(repo, e) {
    this.props.onSelectRepo(repo);
    e.preventDefault();
  }
  render() {
    return (
      <div>
        <nav className="navbar navbar-default">
          <div className="container-fluid">
            <div className="navbar-header">
              <a className="navbar-brand" href="/">Gradle Update</a>
            </div>
          </div>
        </nav>

        {this.props.repos ? (
          <div className="list-group">
            {this.props.user ? (
              <a href={this.props.user.html_url} className="list-group-item">
                <img src={this.props.user.avatar_url} width="32" height="32"/>
                {this.props.user.login}
              </a>
            ) : null}

            {this.props.repos.map((repo) => (
              <a href={`/${repo.full_name}`} className="list-group-item"
                onClick={this.onClickRepo.bind(this, repo)}>
                {repo.full_name}
              </a>
            ))}
          </div>
        ) : null}
      </div>
    );
  }
}
