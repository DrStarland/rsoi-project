import React from "react";
import { Box, Link } from "@chakra-ui/react";

import styles from "./Navbar.module.scss";

import { Logout as LogoutQuery } from "postAPI/accounts/Logout";
import { Login as LoginQuery } from "postAPI/accounts/Login";
import { navItems } from "./items";
import AuthActions from "./AuthActions";
import NoauthActions from "./NoauthActions";

export interface NavbarProps {}

const Navbar: React.FC<NavbarProps> = () => {
    let login = "";
    let role = localStorage.getItem("role") || "";

    const [items, ] = React.useState(navItems[role]);

    // React.useEffect(() => {
    //     console.log(items); 
    //   }, [items]);

    const logout = () => {
        localStorage.clear();
        LogoutQuery();
        window.location.href = '/auth/signin';
    };

    return (
        <Box className={styles.navbar}>
        <Box className={styles.navpages}> {items?.map(item =>
            <Link key={item.name} href={item.ref}> {item.name} </Link>
        )} </Box>
        { role !== '' && <AuthActions login={login} logout={logout} />}
        { role === '' && <NoauthActions /> }
    </Box>
    )
}

export default React.memo(Navbar);
