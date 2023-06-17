import axiosBackend, { AllFilghtsResp } from "..";

const GetNotes = async function(): Promise<AllFilghtsResp> {
    const response = await axiosBackend
        .get(`/notes`);
    return  response.data;
}

export default GetNotes
