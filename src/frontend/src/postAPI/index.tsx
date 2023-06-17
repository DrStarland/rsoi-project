import { Category } from "types/Categories";
import { Note } from "types/Note";
import { Account } from "types/Account";
import axios from "axios";

export const backUrl = "http://localhost:8080/api/v1";
// export const backUrl = "http://gateway-service:8080/api/v1";

const axiosBackend = () => {
    let instance = axios.create({
        baseURL: backUrl
    });

    instance.interceptors.request.use(function (config) {
        const token = localStorage.getItem("authToken");
        if (config.headers && token) {
            config.headers.Authorization = 'Bearer ' + token;
        }

        return config;
    });

    return instance
};

export default axiosBackend();

export type AllNotesResp = {
    // page: number,
    // pageSize: number, 
    // totalElements: number,
    items: Note[]
}

export type AllCategoriesResp = {
    status: number,
    content: Category[]
}

export type AllUsersResp = {
    status: number,
    content: Account[]
}
