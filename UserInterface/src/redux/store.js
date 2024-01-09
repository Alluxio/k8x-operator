import { createStore, combineReducers } from '@reduxjs/toolkit';
import datasetReducer from './reducers/datasetReducer';
import alertReducer from './reducers/alertReducer';
import alluxioReducer from "./reducers/alluxioReducer";


const rootReducer = combineReducers({
    dataset: datasetReducer,
    alert: alertReducer,
    alluxio: alluxioReducer,
});

const store = createStore(rootReducer);

export default store;
