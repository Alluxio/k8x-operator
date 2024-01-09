import { Component } from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';

class DeleteDataset extends Component {
    constructor(props) {
        super(props)
        this.state = {
            httpMethod: 'DELETE',
            open: false
        };
    }

    handleSubmit = () => {
        this.setState({ open: false });
        this.props.handleSendRequest(JSON.stringify(this.props.datasetConfigJSON), this.state.httpMethod)
    };

    handleClickOpen = () => {
        this.setState({ open: true });
    };

    handleClickClose = () => {
        this.setState({ open: false });
    };

    render() {
        return (
            <div>
                <Button variant="outlined" onClick={this.handleClickOpen}>
                    Delete
                </Button>

                <Dialog
                    open={this.state.open}
                    onClose={this.handleClickClose}
                    aria-labelledby="alert-dialog-title"
                    aria-describedby="alert-dialog-description"
                >
                    <DialogTitle id="alert-dialog-title">
                        Do you want to Delete {this.props.datasetConfigJSON.name}?
                    </DialogTitle>
                    <DialogContent>
                        <DialogContentText id="alert-dialog-description">
                            This will Delete {this.props.datasetConfigJSON.name} from Alluxio Cluster.
                        </DialogContentText>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={this.handleClickClose}>No</Button>
                        <Button onClick={this.handleSubmit} autoFocus>Yes</Button>
                    </DialogActions>
                </Dialog>
            </div>
        )
    }
}


export default DeleteDataset