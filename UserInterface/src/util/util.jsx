import {AlertMethod} from "../components/alert/Alert";
import yaml from 'js-yaml'

export const DataType =  {
    JSONString: 1,
    JSONObject: 2,
    YAMLString: 3,
    YAMLObject: 4,
}

// modalStyle contains styling for all modals.
export const modalStyle = {
    position: 'absolute',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: '75%',
    height: '80%',
    bgcolor: 'background.paper',
    border: '1px solid #000',
    borderRadius: ' 10px',
    p: 4,
};


export const FetchHeaders = {'Content-Type':'application/json'}

/**
 * Use HTTP method in requestOptions to communicate with api server.
 *
 * @param requestOptions contains http method. headers and body.
 * @param {string} serverUrl is a relative or points to your development server.
 * @function sendRequest
 * @returns [AlertMethod, AlertMessage]
 */
export async function sendRequest(requestOptions, serverUrl) {
    try {
        const response = await fetch(serverUrl, requestOptions);
        if (!response.ok) {
            try {
                const errorBody = await response.json();
                const errorMessage = `${errorBody.title}: ${errorBody.details}`;
                return [AlertMethod.Error, errorMessage];
            } catch (jsonError) {
                return [AlertMethod.Error, 'An error occurred, but the error message from API-Server could not be parsed.'];
            }
        }
        return [AlertMethod.Success, 'Success!'];
    } catch (error) {
        return [AlertMethod.Error, error.toString()];
    }
}


/**
 * Use GET method to get data
 *
 * @param {string} serverUrl is a relative or points to your development server.
 * @function getRequest
 * @returns [AlertMethod, AlertMessage, Result in JSON Object(Optional)]
 */
export async function getRequest(serverUrl){
    try {
        const response = await fetch(serverUrl);

        if (!response.ok) {
            try {
                const errorBody = await response.json();
                const errorMessage = `${errorBody.title}: ${errorBody.details}`;
                return [AlertMethod.Error, errorMessage];
            } catch (jsonError) {
                return [AlertMethod.Error, 'An error occurred, but the error message from API-Server could not be parsed.'];
            }
        }

        const result = await response.json();
        return [AlertMethod.Success, 'Success!', result];
    } catch (error) {
        console.log(error);
        return [AlertMethod.Error, error.toString()];
    }
}

/**
 * Convert YAML String to JSON String
 *
 * @param {string} yamlStr is a YAML string
 * @function YAMLtoJSON
 * @returns {string} jsonStr
 */
export function YAMLtoJSON(yamlStr) {
    const obj = yaml.load(yamlStr);
    return JSON.stringify(obj);
}

/**
 * Convert JSON String to YAML String
 *
 * @param {JSON} jsonObj is a YAML string
 * @function JSONtoStringifyYAML
 * @returns {string} yamlStr
 */
export function JSONtoStringifyYAML(jsonObj) {
    return yaml.dump(jsonObj)
}


/**
 * Convert Input to StringifyJSON
 *
 * @param dataType is defined in DataType
 * @param inputData can be type that specified in dataType
 * @function convertToStringifyJSON
 * @returns {string} jsonStr
 */
export function convertToStringifyJSON(dataType, inputData) {
    if (dataType === DataType.YAMLString) {
        const parsedYaml = yaml.load(inputData);
        return JSON.stringify(parsedYaml, null, 2);
    }

    if (dataType === DataType.JSONString) {
        const parsedJSON = JSON.parse(inputData);
        return JSON.stringify(parsedJSON);
    }

    if (dataType === DataType.YAMLObject) {
        return JSON.stringify(inputData)
    }
    if (dataType === DataType.JSONObject) {
        return JSON.stringify(inputData)
    }

    return JSON.stringify(JSON)
}


/**
 * Generate HTTP Request Body for Fetch()
 *
 * @param httpMethod
 * @param dataType
 * @param inputData
 * @function generateHttpRequestOptions
 * @returns  requestOptions
 */
export function generateHttpRequestOptions(httpMethod, dataType, inputData) {
    return {
        method: httpMethod,
        headers: FetchHeaders,
        body: convertToStringifyJSON(dataType, inputData),
    }
}

