import axiosBackend, { AllNotesResp } from "..";

const GetNotes = async function(): Promise<AllNotesResp> {
    const response = await axiosBackend
        .get(`/notes`);
    return  response.data;
}

export default GetNotes