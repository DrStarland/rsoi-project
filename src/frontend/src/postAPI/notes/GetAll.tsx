import axiosBackend, { AllNotesResp } from "..";

const GetNotes = async function(): Promise<AllNotesResp> {
    const response = await axiosBackend
        .get(`/notes`);
        // .then(result => {
        //     console.log(result)
        //     var fignya = result.data
        //     return result.data;
        // })
    // response.then(function(result) {
    //     console.log(result)
    //     return response.data;
    // })
    return response.data;
}

export default GetNotes
