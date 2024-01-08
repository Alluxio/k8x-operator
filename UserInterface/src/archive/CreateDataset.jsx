import React, { Component } from 'react';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Modal from '@mui/material/Modal';
import {YAMLtoJSON} from "../util/util";

const style = {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: '75%',
    height: '80%',
    bgcolor: 'background.paper',
    border: '1px solid #000',
    p: 4,
};

class CreateDataset extends Component {
    constructor(props) {
        super(props)
        this.state = {
            httpMethod: 'POST',
            datasetConfigUserInput: '',

            createDatasetModalOpen: false,
        };
    }

    handleChange = (e) => {
        e.preventDefault()
        this.setState(
            {datasetConfigUserInput : e.currentTarget.value}
        )
    }

    handleSubmit = (e) => {
        e.preventDefault();
        // Convert to JSON
        this.props.handleSendRequest(YAMLtoJSON(this.state.datasetConfigUserInput), this.state.httpMethod)
    }

    handleModalOpen = () => {
        this.setState({
            createDatasetModalOpen: true,
        })
    }
    handleModalClose = () => {
        this.setState({
            createDatasetModalOpen: false,
        })
    }


    render() {

        const modalElement = (
            <div className="CreateDataset">
                <h2>Create a new Dataset</h2>
                Put YAML Config Below:
                <form onSubmit={this.handleSubmit}>
                    <label>
                        <textarea
                            rows='30' cols='50' onChange={this.handleChange}
                            spellCheck='false' defaultValue={this.state.datasetConfigUserInput}
                        />
                    </label>
                    <input type="submit" value="Submit DatasetController Config JSON"/>
                </form>
            </div>
        )

        return (
            <div className='Create-Dataset'>
                <Button variant="outlined" onClick={this.handleModalOpen}>
                    Create a New Dataset
                </Button>
                <Modal
                    open={this.state.createDatasetModalOpen}
                    onClose={this.handleModalClose}
                    aria-labelledby="modal-modal-title"
                    aria-describedby="modal-modal-description"
                >
                    <Box sx={style}>
                        {modalElement}
                    </Box>
                </Modal>
            </div>
        )
    }
}

export default CreateDataset
