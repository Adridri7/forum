import { fetchAllUsers } from "./displayUser.js";
import { DisplayRequest } from "./request.js";

export async function CreateRequest(user_uuid, content) {
    console.log("content de la request", content)
    const data = {
        user_uuid: user_uuid,
        content: content,
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
        if (!response.ok) {
            const errorMessage = await response.text();
            throw new Error(`Error: ${response.status} - ${errorMessage}`);
        }

        console.log('Request created successfully');
        alert("Request created successfully!")
    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}

export async function FetchAdminRequest() {
    try {
        const response = await fetch('/api/requests/fetchRequest', {
            headers: {
                "Content-Type": "application/json"
            },
        });

        // Vérifie si la réponse est OK
        if (!response.ok) {
            const errorMessage = await response.text();
            throw new Error(`Error: ${response.status} - ${errorMessage}`);
        }

        const datas = await response.json();
        datas.forEach(data => {
            DisplayRequest(data)
        });

        fetchAllUsers();
        console.log('Request recu', data);
    } catch (error) {
        console.error('There was a problem with the fetch operation:', error);
    }
}