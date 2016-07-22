import React from "react";

export default class extends React.Component {
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
