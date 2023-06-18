import { Box } from "@chakra-ui/react";
import React from "react";
import NoteCard from "../NoteCard";
import { AllNotesResp } from "postAPI"

import styles from "./NoteMap.module.scss";

interface NoteBoxProps {
    searchQuery?: string
    getCall: () => Promise<AllNotesResp>
}

type State = {
    items: AllNotesResp
}

class NoteMap extends React.Component<NoteBoxProps, State> {
    constructor(props) {
        super(props);
    }

    async getAll() {
        var data = await this.props.getCall()
        console.log("data:", data)
        if (data)
            this.setState({ items: data})
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
        console.log("v rendere:", this.state)
        console.log([1, 2, 3, 4, 5, 6])
        if (this.state?.items) {
            const results: any = [];
            console.log("setState: ", this.state.items)


  // 👇️ can use forEach outside of your JSX code
            // if you need to call a function once for each array element
            this.state.items.forEach((item, index) => {
                console.log("note ", item);
                results.push(
                    <NoteCard {...item} key={item.id}/>
                );
            });
            return results;


            // console.log("state", this.state)
            // return (
            //     <Box className={styles.map_box}>
            //         {this.state.map(item => <NoteCard {...item} key={item.id}/>)}
            //     </Box>
            // )
        } else {
            return (<tr></tr>)
        }
        
        // return (
        //     // <Box className={styles.map_box}>
        //     //     {this.state?.items.map(item => <NoteCard {...item} key={item.id}/>)}
        //     // </Box>
        // )
    }
}

const Child = ({data}) => (
    // <tr>
    //   {data.map((x, i) => (<td key={i}>{x}</td>))}
    // </tr>
     <Box className={styles.map_box}>
        {data.map(item => <NoteCard {...item} key={item.id}/>)}
    </Box>
  );


export default React.memo(NoteMap);