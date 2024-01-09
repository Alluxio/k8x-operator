import React, { useState } from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import {DataType} from "../../util/util";

const DeleteObject = ({   objectKind,
                          objectName,
                          handleSendRequest
}) => {
    const [modalOpen, setModalOpen] = useState(false);
    const httpMethod = 'DELETE';

    const handleSubmit = () => {
        setModalOpen(false);
        handleSendRequest(httpMethod, DataType.JSONObject, { name: objectName });
    };

    const handleClickOpen = () => {
        setModalOpen(true);
    };

    const handleClickClose = () => {
        setModalOpen(false);
    };
        return (
            <div>
                <Button variant="outlined" onClick={handleClickOpen} size="large">
                    Delete
                </Button>

                <Dialog
                    open={modalOpen}
                    onClose={handleClickClose}
                    aria-labelledby="alert-dialog-title"
                    aria-describedby="alert-dialog-description"
                >
                    <DialogTitle id="alert-dialog-title">
                        Do you want to Delete {objectName} [{objectKind}]?
                    </DialogTitle>
                    <DialogContent>
                        <DialogContentText id="alert-dialog-description">
                            This will Delete {objectName} from Kubernetes Cluster.
                        </DialogContentText>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={handleClickClose}>No</Button>
                        <Button onClick={handleSubmit} autoFocus>Yes</Button>
                    </DialogActions>
                </Dialog>
            </div>
        )
}

export default DeleteObject
