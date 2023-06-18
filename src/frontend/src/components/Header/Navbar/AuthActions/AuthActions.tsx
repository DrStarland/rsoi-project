import React from "react";
import { Box, Text, Link } from "@chakra-ui/react";

import styles from "../Navbar.module.scss";
import LogoutIcon from "components/Icons/Logout";

export interface AuthActionsProps {
    login: string
    logout: () => void
}
const AuthActions: React.FC<AuthActionsProps> = (props) => {
    let role = localStorage.getItem("role") || "";

    return (
        <Box className={styles['user-act']}>
            { role != "" ? (
            <Link onClick={props.logout}> <LogoutIcon/> </Link>
            ) : (
            <Link onClick={() => window.location.href = '/auth/signin'}> Войти </Link>
            )}
        </Box>
    )
}

export default React.memo(AuthActions);