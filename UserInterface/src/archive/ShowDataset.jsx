import { Component } from 'react';


class ShowDataset extends Component {

    render() {
        const {hasDataset, datasetConfigJSON, datasetStatusJSON} = this.props;

        if (!hasDataset) {
            return <div className="GetDataset">
                <h2>
                    There is no Alive Dataset.
                </h2>
            </div>;
        } else {
            return (
                <div className="GetDataset">

                    <div className='DS Name'>
                        <h3>Data Set Name:</h3>
                        {datasetConfigJSON.name}
                    </div>

                    <div className='DS Path'>
                        <h3>Data Set Path:</h3>
                        {datasetConfigJSON.path}
                    </div>

                    <div className='DS Status'>
                        <h3>Data Set Status:</h3>
                        <h4>Bounded Alluxio Cluster: </h4>
                        {datasetStatusJSON?.boundedAlluxioCluster}
                        <h4>Phase: </h4>
                        {datasetStatusJSON?.phase}
                    </div>

                </div>
            );
        }
    }
}

export default ShowDataset
