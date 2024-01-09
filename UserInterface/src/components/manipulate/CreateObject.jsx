import React, { useState } from 'react';
import { Grid, Box, Button, Modal } from "@mui/material";
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import { modalStyle, DataType } from "../../util/util";
import CodeEditor from '@uiw/react-textarea-code-editor';


const CreateObject = ({   objectKind,
                          handleSendRequest
}) => {
    const [objectConfigInput, setObjectConfigInput] = useState('');
    const [modalOpen, setModalOpen] = useState(false);
    const httpMethod = 'POST';

    const handleSubmit = (e) => {
        e.preventDefault();
        handleSendRequest(httpMethod, DataType.YAMLString, objectConfigInput);
    };

    const handleModalOpen = () => {
        setModalOpen(true);
    };

    const handleModalClose = () => {
        setModalOpen(false);
    };

    const handleFileUpload = (event) => {
        const file = event.target.files[0];
        if (!file) {
            return;
        }

        const reader = new FileReader();
        reader.onload = (e) => {
            setObjectConfigInput(e.target.result)
        };
        reader.readAsText(file);
    };

    const modalElement = (
        <div className="CreateObject">
            <Grid container spacing={2}>
                <Grid item lg={12} xs={12}>
                    <Box sx={{
                        height: '30px',
                        fontSize: '24px',
                        fontWeight: 'bold',
                        width: '100%',}}>
                        Create a new {objectKind}
                    </Box>
                </Grid>
                <Grid item lg={7} xs={12}>
                    <Box sx={{
                        height: '750px',
                        width: '100%',
                        overflow: 'scroll'
                    }}>
                        <CodeEditor
                            value={objectConfigInput}
                            language="yaml"
                            placeholder="Enter or Upload YAML Style Configuration (.yaml/.yml)."
                            onChange={(evn) => setObjectConfigInput(evn.target.value)}
                            padding={15}
                            style={{
                                fontSize: '16px',
                                fontFamily: 'ui-monospace,SFMono-Regular,SF Mono,Consolas,Liberation Mono,Menlo,monospace',
                                minHeight: '100%',
                            }}
                        />
                    </Box>
                </Grid>

                <Grid item lg={5} xs={12}>
                    <Grid container spacing={2}>
                        <Grid item lg={12} xs={6}>
                            <Box sx={{height: '50px'}}>
                                <Button variant="contained" size="large"
                                        component="label"
                                        startIcon={<CloudUploadIcon />}>
                                    Upload {objectKind} Configuration File

                                    <input type="file" hidden
                                           onChange={handleFileUpload} accept=".yaml,.yml"/>
                                </Button>
                            </Box>
                        </Grid>

                        <Grid item lg={6} xs={2}>
                            <Button variant="outlined" size="large" onClick={handleSubmit}>
                                Create
                            </Button>
                        </Grid>
                        <Grid item lg={6} xs={2}>
                            <Button variant="outlined" size="large" onClick={handleModalClose}>
                                Cancel
                            </Button>
                        </Grid>
                    </Grid>
                </Grid>

            </Grid>
        </div>
    )

    return (
        <div className='Create-Object'>
            <Button variant="outlined" onClick={handleModalOpen} sx={{fontSize: '18px'}}>
                Create a new {objectKind}
            </Button>
            <Modal
                open={modalOpen}
                onClose={handleModalClose}
                aria-labelledby="modal-modal-title"
                aria-describedby="modal-modal-description"
            >
                <Box sx={modalStyle}>
                    {modalElement}
                </Box>
            </Modal>
        </div>
    )
}

export default CreateObject
