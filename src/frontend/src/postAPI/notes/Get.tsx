import axiosBackend from "..";
import { Note } from "types/Note";

interface resp {
    status: number
    content: Note
}

const GetNote = async function(id: string): Promise<resp> {
    const response = await axiosBackend
        .get(`/notes/${id}`);
    return {
        status: response.status,
        content: response.data as Note
    };
}

export default GetNote
