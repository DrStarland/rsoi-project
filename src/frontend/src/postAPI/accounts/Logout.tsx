import axios from "axios";
import { backUrl } from "..";

interface resp {
    status: number
}

export const Logout = async function(): Promise<resp> {
    const response = await axios.post(backUrl + `/accounts/logout`)

    return {
        status: response?.status,
    };
}
