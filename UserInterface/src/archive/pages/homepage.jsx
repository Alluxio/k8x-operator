import React from 'react';
import { Link } from 'react-router-dom';

export default function HomePage() {
    return (
        <div className="HomePage">
            This is HomePage

            <ul>
                <li><Link to="/">Home</Link></li>
                <li><Link to="/Alluxio_Cluster">Alluxio Cluster</Link></li>
                <li><Link to="/Alluxio_Dataset">Alluxio Dataset</Link></li>
            </ul>

        </div>
    );
}
