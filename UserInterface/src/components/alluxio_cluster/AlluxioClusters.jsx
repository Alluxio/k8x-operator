import React from 'react';
import AlluxioCluster from "./AlluxioCluster";
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import {useSelector} from 'react-redux';

const AlluxioClusters = ({ handleSendRequest }) => {
    const alluxioList = useSelector(state => state.alluxio.alluxioList);

    return (
        <TableContainer component={Paper}>
            <Table sx={{ minWidth: 650 }} aria-label="simple table">
                <TableHead>
                    <TableRow>
                        <TableCell>Alluxio Cluster Name</TableCell>
                        <TableCell align="left">Phase</TableCell>
                        <TableCell align="right">Setting</TableCell>
                    </TableRow>
                </TableHead>

                <TableBody>
                    {alluxioList.map((originalJSON) => (
                        <AlluxioCluster
                            key={originalJSON['alluxio-cluster-config']['name']}
                            name={originalJSON['alluxio-cluster-config']['name']}
                            originalJSON={originalJSON}
                            handleSendRequest={handleSendRequest}
                        />
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
}

export default AlluxioClusters;
