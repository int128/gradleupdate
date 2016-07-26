import React from "react";
import OAuthSession from "../repository/OAuthSession.jsx";
import AuthorizedIndex from "./AuthorizedIndex.jsx";
import NotAuthorizedIndex from "./NotAuthorizedIndex.jsx";

export default class extends React.Component {
  render() {
    if (OAuthSession.getToken()) {
      return (<AuthorizedIndex/>);
    } else {
      return (<NotAuthorizedIndex/>);
    }
  }
}
