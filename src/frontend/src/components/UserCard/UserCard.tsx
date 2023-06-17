import { Box, HStack, Image, Link, Text, VStack } from "@chakra-ui/react";
import React from "react";

import { Account as AccountI } from "types/Account";

import RoleBox from "components/Boxes/Role/RoleBox";
import LoginBox from "components/LoginBox";
import NoteBox from "components/NoteBox";
import user from "img/user.png";

import styles from "./UserCard.module.scss";

interface UserProps extends AccountI {
    role: string
}

const UserCard: React.FC<UserProps> = (props) => {
  const pathNotes = "/accounts/" + props.login + "/notes";

  return (
    <Box className={styles.main_box}>
      <HStack>
        <VStack>
          <Image src={user} />

          <VStack className={styles.role}>
            <Text>Роль</Text>
            <RoleBox login={props.login} role={props.role}/>
          </VStack>
        </VStack>

        <VStack className={styles.info}>
          <Text>Логин</Text>
          <LoginBox login={props.login} className={styles.login} />

          <Link href={pathNotes}>
            <NoteBox data={"Рецепты"} className={styles.notes} />
          </Link>

        </VStack>
      </HStack>
    </Box>
  );
};

export default UserCard;
