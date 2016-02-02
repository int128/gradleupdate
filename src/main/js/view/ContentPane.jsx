import React from 'react';

import Footer from './Footer.jsx';

export default class extends React.Component {
  render() {
    if (this.props.repo) {
      return (
        <div>
          <h2>{this.props.repo.full_name}</h2>
          <p><img src={`/${this.props.repo.full_name}/status.svg`}/></p>
          <p>{this.props.repo.description}</p>
          <Footer/>
        </div>
      );
    } else {
      return (
        <div>
          <h2>Gradle Update</h2>
          <Footer/>
        </div>
      );
    }
  }
}
