import React from 'react';

import User from './User.jsx';
import Projects from './Projects.jsx';
import Footer from './Footer.jsx';

export default class extends React.Component {
  onClick() {
    this.props.onUnauthorize();
  }
  render() {
    return (
      <div className="container">
        <button className="btn btn-default" onClick={this.onClick.bind(this)}>
          Sign Out
        </button>
        <h2>User</h2>
        <User token={this.props.token}/>
        <Projects token={this.props.token}/>
        <Footer/>
      </div>
    );
  }
}
