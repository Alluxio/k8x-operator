import React from 'react';
import AlluxioController from "./containers/alluxio-controller";
import DatasetController from "./containers/dataset-controller";
import {Grid} from "@mui/material";
import OperatorAlert from "./components/alert/Alert";


function App() {
    return (
        <div className="HOME">
            <Grid container spacing={2}>
                <Grid item lg={0.25} xs={0}/>

                <Grid item lg={5.75} xs={12}>
                    <div className="component" id="aaa">
                        <AlluxioController/>
                    </div>
                </Grid>

                <Grid item lg={5.75} xs={12}>
                    <div className="component" id="bbb">
                        <DatasetController/>
                    </div>
                </Grid>

                <Grid item lg={0.25} xs={0}/>
            </Grid>

            <OperatorAlert/>
        </div>
    );
}

export default App;