import React from "react";
import {Link} from "react-router";
import Footer from "./Footer.jsx";

export default class extends React.Component {
  render() {
    return (
      <div className="container">
        <section className="text-center">
          <h1>404 Not Found</h1>
          <Link to="/" className="btn btn-default">Index</Link>
        </section>

        <Footer/>
      </div>
    );
  }
}
 