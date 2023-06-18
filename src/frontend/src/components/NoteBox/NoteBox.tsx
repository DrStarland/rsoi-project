import React from "react";
import {
    InputProps as IProps,
    Box, Text
} from "@chakra-ui/react";
import Notes from "components/Icons/Notes";

import styles from "./NoteBox.module.scss";

interface InputProps extends IProps {
    data?: string
}

const NoteBox: React.FC<InputProps> = (props) => {
    return (
    <Box className={styles.note_box}> 
        <Box> <Notes className={styles.icon}/> </Box>
        <Text> {props.data} </Text>
    </Box>
    )
}

export default NoteBox;