import { Box, HStack, Link, Text, VStack } from "@chakra-ui/react";
import React from "react";

import { Note as NoteI } from "types/Note";

import styles from "./NoteCard.module.scss";

interface NoteProps extends NoteI {}

const NoteCard: React.FC<NoteProps> = (props) => {
  var path = "/notes/" + props.id;

  return (
    <Link className={styles.link_div} href={path}>
      <Box className={styles.main_box}>
        <Box className={styles.info_box}>
          <VStack>
            <Box className={styles.description_box}>
              <Text>{props.id}</Text>
            </Box>
            <Box className={styles.title_box}>
              <Text>{props.Title}➤</Text>
            </Box>
            <Box className={styles.title_box}>
              <Text>➤{props.Content}</Text>
            </Box>

            <HStack>
              <Box className={styles.description_box}>
                <Text>{"Create"}</Text>
              </Box>
            </HStack>
          </VStack>
        </Box>
      </Box>
    </Link>
  );
};

export default NoteCard;
