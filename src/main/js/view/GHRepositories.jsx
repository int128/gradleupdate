import React from "react";

export default class extends React.Component {
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
