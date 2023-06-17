import axiosBackend from "..";

interface resp {
    status: number
}

type Request = {
    id: number
    Title: string
    Content: string
    CreatedAt: string
}

const CreateNote = async function(note: number, Title: string, Content: string, CreatedAt: string): Promise<resp> {
    let data: Request =  {
        id: note,
        Title: Title,
        Content: Content,
        CreatedAt: CreatedAt
    }
    const response = await axiosBackend.post(`/notes`, data);
    return {
        status: response.status
    };
}
export default CreateNote;
