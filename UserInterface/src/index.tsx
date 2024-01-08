import * as React from "react";
import * as ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App.";
import store from './redux/store'
import { Provider } from 'react-redux'

// @ts-ignore
ReactDOM.createRoot(document.getElementById("root")).render(
    <Provider store={store}>
        <React.StrictMode>
            <App />
        </React.StrictMode>
    </Provider>
);