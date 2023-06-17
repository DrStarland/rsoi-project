import axios from "axios";
import { backUrl } from "..";

interface resp {
    status: number
}

const DeleteNote = async function(id: number): Promise<resp> {
    const response = await axios.delete(backUrl + `/notes/${id}`);
    return {
        status: response.status,
    };
}

export default DeleteNote
