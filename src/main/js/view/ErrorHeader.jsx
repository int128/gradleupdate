import React from "react";

export default class extends React.Component {
  render() {
    if (this.props.message && this.props.kind) {
      return (
        <div className="alert alert-warning" role="alert">
          <strong>{this.props.kind}:</strong> {this.props.message}
        </div>
      );
    } else {
      return null;
    }
  }
}
