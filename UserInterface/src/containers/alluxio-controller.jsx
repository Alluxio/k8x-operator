import React, { Component } from 'react';
import {AlertMethod} from "../components/alert/Alert";
import {
    sendRequest,
    getRequest,
    generateHttpRequestOptions
} from "../util/util";
import AlluxioClusters from "../components/alluxio_cluster/AlluxioClusters";
import CreateObject from "../util/CreateObject";

import { connect } from 'react-redux';
import {setAlert} from "../redux/actions/alertActions";
import {updateAlluxioList} from "../redux/actions/alluxioActions";
import Box from "@mui/material/Box";

const mapStateToProps = state => ({
    alertMethod: state.alert.alertMethod,
    alertMessage: state.alert.alertMessage,
    alluxioList: state.alluxio.alluxioList
});

const mapDispatchToProps = dispatch => ({
    setAlert: (alertMethod, alertMessage) => dispatch(setAlert(alertMethod, alertMessage)),
    updateAlluxioList: (newList) => dispatch(updateAlluxioList(newList)),
});

const serverUrl = '/api/alluxio_cluster'

class AlluxioController extends Component {
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
                const numberOfAlluxio = statusCodeAndMsg[2]['alluxio-clusters'].length
                const newAlluxioList = [];
                for (let i = 0; i < numberOfAlluxio; i++) {
                    newAlluxioList.push(statusCodeAndMsg[2]['alluxio-clusters'][i])
                }
                // Update the dataset list in Redux store
                this.props.updateAlluxioList(newAlluxioList);
            } else {
                this.props.setAlert(statusCodeAndMsg[0], statusCodeAndMsg[1])
            }
        })
    }

    // handleSendRequest is a func that will call sendRequest() and update the panel status
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
            <div className="Alluxio">
                <Box sx={{fontSize: '18px'}}>
                    <h1>Alluxio Controller Panel</h1>
                </Box>
                <CreateObject
                    objectKind={'Alluxio Cluster'}
                    handleSendRequest={this.handleSendRequest}
                />
                <br/>
                <AlluxioClusters
                    handleSendRequest={this.handleSendRequest}
                />
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(AlluxioController);
