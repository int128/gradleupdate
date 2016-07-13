import React from 'react';

import Footer from './Footer.jsx';

export default class extends React.Component {
  onClick() {
    this.props.onAuthorize();
  }
  render() {
    return (
      <div>
        {this.props.error ? (<Error message={this.props.error}/>) : null}

        <div className="jumbotron">
          <div className="container text-center">
            <h1>Gradle Update</h1>
            <p>keeps the latest Gradle wrapper on your GitHub repositories</p>
          </div>
        </div>

        <div className="container text-center">
          <button className="btn btn-default" onClick={this.onClick.bind(this)}>
            <span className="glyphicon glyphicon-off"></span>
            &nbsp;
            Sign in with GitHub Account
          </button>
        </div>

        <Footer/>
      </div>
    );
  }
}

class Error extends React.Component {
  render() {
    return (
      <div className="alert alert-warning">
        <strong>Authorization Error:</strong> {this.props.message}
      </div>
    );
  }
}
