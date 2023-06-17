import { Box } from "@chakra-ui/react";
import React from "react";
import NoteCard from "../NoteCard";
import { AllFilghtsResp } from "postAPI"

import styles from "./NoteMap.module.scss";

interface NoteBoxProps {
    searchQuery?: string
    getCall: (page: number, size: number) => Promise<AllFilghtsResp>
}

type State = AllFilghtsResp

class NoteMap extends React.Component<NoteBoxProps, State> {
    constructor(props) {
        super(props);
    }

    async getAll() {
        var data = await this.props.getCall(0, 20)
        if (data)
            this.setState(data)
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
