import React from "react";
import Snackbar from '@mui/material/Snackbar';
import MuiAlert from '@mui/material/Alert';
import { useSelector, useDispatch } from 'react-redux';
import {clearAlert} from "../../redux/actions/alertActions";

const Alert = React.forwardRef(function Alert(props, ref) {
    return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

export const AlertMethod =  {
    NoStatus : -1,
    Info : 0,
    Success : 1,
    Warning : 2,
    Error : 3,
}

const Severity = [
    "info",
    "success",
    "warning",
    "error"
]

const OperatorAlert = () => {
    const dispatch = useDispatch();
    const alertMethod = useSelector(state => state.alert.alertMethod);
    const alertMessage = useSelector(state => state.alert.alertMessage);

    const handleClose = (event, reason) => {
        if (reason === 'clickaway') {return;}
        dispatch(clearAlert());
    };

    return (
        <Snackbar open={alertMethod !== AlertMethod.NoStatus}
                  autoHideDuration={3500}
                  anchorOrigin={{vertical: 'top', horizontal: 'right'}}
                  onClose={handleClose}>
            <Alert onClose={handleClose}
                   severity={Severity[alertMethod]}
                   sx={{width: '100%', fontSize: '18px'}}>

                {alertMessage}
            </Alert>
        </Snackbar>
    )
}
export default OperatorAlert
