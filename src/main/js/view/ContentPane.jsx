import React from "react";

export default class extends React.Component {
  render() {
    return (
      <section>
        <table className="table">
          <tbody>
            {this.props.repos ? this.props.repos.map((repo) => (
              <tr>
                <td>
                  <img src={`/${repo.full_name}/status.svg`}/>
                </td>
                <td>
                  {repo.owner.login}/{repo.name}
                  <br/>
                  {repo.description}
                  <br/>
                  <button className="btn btn-info"
                    onClick={this.props.updateGradleWrapper.bind(this, repo.full_name)}>
                    Update now
                  </button>
                </td>
              </tr>
            )) : null}
          </tbody>
        </table>
      </section>
    );
  }
}
