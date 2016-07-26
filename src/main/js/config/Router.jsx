import React from "react";
import {Router, Route, browserHistory} from "react-router";
import {SignIn, SignInExchange, SignInError, SignOut} from "../view/OAuth.jsx";
import Index from "../view/Index.jsx";
import NoMatch from "../view/NoMatch.jsx";
import RepositoryStatus from "../view/RepositoryStatus.jsx";

export default class extends React.Component {
  render() {
    return (
      <Router history={browserHistory}>
        <Route path="/signin" component={SignIn} />
        <Route path="/signin/exchange" component={SignInExchange} />
        <Route path="/signin/error" component={SignInError} />
        <Route path="/signout" component={SignOut} />

        <Route path="/:user/:repo/status" component={RepositoryStatus} />
        <Route path="/" component={Index} />

        <Route path="*" component={NoMatch}/>
      </Router>
    );
  }
}
