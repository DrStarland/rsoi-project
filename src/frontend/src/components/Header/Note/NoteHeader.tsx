import React from "react";
import Header from "components/Header";
import { TitlesProps } from "components/Header/Titles/Titles";
import { useParams } from "react-router-dom";
import AuthorBox from "components/Boxes/Author";
import GetNote from "postAPI/notes/Get";


const NoteHeader: React.FC<TitlesProps> = (props) => {
    let {
        title = props.title,
        ...rest
    } = props
    const qParams = useParams()

    return (
        <Header title={"Заметка № " + qParams.id || ""} {...rest}/>
    );
}

export default React.memo(NoteHeader);
