import React, {useState} from "react";
import {Button, Modal, Box, Grid, Typography} from "@mui/material";
import DeleteObject from "./DeleteObject";
import {modalStyle} from "./util";


const ViewObject = ({   objectKind,
                        objectName,
                        objectConfigString,
                        objectStatusString,
                        handleSendRequest
}) => {
    // State using the useState hook
    const [modalOpen, setModalOpen] = useState(false);
    const handleModalOpen = () => setModalOpen(true);
    const handleModalClose = () => setModalOpen(false);

    return (
        <div className='Modal-ViewObject'>
            <Button variant="outlined" onClick={handleModalOpen} size="small">
                Setting
            </Button>

            <Modal
                open={modalOpen}
                onClose={handleModalClose}
                aria-labelledby="modal-modal-title"
                aria-describedby="modal-modal-description"
            >
                <Box sx={modalStyle}>
                    <div className='Setting'>
                        <Grid container spacing={2}>
                            <Grid item lg={12}>
                                <Box sx={{
                                    height: '3vh',
                                    fontSize: '28px',
                                    width: '100%',}}>
                                    {objectName} [{objectKind}]
                                </Box>
                            </Grid>

                            <Grid item lg={5}>
                                <Box sx={{
                                    fontSize: '20px',
                                    width: '100%',
                                    height: '3vh'
                                }}>
                                    Configuration
                                </Box>
                                <Box sx={{
                                    height: '70vh',
                                    overflowY: 'auto',
                                    width: '100%',
                                }}>
                                    <Typography
                                        component="pre"
                                        sx={{
                                            whiteSpace: 'pre-wrap',
                                            overflowX: 'auto',
                                            backgroundColor: '#f0f0f0',
                                            padding: '10px',
                                            borderRadius: '4px',
                                            fontFamily: 'monospace',
                                        }}
                                    >
                                        {objectConfigString}
                                    </Typography>
                                </Box>
                            </Grid>


                            <Grid item lg={5}>
                                <Box sx={{
                                    fontSize: '20px',
                                    width: '100%',
                                    height: '3vh'
                                }}>
                                    Status
                                </Box>
                                <Box sx={{
                                    height: '70vh',
                                    overflowY: 'auto',
                                    width: '100%',
                                }}>
                                    <Typography
                                        component="pre"
                                        sx={{
                                            whiteSpace: 'pre-wrap',
                                            overflowX: 'auto',
                                            backgroundColor: '#f0f0f0',
                                            padding: '10px',
                                            borderRadius: '4px',
                                            fontFamily: 'monospace',
                                        }}
                                    >
                                        {objectStatusString}
                                    </Typography>
                                </Box>
                            </Grid>

                            <Grid item lg={2}>
                                <Box>
                                    <DeleteObject
                                        objectKind={objectKind}
                                        objectName={objectName}
                                        handleSendRequest={handleSendRequest}
                                    />
                                </Box>
                            </Grid>
                        </Grid>

                    </div>
            </Box>
            </Modal>
        </div>
    )
}

export default ViewObject
