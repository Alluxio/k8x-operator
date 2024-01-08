export const setAlert = (alertMethod, alertMessage) => ({
    type: 'SET_ALERT',
    payload: { alertMethod, alertMessage }
});

export const clearAlert = () => ({
    type: 'CLEAR_ALERT'
});
