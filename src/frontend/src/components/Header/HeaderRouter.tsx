import React from "react";
import { Routes, Route } from "react-router";
import { BrowserRouter } from "react-router-dom";
import Header from ".";
import NoteHeader from "./Note";


export const HeaderRouter: React.FC<{}> = () => {
    return <BrowserRouter>
        <Routes>
            <Route path="/" element={<Header title="Все заметки" />} />
            <Route path="/auth/signin" element={<Header title="Вход" />} />
            <Route path="/auth/signup" element={<Header title="Регистрация" undertitle="Регистрация нового пользователя" />} />
            <Route path="/notes" element={<Header title="Все заметки" />} />
            <Route path="/notes/:id" element={<NoteHeader title="" />} />
            <Route path="/tickets" element={<Header title="Личный кабинет" />} />
            <Route path="/statistics" element={<Header title="Статистика сервиса" />} />

            <Route path="*" element={<Header title="Страница не найдена" />} />
        </Routes>
    </BrowserRouter>
}
