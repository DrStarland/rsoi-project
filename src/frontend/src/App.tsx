import * as React from "react";
import theme from "./styles/extendTheme";
import { ChakraProvider, Box, Container } from "@chakra-ui/react";
import { Routes, Route } from "react-router";
import { BrowserRouter } from "react-router-dom";

import Login from "pages/Login";
import SignUp from "pages/Signup";

import { HeaderRouter } from "components/Header";
import AllNotesPage from "pages/Note/AllNotes/AllNotesPage";
import NoteInfoPage from "pages/Note/NoteInfo";
import AllTicketsPage from "pages/Ticket/AllTickets";
import StatisticsPage from "pages/Statistics/StatisticsPage";


interface HomeProps { }
const Home: React.FC<HomeProps> = () => {
  return (
    <Box backgroundColor="bg" h="auto">
      <Container maxW="1000px" minH="95%"
        display="flex"
        paddingX="0px" paddingY="30px"
        alignSelf="center" justifyContent="center"
        textStyle="body"
      >
        <Routing />
      </Container>
    </Box>
  );
};

function Routing() {
  return <BrowserRouter>
    <Routes>
      <Route path="/" element={<AllNotesPage />} />
      <Route path="/auth/signin" element={<Login />} />
      <Route path="/auth/signup" element={<SignUp />} />
      <Route path="/notes" element={<AllNotesPage />} />
      <Route path="/notes/:id" element={<NoteInfoPage />} />
      <Route path="/tickets" element={<AllTicketsPage />} />
      <Route path="/statistics" element={<StatisticsPage />} />

      <Route path="*" element={<NotFound />} />
    </Routes>
  </BrowserRouter>
}

function NotFound() {
  return <h1>404 Page Not Found</h1>
}

export const App = () => {
  return (
    <ChakraProvider theme={theme}>
        <HeaderRouter />
        <Home />
    </ChakraProvider>
  )
};
