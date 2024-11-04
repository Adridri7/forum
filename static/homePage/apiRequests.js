import { UserInfo } from "./app.js";
import { clearList, fetchAllUsers } from "./displayUser.js";
import { DisplayReport, DisplayRequest, initRequestEventListeners } from "./request.js";


let post_uuid_report = null

export async function CreateRequest(user_uuid, content) {
    const data = {
        user_uuid: user_uuid,
        content: content,
        post_uuid: post_uuid_report
    };

    try {
        const response = await fetch('/api/requests/createRequest', {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });

        // Vérifie si la réponse est OK
        await checkResponseStatus(response);


        alert("Request created successfully!")
        post_uuid_report = null
    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}

export async function FetchAdminRequest() {
    const request = document.querySelector('.users-request[data-section="moderation"]');
    clearList(request);
    try {
        const response = await fetch('/api/requests/fetchRequest', {
            headers: {
                "Content-Type": "application/json"
            },
        });

        await checkResponseStatus(response);

        const datas = await response.json();
        for (const data of datas) {
            if (data.post_content !== "") {
                const post = await FetchSinglePosts(data.post_uuid);
                DisplayReport(post[0], data);
            } else {
                DisplayRequest(data);
            }
        }

        fetchAllUsers();
        initRequestEventListeners();
    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}

export async function MarkAsReadRequest(request_uuid) {

    try {
        const response = await fetch('/api/requests/isreadRequest', {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ request_uuid })
        });

        // Vérifie si la réponse est OK
        await checkResponseStatus(response);

    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}

async function FetchSinglePosts(post_uuid) {
    try {
        const response = await fetch('/api/post/getPostDetails', {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ post_uuid })
        });

        // Vérifie si la réponse est OK
        await checkResponseStatus(response);

        const data = response.json();

        return data;

    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}

export async function FetchHistoryRequest() {
    const historyRequest = document.getElementById('history-users-request');
    historyRequest.innerHTML = '';
    const user_uuid = UserInfo.user_uuid
    try {
        const response = await fetch('/api/requests/historyRequest', {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ user_uuid })
        });

        // Vérifie si la réponse est OK
        await checkResponseStatus(response);

        const datas = await response.json();

        for (const data of datas) {
            if (data.post_content !== "") {
                const post = await FetchSinglePosts(data.post_uuid);
                DisplayReport(post[0], data, true);
            } else {
                DisplayRequest(data, true);
            }
        }

    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}

async function checkResponseStatus(response) {
    if (!response.ok) {
        const errorMessage = await response.text();
        throw new Error(`Error: ${response.status} - ${errorMessage}`);
    }
}