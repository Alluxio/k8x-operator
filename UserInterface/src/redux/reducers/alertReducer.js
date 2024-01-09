import {AlertMethod} from "../../components/alert/Alert";

const initialState = {
    alertMethod: AlertMethod.NoStatus,
    alertMessage: '',
};

const alertReducer = (state = initialState, action) => {
    switch (action.type) {
        case 'SET_ALERT':
            return {
                ...state,
                alertMethod: action.payload.alertMethod,
                alertMessage: action.payload.alertMessage,
            };
        case 'CLEAR_ALERT':
            return {
                ...state,
                alertMethod: AlertMethod.NoStatus,
                alertMessage: '',
            };
        default:
            return state;
    }
};

export default alertReducer;