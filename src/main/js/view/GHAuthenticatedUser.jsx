import React from "react";

export default class extends React.Component {
  signOut(e) {
    this.props.signOut();
    e.preventDefault();
  }

  render() {
    return (
      <div>
        {this.props.user ? (
          <div>
            <img src={this.props.user.avatar_url} className="img-circle" width="128" height="128"/>
            <h2>{this.props.user.name}</h2>
            <p>@{this.props.user.login}</p>
          </div>
        ) : null}
        <div>
          <button className="btn btn-default" onClick={this.signOut.bind(this)}>
            Sign Out
          </button>
        </div>
      </div>
    );
  }
}
