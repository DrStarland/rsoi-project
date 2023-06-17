import { Box } from "@chakra-ui/react";
import React from "react";
import NoteCard from "../NoteCard";
import { AllNotesResp } from "postAPI"

import styles from "./NoteMap.module.scss";

interface NoteBoxProps {
    searchQuery?: string
    getCall: () => AllNotesResp
}

type State = AllNotesResp

class NoteMap extends React.Component<NoteBoxProps, State> {
    constructor(props) {
        super(props);
    }

    getAll = () => {
        console.log("state", this.state.items)
        var data = this.props.getCall()
        if (data)
            this.setState(data)
        console.log("state", this.state.items)
    }

    componentDidMount() {
        this.getAll()
    }

    ate(prevProps) {
        if (this.props.searchQuery !== prevProps.searchQuery) {
            this.getAll()
        }
    }

    render() {
        return (
            <Box className={styles.map_box}>
                {this.state?.items.map(item => <NoteCard {...item} key={item.id}/>)}
            </Box>
        )
    }
}

export default React.memo(NoteMap);
