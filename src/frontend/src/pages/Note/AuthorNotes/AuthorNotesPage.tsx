import { Box } from "@chakra-ui/react";
import React from "react";
import { useCookies } from "react-cookie";
import { useParams } from "react-router-dom";
import NoteMap from "../../../components/NoteMap/NoteMap";

import styles from "./AuthorNotesPage.module.scss";

interface AuthorNotesProps {}

const AuthorNotes: React.FC<AuthorNotesProps> = (props) => {
  let [cookie] = useCookies(["role", "login"]);
  const params = useParams();
  let login = params.login ? params.login : cookie.login;

  return (
    <Box className={styles.main_box}>
      <NoteMap getCall={() => {}} />
    </Box>
  );
};

export default AuthorNotes;
