import React from 'react';

import Footer from './Footer.jsx';

export default class extends React.Component {
  render() {
    return (
      <div className="container">
        <section className="text-center">
          {this.props.user ? (
            <div>
              <img src={this.props.user.avatar_url} className="img-circle" width="128" height="128"/>
              <h2>{this.props.user.name}</h2>
              <p>@{this.props.user.login}</p>
            </div>
          ) : null}
        </section>

        <section className="text-center">
          {this.props.repos ? this.props.repos.map((repo) => (
            <div>
              <h3><small>{repo.owner.login}/</small>{repo.name}</h3>
              <p>{repo.description}</p>
              <img src={`/${repo.full_name}/status.svg`}/>
            </div>
          )) : null}
        </section>

        <Footer/>
      </div>
    );
  }
}
