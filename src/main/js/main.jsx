import React from "react";
import {render} from "react-dom";
import Router from "./config/Router.jsx";
import "./main.less";

render(<Router/>, document.getElementById('app'));
