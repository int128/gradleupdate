import React from 'react';

export default class extends React.Component {
  onClick(e) {
    this.props.onSignOut();
    e.preventDefault();
  }
  render() {
    return (
      <nav className="navbar navbar-default">
        <div className="container">
          <ul className="nav navbar-nav">
            <li><a className="navbar-brand" href="/">Gradle Update</a></li>
          </ul>
          <ul className="nav navbar-nav navbar-right">
            <li>
              <a href="#" onClick={this.onClick.bind(this)}>
                <span className="glyphicon glyphicon-log-out"></span>
                &nbsp;
                Sign out
              </a>
            </li>
          </ul>
        </div>
      </nav>
    );
  }
}
