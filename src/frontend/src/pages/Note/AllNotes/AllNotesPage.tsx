import { VStack, Box } from "@chakra-ui/react";
import RoundButton from "components/RoundButton/RoundButton"
import React, { useContext } from "react";
import { SearchContext } from "context/Search";
import NoteMap from "../../../components/NoteMap/NoteMap";

import styles from "./AllNotesPage.module.scss";
import GetNotes from "postAPI/notes/GetAll";

interface AllNotesProps {}

const AllNotesPage: React.FC<AllNotesProps> = (props) => {
  const searchContext = useContext(SearchContext);
  return (
    <>
    {/* Здесь сверху нужна кнопка добавления заметки, такая, как в NoteInfoPage.tsx из папки pages/Note/NoteInfo
    <RoundButton className={styles.basics_button} type="submit" onClick={event => this.submit(event)}>
      Создать заметку
    </RoundButton> */}
    
    <Box className={styles.main_box}>
      {console.log(GetNotes())}
    <NoteMap getCall={GetNotes} />
    </Box>
    </>
  );
};

export default AllNotesPage;
