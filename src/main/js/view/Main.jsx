import React from 'react';

export default class extends React.Component {
  render() {
    return (
      <div>

        <div className="jumbotron">
          <div className="container text-center">
            <h1>Gradle Update</h1>
            <p>keeps the latest Gradle wrapper in your GitHub repositories</p>
          </div>
        </div>

        <div className="container text-center">
          <a href="https://github.com/int128/gradleupdate" className="btn btn-default">
            <i className="fa fa-lg fa-github"></i>
            &nbsp;
            Enable Auto Update on GitHub
          </a>
        </div>

        <div className="footer container text-center">
          &copy; Hidetake Iwata, 2015.
        </div>

      </div>
    );
  }
}
