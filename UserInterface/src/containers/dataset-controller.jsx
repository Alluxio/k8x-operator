import React, { Component} from 'react';
import CreateObject from "../components/manipulate/CreateObject";
import Datasets from "../components/dataset/Datasets";
import {AlertMethod} from "../components/alert/Alert";
import {
    sendRequest,
    getRequest,
    generateHttpRequestOptions
} from "../util/util";
import { connect } from 'react-redux';
import {updateDatasetList} from "../redux/actions/datasetActions";
import {setAlert} from "../redux/actions/alertActions";
import Box from "@mui/material/Box";

const mapStateToProps = state => ({
    datasetList: state.dataset.datasetList,
    alertMethod: state.alert.alertMethod,
    alertMessage: state.alert.alertMessage,
});

const mapDispatchToProps = dispatch => ({
    updateDatasetList: (newList) => dispatch(updateDatasetList(newList)),
    setAlert: (alertMethod, alertMessage) => dispatch(setAlert(alertMethod, alertMessage)),
});

const serverUrl = '/api/dataset'

class DatasetController extends Component {
    componentDidMount() {
        this.interval = setInterval(() => {
            this.handleGetRequest()
        }, 2000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    handleGetRequest = () =>  {
        // Call getRequest
        getRequest(serverUrl).then(statusCodeAndMsg => {
            if (statusCodeAndMsg[0] === AlertMethod.Success) {
                const numberOfDataset = statusCodeAndMsg[2]['datasets'].length
                const newDatasetList = [];
                for (let i = 0; i < numberOfDataset; i++) {
                    newDatasetList.push(statusCodeAndMsg[2]['datasets'][i])
                }
                // Update the dataset list in Redux store
                this.props.updateDatasetList(newDatasetList);
            } else {
                this.props.setAlert(statusCodeAndMsg[0], statusCodeAndMsg[1])
            }
        })
    }

    // handleSendRequest is a func that will call  and update the panel status
    handleSendRequest = (httpMethod, dataType, inputData) =>  {
        // Prep HTTP Request Options
        let requestOptions = {}
        try {
            requestOptions = generateHttpRequestOptions(httpMethod, dataType, inputData)
        } catch (error){
            this.props.setAlert(AlertMethod.Warning, 'Unable to Parse Input. ' + error)
            return
        }

        // Send Request
        sendRequest(requestOptions, serverUrl).then(statusCodeAndMsg => {
            this.props.setAlert(statusCodeAndMsg[0], statusCodeAndMsg[1])
        })
    }

    render() {
        return (
            <div className="Dataset">
                <Box sx={{fontSize: '18px'}}>
                    <h1>Dataset Controller Panel</h1>
                </Box>
                <CreateObject
                    objectKind={'Dataset'}
                    handleSendRequest={this.handleSendRequest}
                />
                <br/>
                <Datasets
                    handleSendRequest={this.handleSendRequest}
                />
            </div>
        );
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(DatasetController);
