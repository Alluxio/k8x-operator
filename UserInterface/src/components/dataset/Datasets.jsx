import React from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import Dataset from "./Dataset";
import {useSelector} from 'react-redux';


function Datasets({handleSendRequest}) {
    const datasetList = useSelector(state => state.dataset.datasetList);

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 500 }} aria-label="simple table">
                <TableHead>
                    <TableRow>
                        <TableCell>Dataset Name</TableCell>
                        <TableCell align="left">Phase</TableCell>
                        <TableCell align="right">Bounded Alluxio Cluster</TableCell>
                        <TableCell align="right">Setting</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {datasetList.map((originalJSON) => (
                        <Dataset
                            key={originalJSON['dataset-config']['name']}
                            name={originalJSON['dataset-config']['name']}
                            originalJSON={originalJSON}
                            handleSendRequest={handleSendRequest}
                        />
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
}

export default Datasets;
