// reducers/datasetReducer.js
const initialState = {
    datasetList: [],
};

const datasetReducer = (state = initialState, action) => {
    switch (action.type) {
        case 'UPDATE_DATASET_LIST':
            return {
                ...state,
                datasetList: action.payload
            };
        default:
            return state;
    }
};

export default datasetReducer;
