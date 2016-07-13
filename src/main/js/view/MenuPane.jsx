import React from "react";

export default class extends React.Component {
  onClick(e) {
    this.props.onSignOut();
    e.preventDefault();
  }
  render() {
    return (
      <section className="text-center">
        {this.props.user ? (
          <div>
            <img src={this.props.user.avatar_url} className="img-circle" width="128" height="128"/>
            <h2>{this.props.user.name}</h2>
            <p>@{this.props.user.login}</p>
          </div>
        ) : null}
      </section>
    );
  }
}
