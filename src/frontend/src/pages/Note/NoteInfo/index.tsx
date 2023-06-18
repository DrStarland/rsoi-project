import React from "react";
import { useNavigate, useParams } from "react-router-dom";
import NoteInfoPage from "./NoteInfoPage";

const NoteInfo = () => {
    let navigate = useNavigate();
    return (
        <NoteInfoPage match={useParams()} navigate={navigate}/>
    )
}

export default NoteInfo;
