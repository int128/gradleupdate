import React from "react";

export default class extends React.Component {
  render() {
    return this.props.latestGradle ?
      (
        <div>
          Latest version is <strong>Gradle {this.props.latestGradle.version}</strong>
        </div>
      ) : null;
  }
}
