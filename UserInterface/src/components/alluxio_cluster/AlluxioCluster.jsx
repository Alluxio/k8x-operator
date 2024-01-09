import React from 'react';
import ViewObject from "../manipulate/ViewObject";
import TableCell from "@mui/material/TableCell";
import TableRow from "@mui/material/TableRow";
import {JSONtoStringifyYAML} from "../../util/util";


const AlluxioCluster = ({name, originalJSON, handleSendRequest}) => {
    const specJSON = originalJSON['alluxio-cluster-config']['spec'];
    const statusJSON = originalJSON['status'];

    return (
        <TableRow
            key={name}
            sx={{'&:last-child td, &:last-child th': {border: 0}}}
        >
            <TableCell>{name}</TableCell>
            <TableCell component="th" scope="row">{statusJSON.phase}</TableCell>
            <TableCell align="right">
                <ViewObject
                    objectKind='AlluxioCluster'
                    objectName={name}
                    objectConfigString={JSONtoStringifyYAML(specJSON)}
                    objectStatusString={JSONtoStringifyYAML(statusJSON)}
                    handleSendRequest={handleSendRequest}
                />
            </TableCell>
        </TableRow>
    );
}

export default AlluxioCluster
