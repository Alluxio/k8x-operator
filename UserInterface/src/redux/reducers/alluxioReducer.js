// reducers/datasetReducer.js
const initialState = {
    alluxioList: []
};

const alluxioReducer = (state = initialState, action) => {
    switch (action.type) {
        case 'UPDATE_ALLUXIO_LIST':
            return {
                ...state,
                alluxioList: action.payload
            };
        default:
            return state;
    }
};

export default alluxioReducer;
