import {postAuthorize, postLogin} from "../api/api";



export async function clickLogin(){

    let loginPage = document.getElementById("LogIn");
    let inputs = Array.from(loginPage.getElementsByTagName('input'));

    let data = {};

    for (const index in inputs) {
        const element = inputs[index];
        if (!element.checkValidity() || element.value == ''){
            return;
        }
        data[element.name] = element.value;

    };
    const responce = await postLogin(JSON.stringify(data));
    console.log(responce.data);
    return responce;
 
}

export async function clickAuthorize(){
    let authorizePage = document.getElementById("Authorize");
    let inputs = Array.from(authorizePage.getElementsByTagName('input'));

    let data = {};

    for (const index in inputs) {
        const element = inputs[index];

        if (element.type == "checkbox"){
            data[element.name] = element.checked;
        }
        else if (!element.checkValidity() || element.value == ''){
            return;
        }
        else{
            data[element.name] = element.value;
        }
    };

    await postAuthorize(JSON.stringify(data));
    

}